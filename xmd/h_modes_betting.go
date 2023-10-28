package xmd

import "fmt"

type XModesBettingRequest struct {
	Issue  string `json:"issue"`
	ModeId int    `json:"modeid"`

	Unix      string `json:"unix"`
	Keycode   string `json:"keycode"`
	PType     string `json:"ptype"`
	DeviceId  string `json:"deviceid"`
	ChannelId string `json:"channelid"`
	UserId    string `json:"userid"`
	Token     string `json:"token"`
}

type XModesBettingResponse struct {
	Status int      `json:"status"`
	Data   struct{} `json:"data"`
	Msg    string   `json:"msg"`
}

func hModesBetting(issue string, modeId int, userBase UserBase) error {
	betRequest := XModesBettingRequest{
		Issue:  issue,
		ModeId: modeId,

		Unix:      userBase.unix,
		Keycode:   userBase.code,
		PType:     "3",
		DeviceId:  userBase.device,
		ChannelId: userBase.channel,
		UserId:    userBase.id,
		Token:     userBase.token,
	}

	var betResponse XModesBettingResponse
	err := hDo(userBase, "POST", URLBetModesBetting, betRequest, &betResponse)
	if err != nil {
		return err
	}

	if betResponse.Status != 0 {
		return fmt.Errorf("%d %s", betResponse.Status, betResponse.Msg)
	}

	return nil
}
