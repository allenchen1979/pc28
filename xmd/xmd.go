package xmd

import (
	"log"
	"strconv"
	"time"
)

func Run(cache *Cache) {
	secs := 53.50
	if cache.userBase.BetMode.IsMode() {
		secs = 42.50
	}

	dua := time.Now().Sub(time.Now().Truncate(time.Minute))
	log.Printf("在 %.4f 秒后开始运行投注 \n", secs-dua.Seconds())
	time.Sleep(time.Second * time.Duration(secs-dua.Seconds()))

	go runTask(cache)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

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
