// Package messages contains structs and functions to encode and decode
// wire protocol messages.
package messages

import (
	"gopkg.in/mgo.v2/bson"
)

// constants representing the different opcodes for the wire protocol.
const (
	OP_UPDATE   int32 = 2001
	OP_INSERT         = 2002
	OP_QUERY          = 2004
	OP_GET_MORE       = 2005
	OP_DELETE         = 2006
)

// constants representing the types of request structs supported by proxy core.
const (
	CommandType string = "command"
	FindType           = "find"
	InsertType         = "insert"
	UpdateType         = "update"
	DeleteType         = "delete"
	GetMoreType        = "getMore"
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
	RequestID   int32
	CommandName string
	Database    string
	Args        bson.M
	Metadata    bson.M
	Docs        []bson.D
}

func (c Command) Type() string {
	return CommandType
}

func (c Command) ID() int32 {
	return c.RequestID
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
	RequestID       int32
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
	return FindType
}

func (f Find) ID() int32 {
	return f.RequestID
}

// the struct for the 'insert' command
type Insert struct {
	RequestID  int32
	Database   string
	Collection string
	Documents  []bson.D
	Ordered    bool
}

func (i Insert) Type() string {
	return InsertType
}

func (i Insert) ID() int32 {
	return i.RequestID
}

type SingleUpdate struct {
	Selector bson.D
	Update   bson.D
	Upsert   bool
	Multi    bool
}

// the struct for the 'update' command
type Update struct {
	RequestID  int32
	Database   string
	Collection string
	Updates    []SingleUpdate
	Ordered    bool
}

func (u Update) Type() string {
	return UpdateType
}

func (u Update) ID() int32 {
	return u.RequestID
}

type SingleDelete struct {
	Selector bson.D
	Limit    int32
}

// struct for 'delete' command
type Delete struct {
	RequestID  int32
	Database   string
	Collection string
	Deletes    []SingleDelete
	Ordered    bool
}

func (d Delete) Type() string {
	return DeleteType
}
func (d Delete) ID() int32 {
	return d.RequestID
}

// struct for 'getMore' command
type GetMore struct {
	RequestID  int32
	Database   string
	CursorID   int64
	Collection string
	BatchSize  int32
}

func (g GetMore) Type() string {
	return GetMoreType
}

func (g GetMore) ID() int32 {
	return g.RequestID
}
