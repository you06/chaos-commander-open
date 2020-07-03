package api

import (
	"github.com/kataras/iris"
	"github.com/ngaut/log"
)

func (hdl *ManagerHandler)GetHistoryList(ctx iris.Context) {
	histories, err := hdl.mgr.GetHistoryList()
	if err != nil {
		log.Errorf("get history list err %v", err)
	}
	ctx.JSON(histories)
}

func (hdl *ManagerHandler)GetHistoryById(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		log.Errorf("Invalid id %s", ctx.Params().Get("id"))
		ctx.WriteString("invalid id")
		return
	}
	history, err := hdl.mgr.GetHistoryById(id)
	if err != nil {
		log.Errorf("get history err %v", err)
		ctx.WriteString("history not found")
	}
	ctx.JSON(history)
}
