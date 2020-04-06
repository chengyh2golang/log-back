package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"task/defs"
	"task/utils"
	"time"
)

var (
	help bool
	env string //"env6"
	podBaseDir string // "/mnt/paas/kubernetes/kubelet/pods/"
	backupDestBaseDir string // "/tmp/pods/"
	restApiUrl string // "http://192.168.250.22:32598/online/podmetadata"
	expired int // 1
	logDestDir  string // "/tmp/backup-task-log"
)

func init()  {
	flag.BoolVar(&help, "help", false, "backup help usage")
	flag.StringVar(&env, "env", "env6","pod's env")
	flag.StringVar(&podBaseDir, "src", "/mnt/paas/kubernetes/kubelet/pods/","pod's home directory")
	flag.StringVar(&backupDestBaseDir, "dst", "/tmp/pods/","backup destination directory")
	flag.StringVar(&restApiUrl, "url", "http://192.168.250.22:32598/online/podmetadata","url which query pod's metadata")
	flag.IntVar(&expired, "expired", 1,"expired time, day")
	flag.StringVar(&logDestDir, "log", "/tmp/backup-task-log","log file directory")
	flag.Usage = usage
}

type archiveLog struct {
	env string
	podBaseDir string
	backupDestDir string
	checkExpired int
	restApiUrl string
	logDestDir string
}

func (a *archiveLog) backup()  {
	//检查指定的log目录是否存在，如果不存在，就创建出这个目录
	logDirExists, err := utils.CheckExists(a.logDestDir)
	if err != nil {
		log.Fatal(err)
	}
	if !logDirExists {
		err := os.MkdirAll(a.logDestDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	dayFormat := time.Now().Format("2006-01-02")
	logFileName := a.logDestDir + "backup-" + dayFormat + ".log"

	exists, err := utils.CheckExists(logFileName)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		file, err := os.Create(logFileName)

		defer func() {
			if err != nil {
				file.Close()
			} else {
				err = file.Close()
			}
		}()
		if err != nil {
			panic(err)
		}
	}

	f, err := os.OpenFile(logFileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	defer f.Close()
	log.SetOutput(f)

	//检查/mnt/paas/kubernetes/kubelet/pods/目录下所有满足条件的文件
	//根据满足条件的文件路径，解析出podId和这个podId对应的需要备份的所有文件列表
	//路径：/mnt/paas/kubernetes/kubelet/pods/
	//路径+：77e00ad0-7033-11ea-bfe7-000c2999f0e6/volumes/kubernetes.io~empty-dir/app-logs
	//文件名：icore-service-uaa-7484d8d4d8-5zr7f.2020-03-12.0.log
	//根据env和podId找restApiUrl获取namespace，deploy，rs，pod名称
	//在备份目标路径上创建目录：/env/namespace/deploy/rs/pod/
	//将备份文件拷贝到目标路径的新建目录上
	//删除本地的文件
	//如果有错，写报错日志到文件中

	//定义一个备份结果数组，用于在log文件中打印信息
	var backupResult []string

	_ = filepath.Walk(a.podBaseDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileName := info.Name()
			//如果info是文件并且path包含有"kubernetes.io~empty-dir"
			//并且info的名字里包含"-20"这样的关键字，就说明这是一个需要做进一步检查是否过期即满足备份的文件

			if !info.IsDir() &&
				strings.Contains(path,defs.EmptyDirName) &&
				strings.Contains(fileName,".20") {

				//检查文件是否需要备份
				if utils.IsNeedBackup(fileName,a.checkExpired) {
					//获取该文件的备份路径
					//根据path获取podID
					podId := utils.FetchPodIdByPath(a.podBaseDir, path)

					//根据env,podId,restApiUrl,backupDestBaseDir获取备份目标路径
					destPath, err := utils.FetchDestPathByEnvAndPodId(a.env, podId, a.restApiUrl, a.backupDestDir)
					if err != nil {
						log.Printf("通过env: %v和pod_id: %v调用api: %v获取信息失败：%v\n",a.env,podId,a.restApiUrl,err)
						return err
					}
					//检查destPath是否存在，如果不存在就创建
					exists, err := utils.CheckExists(destPath)
					if err != nil {
						log.Printf("检查备份目标路径是否存在时，报错：%v\n",err)
						return err
					}
					if !exists {
						err := os.MkdirAll(destPath, os.ModePerm)
						if err != nil {
							log.Printf("创建备份文件夹: %v 失败：%v\n",destPath,err)
							return err
						}
					}
					//备份文件
					err = os.Rename(path, utils.PathWrapper(destPath)+fileName)
					if err != nil {
						log.Printf("执行备份报错: %v\n",err)
						return err
					}
					backupResult = append(backupResult,path)
				}
			}
			return nil
		})
	//输出备份结果信息到log文件中
	if len(backupResult) == 0 {
		log.Println("本次任务，没有找到符合备份条件的文件！")
	} else {
		log.Printf("备份了%v个符合条件的文件:\n",len(backupResult))
		for _,v := range backupResult {
			log.Println(v)
		}
	}

}

func newArchiveLog(env string, podDir string, backupDestDir string, expired int,restApiUrl string,logDir string) *archiveLog {
	return &archiveLog{
		env:env,
		podBaseDir:podDir,
		backupDestDir:backupDestDir,
		checkExpired:expired,
		restApiUrl:restApiUrl,
		logDestDir:logDir,
		}
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
	} else {
		podBaseDir = utils.PathWrapper(podBaseDir)
		backupDestBaseDir = utils.PathWrapper(backupDestBaseDir)
		logDestDir = utils.PathWrapper(logDestDir)
		a := newArchiveLog(env,podBaseDir,backupDestBaseDir,expired,restApiUrl,logDestDir)
		a.backup()
	}

}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `backup version: 1.0
Options:
`)
	flag.PrintDefaults()
}
