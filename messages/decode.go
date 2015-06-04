package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mongodbinc-interns/mongoproxy/buffer"
	"github.com/mongodbinc-interns/mongoproxy/convert"
	. "github.com/mongodbinc-interns/mongoproxy/log"
	"gopkg.in/mgo.v2/bson"
	"io"
	"strings"
)

func splitCommandOpQuery(q bson.D) (string, bson.M) {
	commandName := q[0].Name

	args := bson.M{}

	// throw the command arguments into args. This includes the command
	// name, as some of the commands have an important argument attached
	// to the command definition as well.
	for i := 0; i < len(q); i++ {
		arg := q[i].Name
		argV := q[i].Value

		args[arg] = argV
	}

	return commandName, args
}

// parseNamespace splits a namespace string into the database and collection.
// the first return value is the database, the second the collection. An error
// is returned if either the database or the collection doesn't exist.
func parseNamespace(namespace string) (string, string, error) {
	index := strings.Index(namespace, ".")
	if index < 0 || index >= len(namespace) {
		return "", "", fmt.Errorf("not a namespace")
	}
	database, collection := namespace[0:index], namespace[index+1:]

	// Error if empty database or collection
	if len(database) == 0 {
		return "", "", fmt.Errorf("empty database field")
	}

	if len(collection) == 0 {
		return "", "", fmt.Errorf("empty collection field")
	}

	return database, collection, nil
}

func processCommand(commandName string, database string, args bson.M) Command {
	c := Command{
		CommandName: commandName,
		Database:    database,
		Args:        args,
		Docs:        make([]bson.D, 0),
	}
	return c
}

func processFind(database string, args bson.M) (Find, error) {

	c := args["find"]
	collection, ok := c.(string)
	if !ok {
		// we have issues
		return Find{}, fmt.Errorf("Find command has no collection.")
	}

	f := Find{
		Database:        database,
		Collection:      collection,
		Filter:          convert.ToBSONDoc(args["filter"]),
		Projection:      convert.ToBSONDoc(args["projection"]),
		Skip:            convert.ToInt32(args["skip"]),
		Limit:           convert.ToInt32(args["limit"]),
		Tailable:        convert.ToBool(args["tailable"]),
		OplogReplay:     convert.ToBool(args["oplogReplay"]),
		NoCursorTimeout: convert.ToBool(args["noCursorTimeout"]),
		AwaitData:       convert.ToBool(args["awaitData"]),
		Partial:         convert.ToBool(args["partial"]),
	}

	return f, nil
}

func processInsert(database string, args bson.M) (Insert, error) {
	c := args["insert"]
	collection, ok := c.(string)
	if !ok {
		// we have issues
		return Insert{}, fmt.Errorf("Insert command has no collection.")
	}
	docs := args["documents"]
	Log(ERROR, "%#v", docs)
	documents, ok := docs.([]bson.D)
	if !ok {
		return Insert{}, fmt.Errorf("Insert command has no documents.")
	}
	insert := Insert{
		Database:   database,
		Collection: collection,
		Documents:  documents,
		Ordered:    convert.ToBool(args["ordered"], true),
	}

	return insert, nil
}

func processDelete(database string, args bson.M) (Delete, error) {
	c := args["delete"]
	collection, ok := c.(string)
	if !ok {
		// can't go on without a collection
		return Delete{}, fmt.Errorf("Delete command has no collection.")
	}

	argsDeletesRaw := args["deletes"]
	argsDeletes, ok := argsDeletesRaw.([]bson.M)
	if !ok {
		return Delete{}, fmt.Errorf("Delete command has no deletes.")
	}

	deletes := make([]SingleDelete, 0)
	for i := 0; i < len(argsDeletes); i++ {
		d := argsDeletes[i]
		singleDelete := SingleDelete{
			Selector: convert.ToBSONDoc(d["q"]),
			Limit:    convert.ToInt32(d["limit"]),
		}

		deletes = append(deletes, singleDelete)
	}

	delObj := Delete{
		Database:   database,
		Collection: collection,
		Deletes:    deletes,
		Ordered:    convert.ToBool(args["ordered"], true),
	}

	return delObj, nil
}

