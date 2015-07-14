package main

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
func Start(module *bi.BIModule) *gin.Engine {
	r := gin.New()

	// set up views
	// It's a bit of a hack right now, since we need the function map
	// https://github.com/gin-gonic/gin/issues/228
	if tmpl, err := template.New("projectViews").Funcs(funcMap).ParseGlob("views/*"); err == nil {
		r.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}

	// set up routes
	controllers.Setup(r, module)

	return r
}

// TODO: put in different file so that server.go can go into the 'frontend' package
func main() {
	// initialize BI module
	// TODO: load from a config, rather than hardcoding
	biModule := bi.BIModule{}

	t := make([]string, 2)
	t[0] = bi.Daily
	t[1] = bi.Minutely

	rule := bi.Rule{
		OriginDatabase:    "test",
		OriginCollection:  "foo",
		PrefixDatabase:    "db",
		PrefixCollection:  "metrics",
		TimeGranularities: t,
		ValueField:        "price",
	}

	t2 := make([]string, 5)
	t2[0] = bi.Monthly
	t2[1] = bi.Daily
	t2[2] = bi.Hourly
	t2[3] = bi.Minutely
	t2[4] = bi.Secondly

	rule2 := bi.Rule{
		OriginDatabase:    "test",
		OriginCollection:  "foo",
		PrefixDatabase:    "db",
		PrefixCollection:  "metrics",
		TimeGranularities: t2,
		ValueField:        "amount",
	}
	biModule.Rules = append(biModule.Rules, rule)
	biModule.Rules = append(biModule.Rules, rule2)

	r := Start(&biModule)
	r.Run(":8080")
}
