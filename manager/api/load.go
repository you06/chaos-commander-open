package api

import (
	"github.com/kataras/iris"
	"github.com/ngaut/log"
)

func (hdl *ManagerHandler)GetLoadList(ctx iris.Context) {
	loads, err := hdl.mgr.GetLoadList()
	if err != nil {
		log.Errorf("get load list err %v", err)
	}
	ctx.JSON(loads)
}

func (hdl *ManagerHandler)GetLoadById(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		log.Errorf("Invalid id %s", ctx.Params().Get("id"))
		ctx.WriteString("invalid id")
		return
	}
	load, err := hdl.mgr.GetLoadById(id)
	if err != nil {
		log.Errorf("get load err %v", err)
		ctx.WriteString("load not found")
	}
	ctx.JSON(load)
}
