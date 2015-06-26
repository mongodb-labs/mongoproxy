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

// convertToBSONMapSlice converts an []interface{}, []bson.D, or []bson.M slice to a []bson.M
// slice (assuming that all contents are either bson.M or bson.D objects)
func convertToBSONMapSlice(input interface{}) ([]bson.M, error) {

	inputBSONM, ok := input.([]bson.M)
	if ok {
		return inputBSONM, nil
	}

	inputBSOND, ok := input.([]bson.D)
	if ok {
		// just convert all of the bson.D documents to bson.M
		d := make([]bson.M, len(inputBSOND))
		for i := 0; i < len(inputBSOND); i++ {
			doc := inputBSOND[i]
			d[i] = doc.Map()
		}
		return d, nil
	}

	inputInterface, ok := input.([]interface{})
	if ok {
		d := make([]bson.M, len(inputInterface))
		for i := 0; i < len(inputInterface); i++ {
			doc := inputInterface[i]
			docM, ok2 := doc.(bson.M)
			if !ok2 {
				// check if it's a bson.D
				docD, ok3 := doc.(bson.D)
				if ok3 {
					docM = docD.Map()
				} else {
					// error
					return nil, fmt.Errorf("Slice contents aren't BSON objects")
				}
			}

			d[i] = docM
		}
		return d, nil
	}

	return nil, fmt.Errorf("Unsupported input")
}

// convertToBSONDocSlice converts an []interface{} to a []bson.D slice
// assuming contents are bson.D objects
func convertToBSONDocSlice(input interface{}) ([]bson.D, error) {
	inputBSOND, ok := input.([]bson.D)
	if ok {
		return inputBSOND, nil
	}

	inputInterface, ok := input.([]interface{})
	if ok {
		d := make([]bson.D, len(inputInterface))
		for i := 0; i < len(inputInterface); i++ {
			doc := inputInterface[i]
			docD, ok2 := doc.(bson.D)
			if !ok2 {
				return nil, fmt.Errorf("Slice contents aren't BSON objects")
			}
			d[i] = docD
		}
		return d, nil
	}

	return nil, fmt.Errorf("Unsupported input")
}

func splitCommandOpQuery(q bson.D) (string, bson.M) {
	commandName := q[0].Name

	args := bson.M{}

	// throw the command arguments into args. This includes the command
	// name, as some of the commands have an important argument attached
	// to the command definition as well.
	for i := 0; i < len(q); i++ {
		args[q[i].Name] = q[i].Value
	}

	return commandName, args
}

// parseNamespace splits a namespace string into the database and collection.
// The first return value is the database, the second, the collection. An error
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

func createCommand(header MsgHeader, commandName string, database string, args bson.M) Command {
	c := Command{
		RequestID:   header.RequestID,
		CommandName: commandName,
		Database:    database,
		Args:        args,
	}
	return c
}

