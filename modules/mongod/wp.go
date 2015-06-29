package mongod

import (
	"github.com/mongodbinc-interns/mongoproxy/messages"
	"gopkg.in/mgo.v2/bson"
)

func commandToBSONDoc(c messages.Command) bson.D {
	nameArg, ok := c.Args[c.CommandName]
	if !ok {
		nameArg = 1
	}
	args := bson.D{
		{c.CommandName, nameArg},
	}

	for arg, value := range c.Args {
		if arg != c.CommandName {
			args = append(args, bson.DocElem{arg, value})
		}
	}

	return args
}
