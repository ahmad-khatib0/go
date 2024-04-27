package server

import (
	"net/http"
	"time"
)

// Generators of server-side error messages {ctrl}.

// NoErr indicates successful completion (200).
func NoErr(id, topic string, ts time.Time) *ServerComMessage {
	return NoErrParams(id, topic, ts, nil)
}

// NoErrExplicitTs indicates successful completion with explicit server and incoming request timestamps (200).
func NoErrExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return NoErrParamsExplicitTs(id, topic, serverTs, incomingReqTs, nil)
}

// NoErrReply indicates successful completion as a reply to a client message (200).
func NoErrReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return NoErrExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// NoErrParams indicates successful completion with additional parameters (200).
func NoErrParams(id, topic string, ts time.Time, params any) *ServerComMessage {
	return NoErrParamsExplicitTs(id, topic, ts, ts, params)
}

// NoErrParamsExplicitTs indicates successful completion with additional parameters
// and explicit server and incoming request timestamps (200).
func NoErrParamsExplicitTs(id, topic string, serverTs, incomingReqTs time.Time, params any) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusOK, // 200
			Text:      "ok",
			Topic:     topic,
			Params:    params,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// NoErrParamsReply indicates successful completion with additional parameters
// and explicit server and incoming request timestamps (200).
func NoErrParamsReply(msg *ClientComMessage, ts time.Time, params any) *ServerComMessage {
	return NoErrParamsExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp, params)
}

// NoErrCreated indicated successful creation of an object (201).
func NoErrCreated(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusCreated, // 201
			Text:      "created",
			Topic:     topic,
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// NoErrAccepted indicates request was accepted but not processed yet (202).
func NoErrAccepted(id, topic string, ts time.Time) *ServerComMessage {
	return NoErrAcceptedExplicitTs(id, topic, ts, ts)
}

// NoErrAcceptedExplicitTs indicates request was accepted but not processed yet
// with explicit server and incoming request timestamps (202).
func NoErrAcceptedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusAccepted, // 202
			Text:      "accepted",
			Topic:     topic,
			Timestamp: serverTs,
		}, Id: id,
		Timestamp: incomingReqTs,
	}
}

// NoContentParams indicates request was processed but resulted in no content (204).
func NoContentParams(id, topic string, serverTs, incomingReqTs time.Time, params any) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusNoContent, // 204
			Text:      "no content",
			Topic:     topic,
			Params:    params,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// NoContentParamsReply indicates request was processed but resulted in no content
// in response to a client request (204).
func NoContentParamsReply(msg *ClientComMessage, ts time.Time, params any) *ServerComMessage {
	return NoContentParams(msg.ID, msg.Original, ts, msg.Timestamp, params)
}

// NoErrEvicted indicates that the user was disconnected from topic for no fault of the user (205).
func NoErrEvicted(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusResetContent, // 205
			Text:      "evicted",
			Topic:     topic,
			Timestamp: ts,
		}, Id: id,
	}
}

// NoErrShutdown means user was disconnected from topic because system shutdown is in progress (205).
func NoErrShutdown(ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Code:      http.StatusResetContent, // 205
			Text:      "server shutdown",
			Timestamp: ts,
		},
	}
}

// NoErrDeliveredParams means requested content has been delivered (208).
func NoErrDeliveredParams(id, topic string, ts time.Time, params any) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusAlreadyReported, // 208
			Text:      "delivered",
			Topic:     topic,
			Params:    params,
			Timestamp: ts,
		},
		Id: id,
	}
}

// 3xx

// InfoValidateCredentials requires user to confirm credentials before going forward (300).
func InfoValidateCredentials(id string, ts time.Time) *ServerComMessage {
	return InfoValidateCredentialsExplicitTs(id, ts, ts)
}

