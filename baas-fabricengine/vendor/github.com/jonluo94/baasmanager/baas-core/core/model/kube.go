package model

type K8sData struct {
	Data [][]byte `json:"data"`
}

type ChainDomain struct {
	NodeIps   []string          `json:"nodeIps"`
	NodePorts map[string]string `json:"nodePorts"`
}