func processUpdate(database string, args bson.M) (Update, error) {

	c := args["update"]
	collection, ok := c.(string)
	if !ok {
		// we have issues
		return Update{}, fmt.Errorf("Update command has no collection.")
	}
	updates := make([]SingleUpdate, 0)

	argsUpdatesRaw := args["updates"]
	argsUpdates, ok := argsUpdatesRaw.([]bson.M)
	if !ok {
		return Update{}, fmt.Errorf("Update command has no updates.")
	}

	for i := 0; i < len(argsUpdates); i++ {
		u := argsUpdates[i]
		singleUpdate := SingleUpdate{
			Selector: convert.ToBSONDoc(u["q"]),
			Update:   convert.ToBSONDoc(u["u"]),
			Upsert:   convert.ToBool(u["upsert"]),
			Multi:    convert.ToBool(u["multi"]),
		}

		updates = append(updates, singleUpdate)
	}

	update := Update{
		Database:   database,
		Collection: collection,
		Updates:    updates,
		Ordered:    convert.ToBool(args["ordered"]),
	}

	return update, nil
}

func processGetMore(database string, args bson.M) (GetMore, error) {
	c := args["collection"]
	collection, ok := c.(string)
	if !ok {
		// we have issues
		return GetMore{}, fmt.Errorf("GetMore command has no collection.")
	}
	g := GetMore{
		Database:   database,
		Collection: collection,
		BatchSize:  convert.ToInt32(args["batchSize"]),
		CursorID:   convert.ToInt64(args["getMore"]),
	}

	return g, nil

}

// reads a header from the reader (16 bytes), consistent with wire protocol
func processHeader(reader io.Reader) (MsgHeader, error) {
	// read the message header
	msgHeaderBytes := make([]byte, 16)
	n, err := reader.Read(msgHeaderBytes)
	if err != nil && err != io.EOF {
		return MsgHeader{}, err
	}
	if n == 0 {
		// EOF?
		Log(INFO, "connection closed")
		return MsgHeader{}, err
	}
	mHeader := MsgHeader{}
	err = binary.Read(bytes.NewReader(msgHeaderBytes), binary.LittleEndian, &mHeader)
	if err != nil {
		Log(ERROR, "error decoding from reader: %v\n", err)
		return MsgHeader{}, err
	}
	Log(DEBUG, "request: %#v\n", mHeader)

	// sanity check
	if mHeader.MessageLength <= 15 {
		return MsgHeader{}, fmt.Errorf("Message length not long enough for header")
	}

	return mHeader, nil
}

// anything with OpCode 2004 goes here
func processOpQuery(reader io.Reader, header MsgHeader) (Requester, error) {
	// FLAGS
	flags, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return Command{}, fmt.Errorf("error reading flags: %v\n", err)
	}

	// DATABASE AND COLLECTIONS
	maxStringBytes := header.MessageLength - 16 - 4
	numNamespaceBytes, namespace, err := buffer.ReadNullTerminatedString(reader, maxStringBytes)
	if err != nil {
		return Command{}, fmt.Errorf("error reading null terminated string: %v\n", err)
	}
	database, collection, err := parseNamespace(namespace)

	if err != nil {
		return Command{}, fmt.Errorf("error parsing namespace: %v\n", err)
	}

	// NUMBER TO SKIP
	skip, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return Command{}, fmt.Errorf("error reading NumberToSkip: %v\n", err)
	}

	// NUMBER TO RETURN
	limit, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return Command{}, fmt.Errorf("error reading NumberToReturn: %v\n", err)
	}

	// QUERY
	var docSize int32
	docSize, q, err := buffer.ReadDocument(reader)
	if err != nil {
		return Command{}, fmt.Errorf("error reading query: %v\n", err)
	}
	totalBytesRead := 16 + 4 + numNamespaceBytes + 4 + 4 + docSize

	// PROJECTION
	var projection bson.D
	projection = nil
	if totalBytesRead < header.MessageLength {
		_, projection, err = buffer.ReadDocument(reader)
		if err != nil {
			if err != io.EOF {
				return Find{}, fmt.Errorf("error reading projection: %v\n", err)
			}
		}
	}

	// figure out what kind of struct to actually produce
	switch collection {
	case "$cmd":
		cName, args := splitCommandOpQuery(q)
		var c Requester
		switch cName {
		case "insert":
			// convert documents to an array of bson.D so that the struct
			// knows what to do with them.
			inserts, ok := args["documents"].([]interface{})
			if ok {
				d := make([]bson.D, len(inserts))
				for i := 0; i < len(inserts); i++ {
					doc := inserts[i]
					docD, ok2 := doc.(bson.D)
					if !ok2 {
						docD = bson.D{}
					}

					d[i] = docD
				}

				args["documents"] = d
			}
			c, err = processInsert(database, args)
			if err != nil {
				return nil, err
			}
			break
		case "update":
			// convert updates to an array of bson.M so that the struct
			// knows what to do with them.
			updates, ok := args["updates"].([]interface{})
			if ok {
				u := make([]bson.M, len(updates))
				for i := 0; i < len(updates); i++ {
					doc := updates[i]
					docD, ok2 := doc.(bson.D)
					if !ok2 {
						docD = bson.D{}
					}

					u[i] = docD.Map()
				}

				args["updates"] = u
			}

			c, err = processUpdate(database, args)
			if err != nil {
				return nil, err
			}
			break
		case "delete":
			c, err = processDelete(database, args)
			if err != nil {
				return nil, err
			}
			break
		default:
			c = processCommand(cName, database, args)
		}

		return c, nil
	default:
		// find command
		args := bson.M{}

		// this is to more closely match the command spec
		args["find"] = collection

		args["tailable"] = convert.ReadBit32LE(flags, 1)
		args["slaveOk"] = convert.ReadBit32LE(flags, 2)
		args["oplogReplay"] = convert.ReadBit32LE(flags, 3)
		args["noCursorTimeout"] = convert.ReadBit32LE(flags, 4)
		args["awaitData"] = convert.ReadBit32LE(flags, 5)
		args["partial"] = convert.ReadBit32LE(flags, 7)

		args["skip"] = skip
		args["limit"] = limit

		// the actual query
		args["filter"] = q
		if projection != nil {
			args["projection"] = projection
		}

		f, err := processFind(database, args)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

}

