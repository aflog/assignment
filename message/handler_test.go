package message

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewClientMock(err bool) *clientMock {
	return &clientMock{err: err}
}

type clientMock struct {
	err bool
}

func (cm *clientMock) Send(Message) error {
	if cm.err {
		return errors.New("client mock send error")
	}
	return nil
}

func TestNew(t *testing.T) {
	h, err := NewHandler(nil)
	assert.Error(t, err)
	assert.Nil(t, h)
	h, err = NewHandler(new(clientMock))
	assert.NoError(t, err)
	assert.NotNil(t, h)
}

func TestSendMsg(t *testing.T) {
	for _, test := range []struct {
		name       string
		json       string
		statusCode int
		hErr       bool
	}{
		{
			name:       "empty body",
			json:       "",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			json:       `{recipient:"+31612345678","originator":"MessageBird","message":"This is a test message."}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "invalid recipient",
			json:       `{"recipient":"","originator":"MessageBird","message":"This is a test message."}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "ok",
			json:       `{"recipient":"+31612345678","originator":"MessageBird","message":"This is a test message."}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "concatenated ok",
			json:       `{"recipient":"+31612345678","originator":"MessageBird","message":"text of more then 160 caracters that needs to be splitted in multiple messages and should be prepended by User Data Header, more random text END OF FIRST here starts second message asdiy doi asdoi asd hoiasd oiasd hoiasdh husaudg asdoiha oiasd oiasd oihsad ihasd oiasdi asdoi doiasdi asdi asd END OF SECOND third final part."}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "client error",
			json:       `{"recipient":"+31612345678","originator":"MessageBird","message":"This is a test message."}`,
			statusCode: http.StatusInternalServerError,
			hErr:       true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			h, _ := NewHandler(NewClientMock(test.hErr))
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "http://test.com", bytes.NewBufferString(test.json))
			h.SendMsg(w, r)
			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			assert.Equal(t, test.statusCode, resp.StatusCode)
			if resp.StatusCode == http.StatusOK {
				assert.Empty(t, body)
			} else {
				assert.NotEmpty(t, body)
			}
		})
	}
}
