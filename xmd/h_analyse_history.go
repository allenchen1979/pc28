package xmd

import (
	"errors"
	"fmt"
)

type QHistoryItem struct {
	Issue  string `json:"issue"`
	Result string `json:"lresult"`
	Money  string `json:"tmoney"`
	Member int    `json:"tmember"`
}

type QHistoryData struct {
	Items []QHistoryItem `json:"items"`
}

type QHistory struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`

	Data QHistoryData `json:"data"`
}

type QHistoryRequest struct {
	PageSize  int    `json:"pagesize"`
	Unix      string `json:"unix"`
	KeyCode   string `json:"keycode"`
	PType     string `json:"ptype"`
	DeviceId  string `json:"deviceid"`
	ChannelId string `json:"channelid"`
	UserId    string `json:"userid"`
	Token     string `json:"token"`
}

func hAnalyseHistory(pageSize int, userBase UserBase) ([]QHistoryItem, error) {

	hisRequest := QHistoryRequest{
		PageSize:  pageSize,
		PType:     "3",
		Unix:      userBase.unix,
		KeyCode:   userBase.code,
		DeviceId:  userBase.device,
		ChannelId: userBase.channel,
		UserId:    userBase.id,
		Token:     userBase.token,
	}

	var hisResponse QHistory

	err := hDo(userBase, "POST", URLBetAnalyseHistory, hisRequest, &hisResponse)
	if err != nil {
		return nil, err
	}

	if hisResponse.Status != 0 {
		return nil, fmt.Errorf("%d %s", hisResponse.Status, hisResponse.Msg)
	}

	if len(hisResponse.Data.Items) < 1 {
		return nil, errors.New("empty rows")
	}

	return hisResponse.Data.Items, nil
}
