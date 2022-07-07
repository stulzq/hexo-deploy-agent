package deploy

import (
	"os"

	"github.com/pkg/errors"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	"github.com/stulzq/hexo-deploy-agent/config"
	"github.com/stulzq/hexo-deploy-agent/util"
	"github.com/stulzq/hexo-deploy-agent/web/handlers/deploy/model"
)

var (
	deployConf *model.DeployConfig
	cdnClient  *cdn.Client
)

const (
	uploadPath = "deploy/upload/"
	unzipTmp   = "deploy/tmp/"
)

func init() {
	// load deploy config
	deployConf = &model.DeployConfig{}
	_ = config.GetStruct("deploy", deployConf)

	if util.IsDebug() {
		deployConf.CDN.AccessKey = os.Getenv("AK")
		deployConf.CDN.SecretKey = os.Getenv("SK")
		deployConf.Dingtalk.Url = os.Getenv("DINGTALK_URL")
	}

	// create dir
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		panic(errors.Wrap(err, "mkdir upload dir err"))
	}

	// tencentyun cdn client
	credential := common.NewCredential(
		deployConf.CDN.AccessKey,
		deployConf.CDN.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	cdnClient, _ = cdn.NewClient(credential, "", cpf)
}