// OpCode 2001
func processOpUpdate(reader io.Reader, header MsgHeader) (Requester, error) {
	buffer.ReadInt32LE(reader) // the zero

	// DATABASE AND COLLECTIONS
	maxStringBytes := header.MessageLength - 16 - 4
	_, namespace, err := buffer.ReadNullTerminatedString(reader, maxStringBytes)
	if err != nil {
		return Update{}, fmt.Errorf("error reading null terminated string: %v\n", err)
	}

	database, collection, err := parseNamespace(namespace)
	if err != nil {
		return Update{}, fmt.Errorf("error parsing namespace: %v\n", err)
	}

	// FLAGS
	flags, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return Update{}, fmt.Errorf("error reading flags: %v\n", err)
	}

	// SELECTOR
	_, selector, err := buffer.ReadDocument(reader)
	if err != nil {
		return Update{}, fmt.Errorf("error reading selector: %v\n", err)
	}

	// UPDATE
	_, updator, err := buffer.ReadDocument(reader)
	if err != nil {
		return Update{}, fmt.Errorf("error reading selector: %v\n", err)
	}

	// create a proper update command
	args := bson.M{}
	updateObj := bson.M{}
	updateObj["q"] = selector
	updateObj["u"] = updator
	updateObj["upsert"] = convert.ReadBit32LE(flags, 0)
	updateObj["multi"] = convert.ReadBit32LE(flags, 1)
	args["update"] = collection
	updates := make([]bson.M, 0)
	updates = append(updates, updateObj)
	args["updates"] = updates

	return processUpdate(database, args)
}

// OpCode 2002
func processOpInsert(reader io.Reader, header MsgHeader) (Requester, error) {
	flags, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return Insert{}, fmt.Errorf("error reading flags: %v\n", err)
	}

	// DATABASE AND COLLECTIONS
	maxStringBytes := header.MessageLength - 16 - 4
	numNamespaceBytes, namespace, err := buffer.ReadNullTerminatedString(reader, maxStringBytes)
	if err != nil {
		return Insert{}, fmt.Errorf("error reading null terminated string: %v\n", err)
	}

	database, collection, err := parseNamespace(namespace)
	if err != nil {
		return Insert{}, fmt.Errorf("error parsing namespace: %v\n", err)
	}

	// DOCUMENTS
	totalBytesRead := 16 + 4 + numNamespaceBytes
	docs := make([]bson.D, 0)
	for totalBytesRead < header.MessageLength {
		n, doc, err := buffer.ReadDocument(reader)
		if err != nil {
			if err != io.EOF {
				return Insert{}, fmt.Errorf("error reading document: %v\n", err)
			}
		}
		docs = append(docs, doc)
		totalBytesRead += n
	}

	args := bson.M{}
	args["insert"] = collection
	args["ordered"] = !convert.ReadBit32LE(flags, 0)
	args["documents"] = docs

	return processInsert(database, args)

}

