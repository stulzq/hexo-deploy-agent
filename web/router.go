package web

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stulzq/hexo-deploy-agent/util"
	"github.com/stulzq/hexo-deploy-agent/web/handlers/deploy"
)

func registerRoute(h *server.Hertz) {
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(consts.StatusOK, "pong")
	})

	h.GET("/h", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(consts.StatusOK, util.BytesToStr(ctx.Request.Header.RawHeaders()))
	})

	h.POST("/deploy/upload", deploy.Upload)
}
