package message

import (
	"errors"
	"fmt"
	"math/rand"
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

type Message struct {
	Recipient  string `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
	UDH        string
}

func (m *Message) Validate() error {
	errMsg := []string{}
	if m.Recipient == "" {
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

func (m *Message) concatenate() (ms []Message) {
	csms := randCSMSrefNum()
	ss := splitString(m.Message, concatMsgLength)
	for i, s := range ss {
		udh := getUDH(csms, len(ss), i+1)
		//binaryMsg := intoBinary(s)
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

func (m *Message) ExceedsLimit() bool {
	return utf8.RuneCountInString(m.Message) > maxMsgLength
}

func getUDH(csms string, total, num int) string {
	udh := fmt.Sprintf("%02d%02d%02d%s%02X%02X", udhLength, infoElementID, headerLength, csms, total, num)
	return udh
}

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

func intoBinary(s string) string {
	return s
}
