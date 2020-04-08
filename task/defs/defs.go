package defs

const (
	//Env = "env6"
	//PodBaseDir = "/mnt/paas/kubernetes/kubelet/pods/"
	//BackupDestBaseDir = "/tmp/pods/"
	//RestApiUrl = "http://192.168.250.22:32598/online/podmetadata"
	//Expired = 1
	//LogDestDir = "/tmp/log-back"
	EmptyDirName = "kubernetes.io~empty-dir"
	UrlSuffix= "online/podmetadata"
)

type PodInfo struct {
	PodName    string `json:"pod_name"`
	//PodId      string `json:"pod_id"`
	RsName     string `json:"rs_name"`
	DeployName string `json:"deploy"`
	Namespace  string `json:"namespace"`
}

