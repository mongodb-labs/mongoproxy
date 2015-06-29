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

func insertToBSONDoc(i messages.Insert) bson.D {
	args := bson.D{
		{"insert", i.Collection},
		{"documents", i.Documents},
		{"ordered", i.Ordered},
	}

	return args
}

func updateToBSONDoc(u messages.Update) bson.D {

	updates := make([]bson.M, len(u.Updates))

	for i := 0; i < len(u.Updates); i++ {
		singleUpdate := u.Updates[i]
		updates[i] = bson.M{
			"q":      singleUpdate.Selector,
			"u":      singleUpdate.Update,
			"upsert": singleUpdate.Upsert,
			"multi":  singleUpdate.Multi,
		}
	}

	args := bson.D{
		{"update", u.Collection},
		{"updates", updates},
		{"ordered", u.Ordered},
	}

	return args
}
