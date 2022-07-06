package deploy

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/pkg/errors"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	terrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	"github.com/stulzq/hexo-deploy-agent/config"
	"github.com/stulzq/hexo-deploy-agent/logger"
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
	// load config
	var conf model.DeployConfig
	config.GetStruct("deploy", &conf)
	deployConf = &conf
	if util.IsDebug() {
		deployConf.CDN.AccessKey = os.Getenv("AK")
		deployConf.CDN.SecretKey = os.Getenv("SK")
	}

	// create dir
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		panic(errors.Wrap(err, "mkdir err"))
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

// Upload upload deploy file(zip)
func Upload(_ context.Context, c *app.RequestContext) {
	// get file
	file, err := c.FormFile("f")
	if err != nil {
		logger.Error("[Web][Deploy][Upload] receive file err,", err)
		c.String(500, err.Error())
		return
	}

	// check ext
	ext := filepath.Ext(file.Filename)
	if strings.ToLower(ext) != ".zip" {
		c.String(400, "file ext now allow")
		return
	}

	// save file
	fileSavePath := filepath.Join(uploadPath, file.Filename)

	logger.Info("[Web][Deploy][Upload] receive file: ", file.Filename)
	err = c.SaveUploadedFile(file, fileSavePath)

	if err != nil {
		logger.Error("[Web][Deploy][Upload] save file err,", err)
		c.String(500, err.Error())
		return
	}

	c.String(200, "ok")

	// deploy
	go processDeploy(fileSavePath)
}

func processDeploy(fileSavePath string) {
	logger.Info("[Deploy][Job] start deploy ", fileSavePath)

	// unzip
	unzipTargetDir := deployConf.BlogDir
	if util.IsDebug() {
		unzipTargetDir = unzipTmp
	}
	if err := util.UnZip(fileSavePath, unzipTargetDir); err != nil {
		logger.Error("[Deploy][Job] unzip err,", err)
		return
	}

	logger.Info("[Deploy][Job] unzip success to path: ", unzipTargetDir)

	// cdn
	if deployConf.CDN.Enable {
		request := cdn.NewPurgePathCacheRequest()
		request.Paths = deployConf.CDN.RefreshPaths
		request.FlushType = deployConf.CDN.FlushType

		response, err := cdnClient.PurgePathCache(request)
		if _, ok := err.(*terrors.TencentCloudSDKError); ok {
			logger.Error("[Deploy][Job] an API error has returned: ", err)
			return
		}
		if err != nil {
			logger.Error("[Deploy][Job] cdn request err: ", err)
			return
		}

		logger.Info("[Deploy][Job] cdn response: ", response.ToJsonString())
	}

	logger.Info("[Deploy][Job] successfully")

	os.Remove(fileSavePath)
	os.RemoveAll(unzipTmp)
}
