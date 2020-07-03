package api

import (
	"github.com/kataras/iris"
	"github.com/ngaut/log"
)

func (hdl *ManagerHandler)GetAnsibleList(ctx iris.Context) {
	ansibles, err := hdl.mgr.GetAnsibleList()
	if err != nil {
		log.Errorf("get ansible list err %v", err)
	}
	ctx.JSON(ansibles)
}

func (hdl *ManagerHandler)GetAnsibleById(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		log.Errorf("Invalid id %s", ctx.Params().Get("id"))
		ctx.WriteString("invalid id")
		return
	}
	ansible, err := hdl.mgr.GetAnsibleById(id)
	if err != nil {
		log.Errorf("get ansible err %v", err)
		ctx.WriteString("ansible not found")
	}
	ctx.JSON(ansible)
}
