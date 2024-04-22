package fcm

import "github.com/ahmad-khatib0/go/websockets/chat/internal/push/types"

type Fcm struct{}

func NewFcm() types.Handler {
	return Fcm{}
}
