package server

import (
	"net/http"
	"time"
)

// debugSession is session debug info.
type debugSession struct {
	RemoteAddr string   `json:"remote_addr,omitempty"`
	Ua         string   `json:"ua,omitempty"`
	Uid        string   `json:"uid,omitempty"`
	Sid        string   `json:"sid,omitempty"`
	Clnode     string   `json:"clnode,omitempty"`
	Subs       []string `json:"subs,omitempty"`
}

// debugTopic is a topic debug info.
type debugTopic struct {
	Topic    string   `json:"topic,omitempty"`
	Xorig    string   `json:"xorig,omitempty"`
	IsProxy  bool     `json:"is_proxy,omitempty"`
	PerUser  []string `json:"per_user,omitempty"`
	PerSubs  []string `json:"per_subs,omitempty"`
	Sessions []string `json:"sessions,omitempty"`
}

// debugCachedUser is a user cache entry debug info.
type debugCachedUser struct {
	Uid    string `json:"uid,omitempty"`
	Unread int    `json:"unread,omitempty"`
	Topics int    `json:"topics,omitempty"`
}

// debugDump is server internal state dump for debugging.
type debugDump struct {
	Version   string            `json:"server_version,omitempty"`
	Build     string            `json:"build_id,omitempty"`
	Timestamp time.Time         `json:"ts,omitempty"`
	Sessions  []debugSession    `json:"sessions,omitempty"`
	Topics    []debugTopic      `json:"topics,omitempty"`
	UserCache []debugCachedUser `json:"user_cache,omitempty"`
}

// Wrapper around http.ResponseWriter which detects status set to 400+ and replaces
// default error message with a custom one.
type errorResponseWriter struct {
	status int
	http.ResponseWriter
}
