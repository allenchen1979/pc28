package xmd

type XBet struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type XBetRequest struct {
	Issue     string `json:"issue"`
	GoldEggs  int    `json:"totalgoldeggs"`
	Numbers   int    `json:"numbers"`
	Unix      string `json:"unix"`
	Keycode   string `json:"keycode"`
	PType     string `json:"ptype"`
	DeviceId  string `json:"deviceid"`
	ChannelId string `json:"channelid"`
	Userid    string `json:"userid"`
	Token     string `json:"token"`
}

func hBetting1(issue string, betGold int, result int, userBase UserBase) error {
	if betGold <= 10 {
		return nil
	}

	betRequest := XBetRequest{
		Issue:     issue,
		GoldEggs:  betGold,
		Numbers:   result,
		Unix:      userBase.unix,
		Keycode:   userBase.code,
		PType:     "3",
		DeviceId:  userBase.device,
		ChannelId: userBase.channel,
		Userid:    userBase.id,
		Token:     userBase.token,
	}

	var betResponse XBet
	err := hDo(userBase, "POST", URLBetBetting1, betRequest, &betResponse)
	if err != nil {
		return err
	}

	if betResponse.Status != 0 {
		return err
	}

	return nil
}
