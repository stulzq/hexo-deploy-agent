package web

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/stulzq/hexo-deploy-agent/logger"
)

func Start() {
	hlog.SetLogger(logger.NewHertzFullLogger())
	h := server.Default(server.WithKeepAlive(true), server.WithHostPorts(":9190"), server.WithMaxRequestBodySize(100*1024*1024))

	registerRoute(h)

	h.Spin()
}
