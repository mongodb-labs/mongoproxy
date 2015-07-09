package messages

import (
	"fmt"
)

func ToFindRequest(r Requester) (Find, error) {
	f, ok := r.(Find)
	if !ok {
		return Find{}, fmt.Errorf("Requester was not a find object. Requester received instead: %#v", r)
	}
	return f, nil
}

func ToGetMoreRequest(r Requester) (GetMore, error) {
	g, ok := r.(GetMore)
	if !ok {
		return GetMore{}, fmt.Errorf("Requester was not a getMore object. Requester received instead: %#v", r)
	}
	return g, nil
}

func ToInsertRequest(r Requester) (Insert, error) {
	i, ok := r.(Insert)
	if !ok {
		return Insert{}, fmt.Errorf("Requester was not an insert object. Requester received instead: %#v", r)
	}
	return i, nil
}

func ToUpdateRequest(r Requester) (Update, error) {
	u, ok := r.(Update)
	if !ok {
		return Update{}, fmt.Errorf("Requester was not an update object. Requester received instead: %#v", r)
	}
	return u, nil
}

func ToDeleteRequest(r Requester) (Delete, error) {
	f, ok := r.(Delete)
	if !ok {
		return Delete{}, fmt.Errorf("Requester was not a delete object. Requester received instead: %#v", r)
	}
	return f, nil
}

func ToCommandRequest(r Requester) (Command, error) {
	c, ok := r.(Command)
	if !ok {
		return Command{}, fmt.Errorf("Requester was not a command object. Requester received instead: %#v", r)
	}
	return c, nil
}
