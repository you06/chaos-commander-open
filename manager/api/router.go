package api

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/you06/chaos-commander/manager/core"
)

func CreateRouter(app *iris.Application, mgr *core.Manager) {
	hdl := newManagerHandler(mgr)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
	})
	party := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions)

	party.Get("/ping", hdl.Ping)

	// jobs API
	party.Get("/jobs", hdl.GetJobsList)
	party.Get("/job/{id}", hdl.GetJobById)

	// ansibles API
	party.Get("/ansibles", hdl.GetAnsibleList)
	party.Get("/ansible/{id}", hdl.GetAnsibleById)

	// blue print API
	party.Get("/blueprints", hdl.GetBluePrintList)
	party.Get("/blueprint/{id}", hdl.GetBluePrintById)
	party.Put("/blueprint/{id}", hdl.PutBluePrintById)

	// load API
	party.Get("/loads", hdl.GetLoadList)
	party.Get("/load/{id}", hdl.GetLoadById)

	// resource API
	party.Get("/resources", hdl.GetResourceList)
	party.Get("/resource/{id}", hdl.GetResourceById)

	// history API
	party.Get("/histories", hdl.GetHistoryList)
	party.Get("/history/{id}", hdl.GetHistoryById)

	// app.Get("/*", iris.FileServer(mgr.Config.Root, iris.DirOptions{ShowList: false, Gzip: true}))
}
