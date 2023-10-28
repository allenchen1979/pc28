package xmd

import (
	"log"
	"strconv"
	"time"
)

func Run(cache *Cache) {
	secs := cache.secs
	if cache.userBase.BetMode.IsMode() {
		secs = 42.50
	}

	log.Printf("可用的投注模式包括：\n  %s \n  %s \n  %s \n  %s \n", BetModeCustom, BetModeModeOnly, BetModeModeAll, BetModeHalf)
	log.Printf("当前使用的投注模式为 %q \n", cache.userBase.BetMode)
	log.Printf("当前使用的权重参数为 %.3f \n", cache.sigma)
	log.Printf("当前设定的投注时间段：%q \n", cache.userBase.JoinString())

	dua := time.Now().Sub(time.Now().Truncate(time.Minute))
	log.Printf("在 %.4f 秒后开始运行投注 \n", secs-dua.Seconds())
	time.Sleep(time.Second * time.Duration(secs-dua.Seconds()))

	go runTask(cache)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	log.Println("自动投注正在运行中 ...")
	for {
		select {
		case <-ticker.C:
			go runTask(cache)
		}
	}
}

func runTask(cache *Cache) {
	if err := cache.Sync(200); err != nil {
		log.Println(err.Error())
	}

	isBet, hms := false, time.Now().Format("15:04")

	for _, cs := range cache.userBase.RunTimes {
		if hms >= cs.Start && hms <= cs.End {
			isBet = true
			break
		}
	}

	if !isBet {
		log.Printf("第%s期：不在设定的投注时间段范围内 ****** \n", strconv.Itoa(cache.issue+1))
		return
	}

	if err := fnBet(cache); err != nil {
		log.Println(err.Error())
	}
}
