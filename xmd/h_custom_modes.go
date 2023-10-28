package xmd

import (
	"fmt"
	"strconv"
	"strings"
)

type XCustomModesRequest struct {
	Unix      string `json:"unix"`
	Keycode   string `json:"keycode"`
	PType     string `json:"ptype"`
	DeviceId  string `json:"deviceid"`
	ChannelId string `json:"channelid"`
	UserId    string `json:"userid"`
	Token     string `json:"token"`
}

type XCustomModesResponse struct {
	Status int `json:"status"`
	Data   struct {
		Items []struct {
			Modeid     int    `json:"modeid"`
			Name       string `json:"name"`
			Eggs       string `json:"eggs"`
			Winmodeid  int    `json:"winmodeid"`
			Losemodeid int    `json:"losemodeid"`
		} `json:"items"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func hCustomModes(userBase UserBase) (int, error) {
	betRequest := XCustomModesRequest{
		Unix:      userBase.unix,
		Keycode:   userBase.code,
		PType:     "3",
		DeviceId:  userBase.device,
		ChannelId: userBase.channel,
		UserId:    userBase.id,
		Token:     userBase.token,
	}

	var betResponse XCustomModesResponse
	err := hDo(userBase, "POST", URLBetCustomModes, betRequest, &betResponse)
	if err != nil {
		return 0, err
	}

	if betResponse.Status != 0 {
		return 0, fmt.Errorf("%d %s", betResponse.Status, betResponse.Msg)
	}

	m1 := betResponse.Data.Items[0]
	eggs := strings.Split(m1.Eggs, ",")

	var total int
	for _, egg := range eggs {
		i, err := strconv.Atoi(egg)
		if err != nil {
			return 0, err
		}

		total = total + i
	}

	return total, nil
}
