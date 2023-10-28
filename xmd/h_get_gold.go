package xmd

import (
	"strconv"
	"strings"
)

type UserBaseRequest struct {
	Unix      string `json:"unix"`
	KeyCode   string `json:"keycode"`
	PType     string `json:"ptype"`
	DeviceId  string `json:"deviceid"`
	ChannelId string `json:"channelid"`
	UserId    string `json:"userid"`
	Token     string `json:"token"`
}

type UserBaseResponse struct {
	Status int `json:"status"`
	Data   struct {
		GoldEggs string `json:"goldeggs"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func hGetGold(userBase UserBase) (gold int, err error) {
	userBaseRequest := UserBaseRequest{
		Unix:      userBase.unix,
		KeyCode:   userBase.code,
		PType:     "3",
		DeviceId:  userBase.device,
		ChannelId: userBase.channel,
		UserId:    userBase.id,
		Token:     userBase.token,
	}

	var userBaseResponse UserBaseResponse

	err = hDo(userBase, "POST", URLBetUserBase, userBaseRequest, &userBaseResponse)
	if err != nil {
		return
	}

	if userBaseResponse.Status != 0 {
		return gold, err
	}

	sGold := strings.ReplaceAll(userBaseResponse.Data.GoldEggs, ",", "")
	iGold, err := strconv.Atoi(sGold)
	if err != nil {
		return gold, err
	}

	return iGold, nil
}
