package utils

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"task/defs"
	"time"
)

const (
	//BaseYearTime = "2006"
	BaseDayTime = "2006-01-02"
)

//检查文件是否以.log为结尾
func CheckEndWithDotLog(inputStr string) bool {
	split := strings.Split(inputStr, ".")
	if split[len(split) -1 ] == "log" {
		return true
	} else {
		return false
	}
}

//拷贝文件
func CopyFile(src, des string) (written int64, err error) {
	//获取源文件的权限
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	fi, _ := srcFile.Stat()
	perm := fi.Mode()
	srcFile.Close()

	input, err := ioutil.ReadFile(src)
	if err != nil {
		return 0, err
	}

	err = ioutil.WriteFile(des, input, perm)
	if err != nil {
		return 0, err
	}
	return int64(len(input)), nil
}

//检查文件或者文件夹是否存在
func CheckExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//判断字符串是否是以某个字符为结尾
//用于判断用户输入的路径是/tmp/back还是/tmp/back/
//如果不是以"/"结尾，加上
func PathWrapper(inputStr string) string {
	if strings.LastIndex(inputStr, "/") == len(inputStr) -1 {
		return inputStr
	} else {
		return inputStr + "/"
	}
}

//func FetchLogfileByFullPath(fullPath string) string {
//	split := strings.Split(fullPath,"/")
//	return split[len(split)-1]
//}

//根据超时时间设置：比如7天
//再根据日志文件名：icore-service-uaa-7484d8d4d8-5zr7f.2020-03-12.0.log,解析出该文件生成的日期
// 判断该文件是否需要备份
func IsNeedBackup(fileName string, expiredDay int) bool {
	//获取该文件的日期
	fileDaytimeStr:=  strings.Split(fileName,".")[1]
	fileDaytime, _ := time.Parse(BaseDayTime, fileDaytimeStr)

	nowDaytimeStr := time.Now().Format(BaseDayTime)
	nowDaytime, _ := time.Parse(BaseDayTime, nowDaytimeStr)

	if (fileDaytime.Unix() + int64(86400 * expiredDay) ) <= nowDaytime.Unix() {
		return true
	} else {
		return false
	}
}

//根据podBaseDir和path解析podId
//baseDir="/mnt/paas/kubernetes/kubelet/pods/"
//path=/mnt/paas/kubernetes/kubelet/pods/77e00ad0-7033-11ea-bfe7-000c2999f0e6/volumes/kubernetes.io~empty-dir/app-logs/icore-service-uaa-7484d8d4d8-5zr7f.2020-03-12.0.log
func FetchPodIdByPath(podBaseDir,path string) string {
	podIdStr := strings.Split(path, podBaseDir)[1]
	return strings.Split(podIdStr,"/")[0]
}

//根据env和podId，调用restApiUrl去获取namespace，deploy，rs，podName
//获取信息主要用于创建备份目标的对应目录
//备份目录的路径为：/tmp/back + /env/namespace/deploy/rs/podname/xxx.log
func FetchDestPathByEnvAndPodId(env,podId,restApiUrl,backupDestBaseDir string) (string,error) {
	var podInfo defs.PodInfo
	result := ""

	urlParam := "?env=" + env + "&pod_id=" + podId
	resp, err := http.Get(restApiUrl + defs.UrlSuffix + urlParam)
	if err != nil {
		return result,err
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &podInfo)
	if err != nil {
		return result,err
	}

	result = backupDestBaseDir + env  + "/" + podInfo.Namespace + "/" +
		podInfo.DeployName + "/" + podInfo.RsName + "/" + podInfo.PodName

	return result,nil
}

// srcFile could be a single file or a directory
func ZipFile(srcFile string, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}


		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile) + "/")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if ! info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

///mnt/paas/kubernetes/kubelet/pods/77e00ad0-7033-11ea-bfe7-000c2999f0e6/volumes/kubernetes.io~empty-dir/

/*
func IsDirNameStartWithDot(name string) bool {
	return strings.HasPrefix(name,".")
}

func FetchAllDir(dir string) []string {
	var result []string
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if f.IsDir() && !IsDirNameStartWithDot(f.Name()) {
			result = append(result,f.Name())
		}
	}
	return result
}

err := filepath.Walk(".",
    func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    fmt.Println(path, info.Size())
    return nil
})
if err != nil {
    log.Println(err)
}
*/


