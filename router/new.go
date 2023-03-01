package router

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/andrewarrow/feedback/controller"
	"github.com/andrewarrow/feedback/files"
)

type Router struct {
	Paths    map[string]func(writer http.ResponseWriter)
	Template *template.Template
	Vars     *controller.Vars
	Site     *controller.Site
}

func NewRouter(path string) *Router {
	r := Router{}
	r.Paths = map[string]func(writer http.ResponseWriter){}

	var site controller.Site
	jsonString := files.ReadFile(path)
	json.Unmarshal([]byte(jsonString), &site)
	r.Site = &site

	r.Vars = controller.NewVars(&site)
	r.Template = LoadTemplates()
	//render := controller.NewRender(r.Template, r.Vars, &site)
	r.Paths["models"] = r.ModelsIndex
	//r.Paths["models"] = controller.NewModelsController(render)
	//r.Paths["sessions"] = controller.NewSessionsController(render)
	//for _, model := range r.Site.Models {
	//	r.Paths[fmt.Sprintf("/admin/%s", util.Plural(model.Name))] = "GET"
	//}

	return &r
}
