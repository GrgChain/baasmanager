package models

type K8sData struct {
	Data [][]byte `json:"data"`
}

type ChainDomain struct {
	NodeIps   []string          `json:"nodeIps"`
	NodePorts map[string]string `json:"nodePorts"`
}

type RespData struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
}

type Namespace struct {
	Namespaces string `json:"namespaces"`
}
