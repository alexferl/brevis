package model

import (
	"github.com/marksalpeter/token"
)

type Token struct {
	Id token.Token
}

func NewToken() string {
	t := Token{
		Id: token.New(),
	}
	return t.Id.Encode()
}
