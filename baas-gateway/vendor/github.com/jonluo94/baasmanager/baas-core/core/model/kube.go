package model

type K8sData struct {
	Data [][]byte `json:"data"`
}

type ChainDomain struct {
	NodeIps   []string          `json:"nodeIps"`
	NodePorts map[string]string `json:"nodePorts"`
}

type ChainPod struct {
	Status     string `json:"status"`
	CreateTime string `json:"createTime"`
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	HostIP     string `json:"hostIP"`
	Type       string `json:"type"`
	Cpu        string `json:"cpu"`
	Memory     string `json:"memory"`
}

type Resources struct {
	Node string `json:"node"`
	CPU  float64 `json:"cpu"`
	Memory int `json:"memory"`
}