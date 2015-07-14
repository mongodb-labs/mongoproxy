package frontend

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi/frontend/controllers"
	"html/template"
)

// the funcMap adds helper functions that can be injected into the HTML
// templates
var funcMap = template.FuncMap{
	// marshal converts an interface{} into a JSON object that can be consumed
	// by the browser.
	"marshal": func(v interface{}) template.JS {
		fmt.Printf("%#v", v)
		a, err := json.Marshal(v)
		fmt.Printf("%#v", err)
		return template.JS(a)
	},
}

// Start initializes a new server for the BI Module Frontend.
func Start(module *bi.BIModule, baseDir string) *gin.Engine {
	r := gin.New()

	// set up views
	// It's a bit of a hack right now, since we need the function map
	// https://github.com/gin-gonic/gin/issues/228
	if tmpl, err := template.New("projectViews").Funcs(funcMap).ParseGlob(baseDir + "/views/*"); err == nil {
		r.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}

	// set up routes
	controllers.Setup(r, module, baseDir)

	return r
}
