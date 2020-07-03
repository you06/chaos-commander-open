package api

import (
	"github.com/kataras/iris"
	"github.com/ngaut/log"
)

func (hdl *ManagerHandler)GetBluePrintList(ctx iris.Context) {
	bluePrints, err := hdl.mgr.GetBluePrintList()
	if err != nil {
		log.Errorf("get jobs err %v", err)
	}
	ctx.JSON(bluePrints)
}

func (hdl *ManagerHandler)GetBluePrintById(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		log.Errorf("Invalid id %s", ctx.Params().Get("id"))
		ctx.WriteString("invalid id")
		return
	}
	bluePrint, err := hdl.mgr.GetBluePrintById(id)
	if err != nil {
		log.Errorf("get blue print err %v", err)
		ctx.WriteString("blue print not found")
	}
	ctx.JSON(bluePrint)
}

func (hdl *ManagerHandler)PutBluePrintById(ctx iris.Context) {}
