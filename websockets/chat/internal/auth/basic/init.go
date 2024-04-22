package basic

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
)

func NewAuthBasic(aa auth.AuthArgs) types.AuthHandler {
	return &authenticator{}
}
