package models

type ChainDomain struct {
	NodeIps   []string          `json:"nodeIps"`
	NodePorts map[string]string `json:"nodePorts"`
}

func NewChainDomain(ips []string, ports map[string]string) ChainDomain {
	return ChainDomain{
		NodeIps:   ips,
		NodePorts: ports,
	}
}
