package message

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Client interface {
	Send(Message) error
}

type Handler interface {
	Procces(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	client Client
}

func NewHandler(c Client) (*handler, error) {
	if c == nil {
		return nil, errors.New("provided Client can not be nil")
	}
	return &handler{client: c}, nil
}

func (h *handler) Procces(w http.ResponseWriter, r *http.Request) {
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

	err = msg.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.sendMsg(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) sendMsg(msg Message) error {
	if msg.ExceedsLimit() {
		return h.sendConcatMsg(msg)
	}
	return h.client.Send(msg)
}

func (h *handler) sendConcatMsg(msg Message) error {
	cm := msg.concatenate()
	for _, msg := range cm {
		err := h.client.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
