package types

import (
	"andrew_chat/intenal/domain"
)

type ServerMsg struct {
	Status  int // "connect", "disconnect", "select"
	Server  domain.Server
	Success bool
}

//The message is an instruction to the window manager where to place the window.
type CreateWindowMsg struct {
	Pos   Position
	Model ComponentModel
}

//The message is an instruction to the window manager to delete the window.
type DeleteWindowMsg struct {
	Model ComponentModel
}

type FilterMsg struct{
	Filter string
}

type ErrMsg struct{
	Text string
}

type InputFormMsg struct {
	Fields []InputFieldValue
}