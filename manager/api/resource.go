package api

import (
	"github.com/kataras/iris"
	"github.com/ngaut/log"
)

func (hdl *ManagerHandler)GetResourceList(ctx iris.Context) {
	resources, err := hdl.mgr.GetResourceList()
	if err != nil {
		log.Errorf("get resource list err %v", err)
	}
	ctx.JSON(resources)
}

func (hdl *ManagerHandler)GetResourceById(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		log.Errorf("Invalid id %s", ctx.Params().Get("id"))
		ctx.WriteString("invalid id")
		return
	}
	resource, err := hdl.mgr.GetResourceById(id)
	if err != nil {
		log.Errorf("get resource err %v", err)
		ctx.WriteString("resource not found")
	}
	ctx.JSON(resource)
}
