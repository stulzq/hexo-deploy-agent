package deploy

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	terrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/stulzq/hexo-deploy-agent/logger"
	"github.com/stulzq/hexo-deploy-agent/util"
)

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

	stepUnzip(fileSavePath)

	stepRefreshCDN()

	stepPushMsg()

	stepClean(fileSavePath)

	logger.Info("[Deploy][Job] successfully")
}

func stepUnzip(fileSavePath string) {
	unzipTargetDir := deployConf.BlogDir
	if util.IsDebug() {
		unzipTargetDir = unzipTmp
	}

	if err := util.UnZip(fileSavePath, unzipTargetDir); err != nil {
		logger.Error("[Deploy][Job] unzip err,", err)
		return
	}

	logger.Info("[Deploy][Job] unzip success to path: ", unzipTargetDir)
}

func stepRefreshCDN() {
	if !deployConf.CDN.Enable {
		return
	}

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

func stepPushMsg() {
	if !deployConf.Dingtalk.Enable {
		return
	}

	const msg = `{"msgtype": "text","text": {"content":"Deploy successfully!"}}`

	// send msg https://open.dingtalk.com/document/group/custom-robot-access
	if _, err := http.Post(deployConf.Dingtalk.Url, "application/json", strings.NewReader(msg)); err != nil {
		logger.Error("[Deploy][Job] send msg failed, ", err)
	}
}

func stepClean(fileSavePath string) {
	os.Remove(fileSavePath)
	os.RemoveAll(unzipTmp)
}
