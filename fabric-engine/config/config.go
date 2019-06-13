package config

type UserBaasConfig struct {
	ArtifactPath  string
	K8sConfigPath string
	DataPath      string
	TemplatePath  string
}

func NewUserBaasConfig(artifactPath, k8sConfig, dataPath, templatePath string) *UserBaasConfig {
	return &UserBaasConfig{
		ArtifactPath:  artifactPath,
		K8sConfigPath: k8sConfig,
		DataPath:      dataPath,
		TemplatePath:  templatePath,
	}
}
