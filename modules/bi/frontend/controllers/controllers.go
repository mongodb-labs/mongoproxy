package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"github.com/mongodbinc-interns/mongoproxy/modules/bi"
	"gopkg.in/mgo.v2"
	"net/http"
	"strconv"
	"time"
)

var biModule bi.BIModule
var mongoSession *mgo.Session
var mongoDBDialInfo = &mgo.DialInfo{
	// TODO: Allow configurable connection info
	Addrs:    []string{"localhost:27017"},
	Timeout:  60 * time.Second,
	Database: "test",
}

type metricParam struct {
	RuleIndex   int64
	Granularity string
	Start       time.Time
	End         time.Time
}

func parseMetricParams(c *gin.Context) (*metricParam, error) {
	ruleIndex, err := strconv.ParseInt(c.Param("ruleIndex"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid ruleIndex: not a number")
	}
	if ruleIndex < 0 || ruleIndex >= int64(len(biModule.Rules)) {
		return nil, fmt.Errorf("Invalid ruleIndex: out of range")
	}

	granularity := c.Param("granularity")

	var start time.Time
	err = start.UnmarshalText([]byte(c.Param("start")))
	if err != nil {
		return nil, fmt.Errorf("Invalid start time: %v", c.Param("start"))
	}

	var end time.Time
	err = end.UnmarshalText([]byte(c.Param("end")))
	if err != nil {
		return nil, fmt.Errorf("Invalid end time: %v", c.Param("end"))
	}

	return &metricParam{
		ruleIndex, granularity, start, end,
	}, nil
}

func init() {
	var err error
	mongoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		Log(ERROR, "%#v\n", err)
		return
	}
}

func getMain(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", gin.H{
		"module": biModule,
	})
}

func getMetric(c *gin.Context) {
	params, err := parseMetricParams(c)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}

	docs, err := getDataOverRange(mongoSession, biModule.Rules[params.RuleIndex],
		params.Granularity, params.Start, params.End)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, docs)
}

func Setup(r *gin.Engine, source bi.BIModule) *gin.Engine {
	biModule = source

	r.GET("/", getMain)
	r.GET("/data/:ruleIndex/:granularity/:start/:end", getMetric)

	r.Static("/public", "./public")
	return r
}
