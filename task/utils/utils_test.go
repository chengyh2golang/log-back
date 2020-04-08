package utils

import (
	"fmt"
	"testing"
)

const (

	dirName1 = ".abc"
	dirName2 = "abc"
	dirName3 = "/tmp/log-bak"
	dirName4 = "/tmp/log-bak/"
	fileName = "icore-service-uaa-7484d8d4d8-5zr7f.2020-03-12.0.log"
	fileName1 = "command-center.log.2020-03-28.0"
	fileName2 = "command-center.log.2020-03-28.0.lck"
	path = "/mnt/paas/kubernetes/kubelet/pods/77e00ad0-7033-11ea-bfe7-000c2999f0e6/volumes/kubernetes.io~empty-dir/app-logs/icore-service-uaa-7484d8d4d8-5zr7f.2020-03-12.0.log"
)



func TestIsNeedBackup(t *testing.T) {
	//fmt.Println(IsNeedBackup(fileName,7))
	fmt.Println(CheckEndWithDotLog(fileName))
	fmt.Println(CheckEndWithDotLog(fileName1))
	fmt.Println(CheckEndWithDotLog(fileName2))
}

func TestPathInputWrapper(t *testing.T) {

	fmt.Println(PathWrapper(dirName4))
	fmt.Println(PathWrapper(dirName3))
}

//func TestFetchPodMetaDataByEnvAndPodId(t *testing.T) {
//	info, err := FetchDestPathByEnvAndPodId("env6", "77e00ad0-7033-11ea-bfe7-000c2999f0e6")
//	if err != nil {
//		t.Errorf("fetch pod metadata failed: %v",err)
//	}
//	fmt.Printf("%+v",info)
//}


