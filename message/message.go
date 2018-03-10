package message

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	maxMsgLength    = 160
	concatMsgLength = 153
	udhLength       = 05
	infoElementID   = 00
	headerLength    = 03
)

//Message holds information of one message
type Message struct {
	Recipient  string `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
	UDH        string
}

//Validate checks wheather the message is composed by a valid phone number and
//that nor recipient neither message text are empty and returns error if any of
//the checks doesn't pass
func (m *Message) Validate() error {
	errMsg := []string{}
	if err := validInterNum(m.Recipient); err != nil {
		errMsg = append(errMsg, fmt.Sprintf("%s is not a valid recipient number", m.Recipient))
	}
	if m.Originator == "" {
		errMsg = append(errMsg, "originator can not be empty")
	}
	if m.Message == "" {
		errMsg = append(errMsg, "message can not be empty")
	}
	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg[:], ", "))
	}
	return nil
}

//ExceedsLimit checks whether the lenght of the message exceeds the maximum limit
func (m *Message) ExceedsLimit() bool {
	return utf8.RuneCountInString(m.Message) > maxMsgLength
}

//Concatenate divides the message into an array of messages according to the
//maximum length defined for concatenated messages,sets the UDH for all created
//messages and format it's text as binary (hexadecimal representation)
func (m *Message) Concatenate() (ms []Message) {
	csms := randCSMSrefNum()
	ss := splitString(m.Message, concatMsgLength)
	for i, s := range ss {
		udh := createUDH(csms, len(ss), i+1)
		nm := Message{
			Recipient:  m.Recipient,
			Originator: m.Originator,
			Message:    fmt.Sprintf("%x", s),
			UDH:        udh,
		}
		ms = append(ms, nm)
	}
	return
}

//createUDH returns concatenated message UDH for a specified csms, total amount
//of message and the position of the message in the sequence
func createUDH(csms string, total, num int) string {
	udh := fmt.Sprintf("%02d%02d%02d%s%02X%02X", udhLength, infoElementID, headerLength, csms, total, num)
	return udh
}

//splitString splits string by specified length
func splitString(s string, length int) (ss []string) {
	if length < 1 {
		return
	}
	r := []rune(s)
	for from := 0; from < len(r); from += length {
		till := from + length
		if till > len(r) {
			till = len(r)
		}
		ss = append(ss, string(r[from:till]))
	}
	return
}

//returns a random CSMS reference number
func randCSMSrefNum() string {
	ri := rand.Intn(256) //int 0-255 equals hex 00-FF
	return fmt.Sprintf("%02X", ri)
}

//validInterNum validates if string is an international number and returns an
//error in case the received string doesn't comply with the specified format,
//a valid format is '+'followed by 7-15 digits
func validInterNum(n string) error {
	if _, err := strconv.Atoi(n); err != nil {
		return errors.New("is not a number")
	}
	if n[0:1] != "+" {
		return errors.New("is not in internatinal format, number has to start with '+'")
	}
	if len(n) < 8 || len(n) > 16 {
		return errors.New("number does not have correct length")
	}
	return nil
}
