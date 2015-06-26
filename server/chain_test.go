package server

import (
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/messages"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

var msgOne = bson.M{"ok": 1}
var msgTwo = bson.M{"message": "two"}

// A MockReq emulates a messages.Requester interface. For testing only.
type MockReq struct {
}

func (f MockReq) Type() string {
	return "request"
}

type MockRes struct {
	Data []bson.M
}

// A MockRes emulates a messages.Responder interface, simply pushing
// written documents into a slice that can be examined. For testing only.
func (f *MockRes) Type() string {
	return "response"
}

func (f *MockRes) Write(res messages.ResponseWriter) {
	f.Data = append(f.Data, res.ToBSON())
}

func (f *MockRes) Error(code int32, message string) {
	fmt.Println("We got an error!")
}

type ModuleOne struct {
}

func (m ModuleOne) Process(req messages.Requester, res messages.Responder, next PipelineFunc) {
	r := messages.CommandResponse{}
	r.Reply = msgOne
	res.Write(r)
	next(req, res)
}

type ModuleTwo struct {
}

func (m ModuleTwo) Process(req messages.Requester, res messages.Responder, next PipelineFunc) {
	next(req, res)
	r := messages.CommandResponse{}
	r.Reply = msgTwo
	res.Write(r)
}

func TestModuleChaining(t *testing.T) {
	Convey("Create a chain", t, func() {

		Convey("with no modules", func() {

			chain := CreateChain()
			r := MockReq{}
			w := &MockRes{
				Data: make([]bson.M, 0),
			}
			pipeline := BuildPipeline(chain)
			pipeline(r, w)

			So(len(w.Data), ShouldEqual, 0)

		})
		Convey("with a single module", func() {

			m1 := ModuleOne{}
			chain := CreateChain()
			chain.AddModule(m1)
			r := MockReq{}
			w := &MockRes{
				Data: make([]bson.M, 0),
			}
			pipeline := BuildPipeline(chain)
			pipeline(r, w)

			So(w.Data[0], ShouldEqual, msgOne)

		})
		Convey("with multiple modules", func() {
			Convey("in the correct order", func() {
				m1 := ModuleOne{}
				m2 := ModuleTwo{}
				chain := CreateChain()
				chain.AddModule(m2)
				chain.AddModule(m1)
				r := MockReq{}
				w := &MockRes{
					Data: make([]bson.M, 0),
				}
				pipeline := BuildPipeline(chain)
				pipeline(r, w)

				So(w.Data[0], ShouldEqual, msgOne)
				So(w.Data[1], ShouldEqual, msgTwo)
			})
		})

	})
}
