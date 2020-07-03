package api

import (
	"github.com/kataras/iris"
	"github.com/ngaut/log"
)

func (hdl *ManagerHandler)GetJobsList(ctx iris.Context) {
	jobs, err := hdl.mgr.GetJobsList()
	if err != nil {
		log.Errorf("get jobs err %v", err)
	}
	ctx.JSON(jobs)
}

func (hdl *ManagerHandler)GetJobById(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		log.Errorf("Invalid id %s", ctx.Params().Get("id"))
		ctx.WriteString("invalid id")
		return
	}
	job, err := hdl.mgr.GetJobById(id)
	if err != nil {
		log.Errorf("get job err %v", err)
		ctx.WriteString("job not found")
	}
	ctx.JSON(job)
}
