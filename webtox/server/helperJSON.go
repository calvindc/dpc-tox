package main

import (
	"encoding/json"
	"github.com/calvindc/dpc-tox"
	"net/http"
)

// rejectWithErrorJSON writes an error encoded as JSON to a http.ResponseWriter
// w        the http.ResponseWriter
// code     an error code that identifies the error
// message  a message explaining what went wrong (should be human readable)
func rejectWithErrorJSON(w http.ResponseWriter, code string, message string) {
	type Err struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	e := Err{Code: code, Message: message}
	jsonErr, _ := json.Marshal(e)
	http.Error(w, string(jsonErr), 422)
}

// rejectWithDefaultErrorJSON writes a default error encoded as JSON to a
// http.ResponseWriter. rejectWithDefaultErrorJSON(w) is equivalent to
// rejectWithErrorJSON(w, "unknown", "An unknown error occoured."))
// w  the http.ResponseWriter
func rejectWithDefaultErrorJSON(w http.ResponseWriter) {
	type Err struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	e := Err{Code: "unknown", Message: "An unknown error occoured."}
	jsonErr, _ := json.Marshal(e)
	http.Error(w, string(jsonErr), 422)
}

// rejectWithFriendErrorJSON writes a dpc_tox.ToxErrFriendAdd error encoded as
// JSON to a http.ResponseWriter
// w    the http.ResponseWriter
// err  the dpc_tox.ToxErrFriendAdd error to be encoded
func rejectWithFriendErrorJSON(w http.ResponseWriter, err error) {
	switch err {
	case dpc_tox.ErrFriendAddNoMessage:
		rejectWithErrorJSON(w, "no_message", "An invitation message is required.")
		return
	case dpc_tox.ErrFriendAddTooLong:
		rejectWithErrorJSON(w, "invalid_message", "The message you entered is too long.")
		return
	case dpc_tox.ErrFriendAddOwnKey:
		fallthrough
	case dpc_tox.ErrFriendAddBadChecksum:
		fallthrough
	case dpc_tox.ErrFriendAddSetNewNospam:
		rejectWithErrorJSON(w, "invalid_toxid", "The Tox ID you entered is invalid.")
		return
	case dpc_tox.ErrFriendAddAlreadySent:
		rejectWithErrorJSON(w, "already_send", "A friend request to this person has already send.")
		return
	default:
		rejectWithDefaultErrorJSON(w)
		return
	}
}

// createSimpleJSONEvent creates a simple JSON event used in a WS connection
// name  the name of the type of the event
func createSimpleJSONEvent(name string) string {
	type jsonEvent struct {
		Type string `json:"type"`
	}

	e, _ := json.Marshal(jsonEvent{
		Type: name,
	})

	return string(e)
}
