package server

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

func usersUpdateUnread(uid types.Uid, val int, inc bool) {
	if globals.usersUpdate == nil || (val == 0 && inc) {
		return
	}

	upd := &UserCacheReq{UserId: uid, Unread: val, Inc: inc}
	if globals.cluster.isRemoteTopic(uid.UserId()) {
		// Send request to remote node which owns the user.
		globals.cluster.routeUserReq(upd)
	} else {
		select {
		case globals.usersUpdate <- upd:
		default:
		}
	}
}