// OpCode 2005
func processOpGetMore(reader io.Reader, header MsgHeader) (Requester, error) {
	buffer.ReadInt32LE(reader) // the zero

	// DATABASE AND COLLECTIONS
	maxStringBytes := header.MessageLength - 16 - 4
	_, namespace, err := buffer.ReadNullTerminatedString(reader, maxStringBytes)
	if err != nil {
		return GetMore{}, fmt.Errorf("error reading null terminated string: %v\n", err)
	}

	database, collection, err := parseNamespace(namespace)
	if err != nil {
		return GetMore{}, fmt.Errorf("error parsing namespace: %v\n", err)
	}

	// NUM TO RETURN
	batchSize, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return GetMore{}, fmt.Errorf("error parsing batch size: %v\n", err)
	}

	// CURSOR ID
	cursorID, err := buffer.ReadInt64LE(reader)
	if err != nil {
		return GetMore{}, fmt.Errorf("error parsing cursor ID: %v\n", err)
	}

	args := bson.M{}
	args["getMore"] = cursorID
	args["collection"] = collection
	args["batchSize"] = batchSize

	return processGetMore(database, args)
}

// OpCode 2006
func processOpDelete(reader io.Reader, header MsgHeader) (Requester, error) {
	buffer.ReadInt32LE(reader) // the zero

	// DATABASE AND COLLECTIONS
	maxStringBytes := header.MessageLength - 16 - 4
	_, namespace, err := buffer.ReadNullTerminatedString(reader, maxStringBytes)
	if err != nil {
		return Delete{}, fmt.Errorf("error reading null terminated string: %v\n", err)
	}

	database, collection, err := parseNamespace(namespace)
	if err != nil {
		return Delete{}, fmt.Errorf("error parsing namespace: %v\n", err)
	}

	// FLAGS
	flags, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return Delete{}, fmt.Errorf("error reading flags: %v\n", err)
	}

	// SELECTOR
	_, selector, err := buffer.ReadDocument(reader)
	if err != nil {
		return Delete{}, fmt.Errorf("error reading selector: %v\n", err)
	}

	args := bson.M{}
	args["delete"] = collection
	deletes := make([]bson.M, 1)
	delObj := bson.M{}
	delObj["q"] = selector

	if convert.ReadBit32LE(flags, 0) {
		delObj["limit"] = 1
	} else {
		delObj["limit"] = 0
	}

	deletes[0] = delObj
	args["deletes"] = deletes

	return processDelete(database, args)
}

// Decodes a wire protocol message from a connection into a Requester to pass
// onto modules, a struct containing the header of the original message, and an error.
// It returns a non-nil error if reading from the connection
// fails in any way
func Decode(reader io.Reader) (Requester, MsgHeader, error) {
	mHeader, err := processHeader(reader)
	if err != nil {
		return nil, MsgHeader{}, err
	}

	switch mHeader.OpCode {
	case 2001:
		opu, err := processOpUpdate(reader, mHeader)
		if err != nil {
			return nil, MsgHeader{}, err
		}
		return opu, mHeader, nil
	case 2002:
		opi, err := processOpInsert(reader, mHeader)
		if err != nil {
			return nil, MsgHeader{}, err
		}
		return opi, mHeader, nil
	case 2004:
		opq, err := processOpQuery(reader, mHeader)
		if err != nil {
			return nil, MsgHeader{}, err
		}
		return opq, mHeader, nil
	case 2005:
		opg, err := processOpGetMore(reader, mHeader)
		if err != nil {
			return nil, MsgHeader{}, err
		}
		return opg, mHeader, nil
	case 2006:
		opd, err := processOpDelete(reader, mHeader)
		if err != nil {
			return nil, MsgHeader{}, err
		}
		return opd, mHeader, nil
	default:
		return nil, MsgHeader{}, fmt.Errorf("unimplemented operation: %#v", mHeader)
	}
}
