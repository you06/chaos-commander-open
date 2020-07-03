package api

import (
	"github.com/kataras/iris"
	"github.com/you06/chaos-commander/manager/core"
)

// ManagerHandler is manager api handler
type ManagerHandler struct {
	mgr *core.Manager
}

func newManagerHandler(mgr *core.Manager) *ManagerHandler {
	return &ManagerHandler{
		mgr: mgr,
	}
}

func (hdl *ManagerHandler)Ping(ctx iris.Context) {
	ctx.WriteString("pong")
}
