package message

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

//Client defines an interface for a connection to a third party sms service
type Client interface {
	Send(Message) error
}

//Handler defines an interface of a message handler
type Handler interface {
	SendMsg(w http.ResponseWriter, r *http.Request)
}

//handler holds set up for a message hadler and implements the Handler interface
type handler struct {
	client Client
}

//NewHandler creates and sets up a new message handler and returns it
//as a Handler interface, will return an error in case when an invalid Client
//is provided
func NewHandler(c Client) (Handler, error) {
	if c == nil {
		return nil, errors.New("provided Client can not be nil")
	}
	return &handler{client: c}, nil
}

//SendMsg reads message from a request body, validates it and sends it to the third
//party sms service. Returns 202 HTTP status in case of success, 400 in case of
//an ivalid message format and 500 for failure that occurs for other reasons.
func (h *handler) SendMsg(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Unmarshal
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validate message data
	err = msg.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Send message
	err = h.sendMsgToService(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//sendMsgToService sends message to the third party sms service
func (h *handler) sendMsgToService(msg Message) error {
	if msg.ExceedsLimit() {
		return h.sendConcatMsgToService(msg)
	}
	return h.client.Send(msg)
}

//sendConcatMsgToService sends the message received to the third party sms service
//as a concatenated message
func (h *handler) sendConcatMsgToService(msg Message) error {
	cm := msg.Concatenate()
	for _, msg := range cm {
		err := h.client.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
