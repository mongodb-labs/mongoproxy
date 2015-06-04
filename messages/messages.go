// Package messages contains structs and functions to encode and decode
// wire protocol messages.
package messages

import (
	"gopkg.in/mgo.v2/bson"
)

// a struct to represent a wire protocol message header.
type MsgHeader struct {
	MessageLength int32
	RequestID     int32
	ResponseTo    int32
	OpCode        int32
}

// struct for a generic command, the default Requester sent from proxy
// core to modules
type Command struct {
	CommandName string
	Database    string
	Args        bson.M
	Metadata    bson.M
	Docs        []bson.D
}

func (c Command) Type() string {
	return "command"
}

// GetArg takes the name of an argument for the command and returns the
// value of that argument.
func (c Command) GetArg(arg string) interface{} {
	a, ok := c.Args[arg]
	if !ok {
		return nil
	}
	return a
}

// the struct for the 'find' command.
type Find struct {
	Database        string
	Collection      string
	Filter          bson.D
	Sort            bson.D
	Projection      bson.D
	Skip            int32
	Limit           int32
	Tailable        bool
	OplogReplay     bool
	NoCursorTimeout bool
	AwaitData       bool
	Partial         bool
}

func (f Find) Type() string {
	return "find"
}

// the struct for the 'insert' command
type Insert struct {
	Database   string
	Collection string
	Documents  []bson.D
	Ordered    bool
}

func (i Insert) Type() string {
	return "insert"
}

type SingleUpdate struct {
	Selector bson.D
	Update   bson.D
	Upsert   bool
	Multi    bool
}

// the struct for the 'update' command
type Update struct {
	Database   string
	Collection string
	Updates    []SingleUpdate
	Ordered    bool
}

func (u Update) Type() string {
	return "update"
}

type SingleDelete struct {
	Selector bson.D
	Limit    int32
}

// struct for 'delete' command
type Delete struct {
	Database   string
	Collection string
	Deletes    []SingleDelete
	Ordered    bool
}

func (d Delete) Type() string {
	return "delete"
}

// struct for 'getMore' command
type GetMore struct {
	Database   string
	CursorID   int64
	Collection string
	BatchSize  int32
}

func (g GetMore) Type() string {
	return "getMore"
}