func createFind(header MsgHeader, database string, args bson.M) (Find, error) {

	c := args["find"]
	collection, ok := c.(string)
	if !ok {
		// we have issues
		return Find{}, fmt.Errorf("Find command has no collection.")
	}

	f := Find{
		RequestID:       header.RequestID,
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

func createInsert(header MsgHeader, database string, args bson.M) (Insert, error) {
	c := args["insert"]
	collection, ok := c.(string)
	if !ok {
		// we have issues
		return Insert{}, fmt.Errorf("Insert command has no collection.")
	}
	docs := args["documents"]
	documents, ok := docs.([]bson.D)
	if !ok {
		return Insert{}, fmt.Errorf("Insert command has no documents.")
	}
	insert := Insert{
		RequestID:  header.RequestID,
		Database:   database,
		Collection: collection,
		Documents:  documents,
		Ordered:    convert.ToBool(args["ordered"], true),
	}

	return insert, nil
}

func createDelete(header MsgHeader, database string, args bson.M) (Delete, error) {
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

	deletes := make([]SingleDelete, len(argsDeletes))
	for i := 0; i < len(argsDeletes); i++ {
		d := argsDeletes[i]
		singleDelete := SingleDelete{
			Selector: convert.ToBSONDoc(d["q"]),
			Limit:    convert.ToInt32(d["limit"]),
		}

		deletes[i] = singleDelete
	}

	delObj := Delete{
		RequestID:  header.RequestID,
		Database:   database,
		Collection: collection,
		Deletes:    deletes,
		Ordered:    convert.ToBool(args["ordered"], true),
	}

	return delObj, nil
}

func createUpdate(header MsgHeader, database string, args bson.M) (Update, error) {

	c := args["update"]
	collection, ok := c.(string)
	if !ok {
		// we have issues
		return Update{}, fmt.Errorf("Update command has no collection.")
	}

	argsUpdatesRaw := args["updates"]
	argsUpdates, ok := argsUpdatesRaw.([]bson.M)
	if !ok {
		return Update{}, fmt.Errorf("Update command has no updates.")
	}

	updates := make([]SingleUpdate, len(argsUpdates))
	for i := 0; i < len(argsUpdates); i++ {
		u := argsUpdates[i]
		singleUpdate := SingleUpdate{
			Selector: convert.ToBSONDoc(u["q"]),
			Update:   convert.ToBSONDoc(u["u"]),
			Upsert:   convert.ToBool(u["upsert"]),
			Multi:    convert.ToBool(u["multi"]),
		}

		updates[i] = singleUpdate
	}

	update := Update{
		RequestID:  header.RequestID,
		Database:   database,
		Collection: collection,
		Updates:    updates,
		Ordered:    convert.ToBool(args["ordered"]),
	}

	return update, nil
}

func createGetMore(header MsgHeader, database string, args bson.M) (GetMore, error) {

	c := args["collection"]
	collection, ok := c.(string)
	if !ok {
		// we have issues
		return GetMore{}, fmt.Errorf("GetMore command has no collection.")
	}
	g := GetMore{
		RequestID:  header.RequestID,
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
	// flags
	flags, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading flags: %v", err)
	}

	// namespace

	// cut off the string at the remaining message length in case it is not
	// null terminated.
	maxStringBytes := header.MessageLength - // length of the whole wire protocol message
		16 - // bytes representing the header
		4 // bytes representing the flags
	numNamespaceBytes, namespace, err := buffer.ReadNullTerminatedString(reader, maxStringBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading null terminated string: %v", err)
	}
	database, collection, err := parseNamespace(namespace)

	if err != nil {
		return nil, fmt.Errorf("error parsing namespace: %v", err)
	}

	// numberToSkip
	skip, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading NumberToSkip: %v", err)
	}

	// numberToReturn
	limit, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading NumberToReturn: %v", err)
	}

	// query
	var docSize int32
	docSize, q, err := buffer.ReadDocument(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading query: %v", err)
	}
	totalBytesRead := 16 + // bytes representing the header
		4 + // bytes representing flags
		numNamespaceBytes + // bytes representing the namespace
		4 + // bytes representing numberToSkip (skip)
		4 + // bytes representing numberToReturn (limit)
		docSize // bytes representing the query

	// optional projection
	var projection bson.D
	if totalBytesRead < header.MessageLength {
		_, projection, err = buffer.ReadDocument(reader)
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("error reading projection: %v", err)
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
			i, err := convertToBSONDocSlice(args["documents"])

			if err != nil {
				return nil, err
			}

			args["documents"] = i

			c, err = createInsert(header, database, args)
			if err != nil {
				return nil, err
			}
			break
		case "update":
			// convert updates to an array of bson.M so that the struct
			// knows what to do with them.
			u, err := convertToBSONMapSlice(args["updates"])

			if err != nil {
				return nil, err
			}

			args["updates"] = u

			c, err = createUpdate(header, database, args)
			if err != nil {
				return nil, err
			}
			break
		case "delete":

			d, err := convertToBSONMapSlice(args["deletes"])
			if err != nil {
				return nil, err
			}

			args["deletes"] = d

			c, err = createDelete(header, database, args)
			if err != nil {
				return nil, err
			}
			break
		default:
			c = createCommand(header, cName, database, args)
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
		args["projection"] = projection

		f, err := createFind(header, database, args)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

}

// OpCode 2001
func processOpUpdate(reader io.Reader, header MsgHeader) (Requester, error) {
	buffer.ReadInt32LE(reader) // the zero (not used in wire protocol)

	// namespace

	// cut off the string at the remaining message length in case it is not
	// null terminated.
	maxStringBytes := header.MessageLength -
		16 - // bytes representing the header
		4 // the zero
	_, namespace, err := buffer.ReadNullTerminatedString(reader, maxStringBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading null terminated string: %v", err)
	}

	database, collection, err := parseNamespace(namespace)
	if err != nil {
		return nil, fmt.Errorf("error parsing namespace: %v", err)
	}

	// flags
	flags, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading flags: %v", err)
	}

	// selector
	_, selector, err := buffer.ReadDocument(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading selector: %v", err)
	}

	// update
	_, updator, err := buffer.ReadDocument(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading selector: %v", err)
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

	return createUpdate(header, database, args)
}

// OpCode 2002
func processOpInsert(reader io.Reader, header MsgHeader) (Requester, error) {
	flags, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading flags: %v", err)
	}

	// namespace

	// cut off the string at the remaining message length in case it is not
	// null terminated.
	maxStringBytes := header.MessageLength -
		16 - // bytes representing the header
		4 // bytes representing flags
	numNamespaceBytes, namespace, err := buffer.ReadNullTerminatedString(reader, maxStringBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading null terminated string: %v", err)
	}

	database, collection, err := parseNamespace(namespace)
	if err != nil {
		return nil, fmt.Errorf("error parsing namespace: %v", err)
	}

	// documents to insert
	totalBytesRead := 16 + 4 + numNamespaceBytes
	docs := make([]bson.D, 0)
	for totalBytesRead < header.MessageLength {
		n, doc, err := buffer.ReadDocument(reader)
		if err != nil {
			if err != io.EOF {
				return Insert{}, fmt.Errorf("error reading document: %v", err)
			}
		}
		docs = append(docs, doc)
		totalBytesRead += n
	}

	args := bson.M{}
	args["insert"] = collection
	args["ordered"] = !convert.ReadBit32LE(flags, 0)
	args["documents"] = docs

	return createInsert(header, database, args)

}

