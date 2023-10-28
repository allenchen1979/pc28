package xmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Cache struct {
	secs     float64
	sigma    float64
	userBase UserBase

	issue  int
	result int
	money  int
}

func NewCache(dir string) (*Cache, error) {
	bs, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := json.Unmarshal(bs, &conf); err != nil {
		return nil, err
	}

	userBase := NewUserBase(
		conf.BetMode, conf.RunTimes, conf.Cookie, conf.UserAgent,
		conf.Unix, conf.KeyCode, conf.DeviceId, conf.ChannelId, conf.UserId, conf.Token,
	)

	cache := &Cache{
		secs:     conf.Secs,
		sigma:    conf.Sigma,
		userBase: userBase,

		issue:  -1,
		result: -1,
		money:  -1,
	}

	return cache, nil
}

func (o *Cache) Sync(size int) error {
	items, err := hAnalyseHistory(size, o.userBase)
	if err != nil {
		return err
	}

	for i, item := range items {
		issue, err := strconv.Atoi(item.Issue)
		if err != nil {
			return err
		}

		result, err := strconv.Atoi(item.Result)
		if err != nil {
			return err
		}

		money, err := strconv.Atoi(strings.ReplaceAll(item.Money, ",", ""))
		if err != nil {
			return err
		}

		if i == 0 {
			o.issue = issue
			o.result = result
			o.money = money
		}
	}

	return nil
}