// InfoValidateCredentialsExplicitTs requires user to confirm credentials before going forward
// with explicit server and incoming request timestamps (300).
func InfoValidateCredentialsExplicitTs(id string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusMultipleChoices, // 300
			Text:      "validate credentials",
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// InfoChallenge requires user to respond to presented challenge before login can be completed (300).
func InfoChallenge(id string, ts time.Time, challenge []byte) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusMultipleChoices, // 300
			Text:      "challenge",
			Params:    map[string]any{"challenge": challenge},
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// InfoAuthReset is sent in response to request to reset authentication when it was completed
// but login was not performed (301).
func InfoAuthReset(id string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusMovedPermanently, // 301
			Text:      "auth reset",
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// InfoUseOther is a response to a subscription request redirecting client to another topic (303).
func InfoUseOther(id, topic, other string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusSeeOther, // 303
			Text:      "use other",
			Topic:     topic,
			Params:    map[string]string{"topic": other},
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// InfoUseOtherReply is a response to a subscription request redirecting client to another topic (303).
func InfoUseOtherReply(msg *ClientComMessage, other string, ts time.Time) *ServerComMessage {
	return InfoUseOther(msg.ID, msg.Original, other, ts, msg.Timestamp)
}

// InfoAlreadySubscribed response means request to subscribe was ignored because user is already subscribed (304).
func InfoAlreadySubscribed(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusNotModified, // 304
			Text:      "already subscribed",
			Topic:     topic,
			Timestamp: ts,
		},
		Id: id, Timestamp: ts,
	}
}

// InfoNotJoined response means request to leave was ignored because user was not subscribed (304).
func InfoNotJoined(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusNotModified, // 304
			Text:      "not joined",
			Topic:     topic,
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// InfoNoAction response means request was ignored because the object was already in the desired state
// with explicit server and incoming request timestamps (304).
func InfoNoAction(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusNotModified, // 304
			Text:      "no action",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// InfoNoActionReply response means request was ignored because the object was already in the desired state
// in response to a client request (304).
func InfoNoActionReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return InfoNoAction(msg.ID, msg.Original, ts, msg.Timestamp)
}

// InfoNotModified response means update request was a noop (304).
func InfoNotModified(id, topic string, ts time.Time) *ServerComMessage {
	return InfoNotModifiedExplicitTs(id, topic, ts, ts)
}

// InfoNotModifiedReply response means update request was a noop
// in response to a client request (304).
func InfoNotModifiedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return InfoNotModifiedExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// InfoNotModifiedExplicitTs response means update request was a noop
// with explicit server and incoming request timestamps (304).
func InfoNotModifiedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusNotModified, // 304
			Text:      "not modified",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// InfoFound redirects to a new resource (307).
func InfoFound(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusTemporaryRedirect, // 307
			Text:      "found",
			Topic:     topic,
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// 4xx Errors

// ErrMalformed request malformed (400).
func ErrMalformed(id, topic string, ts time.Time) *ServerComMessage {
	return ErrMalformedExplicitTs(id, topic, ts, ts)
}

// ErrMalformedReply request malformed
// in response to a client request (400).
func ErrMalformedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrMalformedExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrMalformedExplicitTs request malformed with explicit server and incoming request timestamps (400).
func ErrMalformedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusBadRequest, // 400
			Text:      "malformed",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrAuthRequired authentication required  - user must authenticate first (401).
func ErrAuthRequired(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusUnauthorized, // 401
			Text:      "authentication required",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrAuthRequiredReply authentication required  - user must authenticate first
// in response to a client request (401).
func ErrAuthRequiredReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrAuthRequired(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrAuthFailed authentication failed
// with explicit server and incoming request timestamps (401).
func ErrAuthFailed(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusUnauthorized, // 401
			Text:      "authentication failed",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrAuthUnknownScheme authentication scheme is unrecognized or invalid (401).
func ErrAuthUnknownScheme(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusUnauthorized, // 401
			Text:      "unknown authentication scheme",
			Topic:     topic,
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// ErrPermissionDenied user is authenticated but operation is not permitted (403).
func ErrPermissionDenied(id, topic string, ts time.Time) *ServerComMessage {
	return ErrPermissionDeniedExplicitTs(id, topic, ts, ts)
}

// ErrPermissionDeniedExplicitTs user is authenticated but operation is not permitted
// with explicit server and incoming request timestamps (403).
func ErrPermissionDeniedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusForbidden, // 403
			Text:      "permission denied",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrPermissionDeniedReply user is authenticated but operation is not permitted
// with explicit server and incoming request timestamps in response to a client request (403).
func ErrPermissionDeniedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrPermissionDeniedExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrAPIKeyRequired  valid API key is required (403).
func ErrAPIKeyRequired(ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Code:      http.StatusForbidden,
			Text:      "valid API key required",
			Timestamp: ts,
		},
	}
}

// ErrSessionNotFound  valid API key is required (403).
func ErrSessionNotFound(ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Code:      http.StatusForbidden,
			Text:      "invalid or expired session",
			Timestamp: ts,
		},
	}
}

// ErrTopicNotFound topic is not found
// with explicit server and incoming request timestamps (404).
func ErrTopicNotFound(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusNotFound,
			Text:      "topic not found", // 404
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrTopicNotFoundReply topic is not found
// with explicit server and incoming request timestamps
// in response to a client request (404).
func ErrTopicNotFoundReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrTopicNotFound(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrUserNotFound user is not found
// with explicit server and incoming request timestamps (404).
func ErrUserNotFound(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusNotFound, // 404
			Text:      "user not found",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrUserNotFoundReply user is not found
// with explicit server and incoming request timestamps in response to a client request (404).
func ErrUserNotFoundReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrUserNotFound(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrNotFound is an error for missing objects other than user or topic (404).
func ErrNotFound(id, topic string, ts time.Time) *ServerComMessage {
	return ErrNotFoundExplicitTs(id, topic, ts, ts)
}

// ErrNotFoundExplicitTs is an error for missing objects other than user or topic
// with explicit server and incoming request timestamps (404).
func ErrNotFoundExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusNotFound, // 404
			Text:      "not found",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrNotFoundReply is an error for missing objects other than user or topic
// with explicit server and incoming request timestamps in response to a client request (404).
func ErrNotFoundReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrNotFoundExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrOperationNotAllowed a valid operation is not permitted in this context (405).
func ErrOperationNotAllowed(id, topic string, ts time.Time) *ServerComMessage {
	return ErrOperationNotAllowedExplicitTs(id, topic, ts, ts)
}

// ErrOperationNotAllowedExplicitTs a valid operation is not permitted in this context
// with explicit server and incoming request timestamps (405).
func ErrOperationNotAllowedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusMethodNotAllowed, // 405
			Text:      "operation or method not allowed",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrOperationNotAllowedReply a valid operation is not permitted in this context
// with explicit server and incoming request timestamps (405).
func ErrOperationNotAllowedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrOperationNotAllowedExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrInvalidResponse indicates that the client's response in invalid
// with explicit server and incoming request timestamps (406).
func ErrInvalidResponse(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusNotAcceptable, // 406
			Text:      "invalid response",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrAlreadyAuthenticated invalid attempt to authenticate an already authenticated session
// Switching users is not supported (409).
func ErrAlreadyAuthenticated(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusConflict, // 409
			Text:      "already authenticated",
			Topic:     topic,
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// ErrDuplicateCredential attempt to create a duplicate credential
// with explicit server and incoming request timestamps (409).
func ErrDuplicateCredential(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusConflict, // 409
			Text:      "duplicate credential",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrAttachFirst must attach to topic first in response to a client message (409).
func ErrAttachFirst(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        msg.ID,
			Code:      http.StatusConflict, // 409
			Text:      "must attach first",
			Topic:     msg.Original,
			Timestamp: ts,
		},
		Id:        msg.ID,
		Timestamp: msg.Timestamp,
	}
}

// ErrAlreadyExists the object already exists (409).
func ErrAlreadyExists(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusConflict, // 409
			Text:      "already exists",
			Topic:     topic,
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// ErrCommandOutOfSequence invalid sequence of comments, i.e. attempt to {sub} before {hi} (409).
func ErrCommandOutOfSequence(id, unused string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusConflict, // 409
			Text:      "command out of sequence",
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// ErrGone topic deleted or user banned (410).
func ErrGone(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusGone, // 410
			Text:      "gone",
			Topic:     topic,
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// ErrTooLarge packet or request size exceeded the limit (413).
func ErrTooLarge(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusRequestEntityTooLarge, // 413
			Text:      "too large",
			Topic:     topic,
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}

// ErrPolicy request violates a policy (e.g. password is too weak or too many subscribers) (422).
func ErrPolicy(id, topic string, ts time.Time) *ServerComMessage {
	return ErrPolicyExplicitTs(id, topic, ts, ts)
}

// ErrPolicyExplicitTs request violates a policy (e.g. password is too weak or too many subscribers)
// with explicit server and incoming request timestamps (422).
func ErrPolicyExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusUnprocessableEntity, // 422
			Text:      "policy violation",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrPolicyReply request violates a policy (e.g. password is too weak or too many subscribers)
// with explicit server and incoming request timestamps in response to a client request (422).
func ErrPolicyReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrPolicyExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrCallBusyExplicitTs indicates a "busy" reply to a video call request (486).
func ErrCallBusyExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      486, // Busy here.
			Text:      "busy here",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrCallBusyReply indicates a "busy" reply in response to a video call request (486)
func ErrCallBusyReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrCallBusyExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrUnknown database or other server error (500).
func ErrUnknown(id, topic string, ts time.Time) *ServerComMessage {
	return ErrUnknownExplicitTs(id, topic, ts, ts)
}

// ErrUnknownExplicitTs database or other server error with explicit server and incoming request timestamps (500).
func ErrUnknownExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusInternalServerError, // 500
			Text:      "internal error",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrUnknownReply database or other server error in response to a client request (500).
func ErrUnknownReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrUnknownExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrNotImplemented feature not implemented with explicit server and incoming request timestamps (501).
// TODO: consider changing status code to 4XX.
func ErrNotImplemented(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusNotImplemented, // 501
			Text:      "not implemented",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrNotImplementedReply feature not implemented error in response to a client request (501).
func ErrNotImplementedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrNotImplemented(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrClusterUnreachableReply in-cluster communication has failed error as response to a client request (502).
func ErrClusterUnreachableReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrClusterUnreachableExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrClusterUnreachable in-cluster communication has failed error with explicit server and
// incoming request timestamps (502).
func ErrClusterUnreachableExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusBadGateway, // 502
			Text:      "cluster unreachable",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrServiceUnavailableReply server overloaded error in response to a client request (503).
func ErrServiceUnavailableReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrServiceUnavailableExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrServiceUnavailableExplicitTs server overloaded error with explicit server and
// incoming request timestamps (503).
func ErrServiceUnavailableExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusServiceUnavailable, // 503
			Text:      "service unavailable",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrLocked operation rejected because the topic is being deleted (503).
func ErrLocked(id, topic string, ts time.Time) *ServerComMessage {
	return ErrLockedExplicitTs(id, topic, ts, ts)
}

// ErrLockedReply operation rejected because the topic is being deleted in response
// to a client request (503).
func ErrLockedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrLockedExplicitTs(msg.ID, msg.Original, ts, msg.Timestamp)
}

// ErrLockedExplicitTs operation rejected because the topic is being deleted
// with explicit server and incoming request timestamps (503).
func ErrLockedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusServiceUnavailable, // 503
			Text:      "locked",
			Topic:     topic,
			Timestamp: serverTs,
		},
		Id:        id,
		Timestamp: incomingReqTs,
	}
}

// ErrVersionNotSupported invalid (too low) protocol version (505).
func ErrVersionNotSupported(id string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{
		Ctrl: &MsgServerCtrl{
			Id:        id,
			Code:      http.StatusHTTPVersionNotSupported, // 505
			Text:      "version not supported",
			Timestamp: ts,
		},
		Id:        id,
		Timestamp: ts,
	}
}
