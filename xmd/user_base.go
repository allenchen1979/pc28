package xmd

import (
	"fmt"
	"strings"
)

type UserBase struct {
	BetMode  BetMode
	RunTimes []RunTime

	cookie  string
	agent   string
	unix    string
	code    string
	device  string
	channel string
	id      string
	token   string
}

func (o UserBase) JoinString() string {
	ss := make([]string, 0, len(o.RunTimes))
	for _, cs := range o.RunTimes {
		ss = append(ss, fmt.Sprintf("%s ~ %s", cs.Start, cs.End))
	}

	return strings.Join(ss, ",")
}

func NewUserBase(betMode BetMode, runTimes []RunTime, cookie string, agent string, unix string, code string, device string, channel string, id string, token string) UserBase {
	return UserBase{
		BetMode:  betMode,
		RunTimes: runTimes,

		cookie:  cookie,
		agent:   agent,
		unix:    unix,
		code:    code,
		device:  device,
		channel: channel,
		id:      id,
		token:   token,
	}
}