// OpCode 2005
func processOpGetMore(reader io.Reader, header MsgHeader) (Requester, error) {
	buffer.ReadInt32LE(reader) // the zero (not used in wire protocol)

	// namespace

	// cut off the string at the remaining message length in case it is not
	// null terminated.
	maxStringBytes := header.MessageLength -
		16 - // bytes representing the header
		4 // the zero
	_, namespace, err := buffer.ReadNullTerminatedString(reader, maxStringBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading null terminated string: %v", err)
	}

	database, collection, err := parseNamespace(namespace)
	if err != nil {
		return nil, fmt.Errorf("error parsing namespace: %v", err)
	}

	// numToReturn
	batchSize, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing batch size: %v", err)
	}

	// cursorID
	cursorID, err := buffer.ReadInt64LE(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing cursor ID: %v", err)
	}

	args := bson.M{}
	args["getMore"] = cursorID
	args["collection"] = collection
	args["batchSize"] = batchSize

	return createGetMore(header, database, args)
}

// OpCode 2006
func processOpDelete(reader io.Reader, header MsgHeader) (Requester, error) {
	buffer.ReadInt32LE(reader) // the zero (not used in wire protocol)

	// namespace

	// cut off the string at the remaining message length in case it is not
	// null terminated.
	maxStringBytes := header.MessageLength -
		16 - // bytes representing the header
		4 // the zero
	_, namespace, err := buffer.ReadNullTerminatedString(reader, maxStringBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading null terminated string: %v", err)
	}

	database, collection, err := parseNamespace(namespace)
	if err != nil {
		return nil, fmt.Errorf("error parsing namespace: %v", err)
	}

	// flags
	flags, err := buffer.ReadInt32LE(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading flags: %v", err)
	}

	// selector
	_, selector, err := buffer.ReadDocument(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading selector: %v", err)
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

	return createDelete(header, database, args)
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
	case OP_UPDATE:
		opu, err := processOpUpdate(reader, mHeader)
		if err != nil {
			return nil, MsgHeader{}, err
		}
		return opu, mHeader, nil
	case OP_INSERT:
		opi, err := processOpInsert(reader, mHeader)
		if err != nil {
			return nil, MsgHeader{}, err
		}
		return opi, mHeader, nil
	case OP_QUERY:
		opq, err := processOpQuery(reader, mHeader)
		if err != nil {
			return nil, MsgHeader{}, err
		}
		return opq, mHeader, nil
	case OP_GET_MORE:
		opg, err := processOpGetMore(reader, mHeader)
		if err != nil {
			return nil, MsgHeader{}, err
		}
		return opg, mHeader, nil
	case OP_DELETE:
		opd, err := processOpDelete(reader, mHeader)
		if err != nil {
			return nil, MsgHeader{}, err
		}
		return opd, mHeader, nil
	default:
		return nil, MsgHeader{}, fmt.Errorf("unimplemented operation: %#v", mHeader)
	}
}
