package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"task/utils"
	"time"
)

var (
	help bool
	env string //"env6"
	backupDestBaseDir string // "/tmp/pods/"
)

func init()  {
	flag.BoolVar(&help, "help", false, "zipfile help usage")
	flag.StringVar(&env, "env", "env6","pod's env")
	flag.StringVar(&backupDestBaseDir, "dst", "/tmp/pods/","backup destination directory")
	//flag.StringVar(&restApiUrl, "url", "http://192.32.14.181:30954/","url which query pod's metadata")
	flag.Usage = usage
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `zipfile version: 1.0
Options:
`)
	flag.PrintDefaults()
}

func main() {
	//所有文件备份完成后，开始执行备份目标文件夹的压缩
	//首先判断要压缩的目录是否存在
	flag.Parse()
	if help {
		flag.Usage()
	} else {
		dayFormat := time.Now().Format("2006-01-02")

		backupDestBaseDir = utils.PathWrapper(backupDestBaseDir)

		ok, err := utils.CheckExists(backupDestBaseDir + env)
		if err != nil {
			log.Printf("检查目录：%v 报错: %v\n",backupDestBaseDir + env,err)
		}
		//如果/tmp/log-backup中存在env目录，就开始执行压缩
		if ok {
			err = utils.ZipFile(backupDestBaseDir + env, backupDestBaseDir + env + "-" + dayFormat + ".zip")
			if err == nil {
				//压缩完成后，删除对应的env目录
				err = os.RemoveAll(backupDestBaseDir + env)
				if err != nil {
					log.Printf("删除：%v 报错: %v\n",backupDestBaseDir + env,err)
				}
			} else {
				log.Printf("创建压缩文件报错: %v\n",err)
			}
		}
	}


}
