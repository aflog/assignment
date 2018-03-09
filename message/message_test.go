package message

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUDH(t *testing.T) {
	udh := getUDH("a7", 3, 2)
	assert.Equal(t, "050003a70302", udh)
}

func TestSplitString(t *testing.T) {
	for _, test := range []struct {
		name   string
		s      string
		length int
		res    []string
	}{
		{
			name:   "weird caracters",
			s:      "a#f界世",
			length: 2,
			res:    []string{"a#", "f界", "世"},
		},
		{
			name:   "empty string",
			s:      "",
			length: 2,
			res:    nil,
		},
		{
			name:   "zero length",
			s:      "abcd",
			length: 0,
			res:    nil,
		},
		{
			name:   "empty message",
			s:      "text of more then 160 caracters that needs to be splitted in multiple messages and should be prepended by User Data Header, more random text END OF FIRST here starts second message asdiy doi asdoi asd hoiasd oiasd hoiasdh husaudg asdoiha oiasd oiasd oihsad ihasd oiasdi asdoi doiasdi asdi asd END OF SECOND third final part",
			length: 153,
			res: []string{
				"text of more then 160 caracters that needs to be splitted in multiple messages and should be prepended by User Data Header, more random text END OF FIRST",
				" here starts second message asdiy doi asdoi asd hoiasd oiasd hoiasdh husaudg asdoiha oiasd oiasd oihsad ihasd oiasdi asdoi doiasdi asdi asd END OF SECOND",
				" third final part",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			res := splitString(test.s, test.length)
			assert.Equal(t, test.res, res)
		})
	}
}

func TestValid(t *testing.T) {
	for _, test := range []struct {
		name string
		msg  Message
		err  bool
	}{
		{
			name: "empty message",
			err:  true,
		},
		{
			name: "empty message text",
			msg:  Message{Recipient: "123", Originator: "test originator"},
			err:  true,
		},
		{
			name: "empty originator",
			msg:  Message{Recipient: "123", Message: "test message"},
			err:  true,
		},
		{
			name: "empty recipient",
			msg:  Message{Message: "test message", Originator: "test originator"},
			err:  true,
		},
		{
			name: "ok",
			msg:  Message{Recipient: "123", Message: "test message", Originator: "test originator"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := test.msg.Validate()
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConcatenate(t *testing.T) {
	msg := Message{
		Originator: "test originator",
		Recipient:  "123456789",
		Message:    "text of more then 160 caracters that needs to be splitted in multiple messages and should be prepended by User Data Header, more random text END OF FIRST here starts second message asdiy doi asdoi asd hoiasd oiasd hoiasdh husaudg asdoiha oiasd oiasd oihsad ihasd oiasdi asdoi doiasdi asdi asd END OF SECOND third final part",
	}
	expected := []string{
		"text of more then 160 caracters that needs to be splitted in multiple messages and should be prepended by User Data Header, more random text END OF FIRST",
		" here starts second message asdiy doi asdoi asd hoiasd oiasd hoiasdh husaudg asdoiha oiasd oiasd oihsad ihasd oiasdi asdoi doiasdi asdi asd END OF SECOND",
		" third final part",
	}
	res := msg.concatenate()
	assert.Len(t, res, 3)
	csms := ""
	for i, m := range res {
		assert.Equal(t, msg.Recipient, m.Recipient)
		assert.Equal(t, msg.Originator, m.Originator)
		byteMsg := []byte(m.Message)
		if i == 0 {
			csms = string(byteMsg[6:8])
		}
		assert.Equal(t, expected[i], string(byteMsg[12:]))
		assert.Equal(t, "050003", string(byteMsg[:6]))
		assert.Equal(t, fmt.Sprintf("03%02d", i), string(byteMsg[8:12]))
		assert.Equal(t, csms, string(byteMsg[6:8]))
	}
}
