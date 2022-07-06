package model

type DeployConfig struct {
	BlogDir string `yaml:"blog_dir"`
	CDN     struct {
		Enable       bool      `yaml:"enable"`
		AccessKey    string    `yaml:"accessKey"`
		SecretKey    string    `yaml:"secretKey"`
		RefreshPaths []*string `yaml:"refreshPaths"`
		FlushType    *string   `yaml:"flushType"`
	} `yaml:"cdn"`
	Dingtalk struct {
		Enable bool   `yaml:"enable"`
		Url    string `yaml:"url"`
	} `yaml:"dingtalk"`
}
