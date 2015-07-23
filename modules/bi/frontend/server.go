package frontend

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi/frontend/controllers"
	"gopkg.in/mgo.v2/bson"
	"html/template"
)

// the funcMap adds helper functions that can be injected into the HTML
// templates
var funcMap = template.FuncMap{
	// marshal converts an interface{} into a JSON object that can be consumed
	// by the browser.
	"marshal": func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	},
}

// Start initializes a new server for the BI Module Frontend.
func Start(config bson.M, baseDir string, configLocation *controllers.ConfigLocation) (*gin.Engine, error) {
	r := gin.New()

	// set up views
	// It's a bit of a hack right now, since we need the function map
	// https://github.com/gin-gonic/gin/issues/228
	if tmpl, err := template.New("projectViews").Funcs(funcMap).ParseGlob(baseDir + "/views/*"); err == nil {
		r.SetHTMLTemplate(tmpl)
	} else {
		return nil, err
	}

	// set up routes
	controllers.SetConfigSaveLocation(configLocation)
	err := controllers.Setup(r, config, baseDir)
	if err != nil {
		return nil, err
	}

	return r, nil
}
