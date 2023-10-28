package xmd

type BetMode string

func (m BetMode) String() string {
	switch m {
	case BetModeCustom:
		return "[Custom]  每期均进行投注"
	case BetModeModeAll:
		return "[Mode All]  当所选数字超过一定权重时，使用自定义投注模式投注，其他数字额外进行投注"
	case BetModeModeOnly:
		return "[Mode Only]  当所选数字超过一定权重时，仅使用自定义投注模式投注"
	case BetModeHalf:
		return "[Half]  按照权重大小顺序，选择进行一半的数字进行投注"
	default:
		return "<Undefined>"
	}
}

func (m BetMode) IsMode() bool {
	return m == BetModeModeAll || m == BetModeModeOnly
}

var (
	BetModeCustom   BetMode = "Custom"
	BetModeModeAll  BetMode = "Mode All"
	BetModeModeOnly BetMode = "Mode Only"
	BetModeHalf     BetMode = "Half"
)

type RunTime struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Config struct {
	BetMode   BetMode   `json:"bet_mode"`
	Sigma     float64   `json:"sigma"`
	RunTimes  []RunTime `json:"run_times"`
	Cookie    string    `json:"cookie"`
	UserAgent string    `json:"user_agent"`
	UserId    string    `json:"user_id"`
	Token     string    `json:"token"`
	Unix      string    `json:"unix"`
	KeyCode   string    `json:"key_code"`
	DeviceId  string    `json:"device_id"`
	ChannelId string    `json:"channel_id"`
}
