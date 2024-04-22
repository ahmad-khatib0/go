package auth

import (
	dt "github.com/ahmad-khatib0/go/websockets/chat/internal/db/types"
	st "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type AuthArgs struct {
	DB     dt.Adapter
	Conf   any
	Logger *logger.Logger
	UGen   *st.UidGenerator
}
