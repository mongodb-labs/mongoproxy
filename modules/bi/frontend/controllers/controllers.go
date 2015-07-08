package controllers

import (
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
	Log(INFO, "%v", c.Param("ruleIndex"))
	ruleIndex, err := strconv.ParseInt(c.Param("ruleIndex"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid rule: not a number",
			"err":     err,
		})
		return
	}
	if ruleIndex < 0 || ruleIndex >= int64(len(biModule.Rules)) {
		c.JSON(400, gin.H{
			"message":   "Invalid rule: out of range",
			"ruleIndex": ruleIndex,
		})
		return
	}

	granularity := c.Param("granularity")
	valueField := c.Param("valueField")

	var start time.Time
	err = start.UnmarshalText([]byte(c.Param("start")))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid start time",
			"start":   c.Param("start"),
			"err":     err,
		})
		return
	}

	var end time.Time
	err = end.UnmarshalText([]byte(c.Param("end")))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid end time",
			"err":     err,
		})
		return
	}

	docs, err := getDataOverRange(mongoSession, biModule.Rules[ruleIndex], granularity, valueField, start, end)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error retrieving data",
			"err":     err,
		})
		return
	}
	c.JSON(200, docs)
}

func Setup(r *gin.Engine, source bi.BIModule) *gin.Engine {
	biModule = source

	r.GET("/", getMain)
	r.GET("/data/:ruleIndex/:granularity/:valueField/:start/:end", getMetric)

	r.Static("/public", "./public")
	return r
}
