package message

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUDH(t *testing.T) {
	udh := createUDH("a7", 3, 2)
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
			msg:  Message{Recipient: "+31123456789", Message: "test message", Originator: "test originator"},
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
		fmt.Sprintf("%x", "text of more then 160 caracters that needs to be splitted in multiple messages and should be prepended by User Data Header, more random text END OF FIRST"),
		fmt.Sprintf("%x", " here starts second message asdiy doi asdoi asd hoiasd oiasd hoiasdh husaudg asdoiha oiasd oiasd oihsad ihasd oiasdi asdoi doiasdi asdi asd END OF SECOND"),
		fmt.Sprintf("%x", " third final part"),
	}
	res := msg.Concatenate()
	assert.Len(t, res, 3)
	csms := ""
	for i, m := range res {
		assert.Equal(t, msg.Recipient, m.Recipient)
		assert.Equal(t, msg.Originator, m.Originator)
		if i == 0 {
			csms = m.UDH[6:8]
		}
		assert.Equal(t, expected[i], m.Message)
		assert.Equal(t, "050003", m.UDH[:6])
		assert.Equal(t, fmt.Sprintf("03%02d", i+1), m.UDH[8:12])
		assert.Equal(t, csms, m.UDH[6:8])
	}
}

func TestValidNumber(t *testing.T) {
	for _, test := range []struct {
		name string
		num  string
		err  bool
	}{
		{
			name: "empty",
			num:  "",
			err:  true,
		},
		{
			name: "not international format",
			num:  "0678123456",
			err:  true,
		},
		{
			name: "with white spaces",
			num:  "+31 123456789",
			err:  true,
		},
		{
			name: "non digits",
			num:  "+asghjkjkjhsdkjasd",
			err:  true,
		},
		{
			name: "short",
			num:  "+311234",
			err:  true,
		},
		{
			name: "long",
			num:  "+3112345678901234",
			err:  true,
		},
		{
			name: "correct format",
			num:  "+31123456789",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := validInterNum(test.num)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
