package xmd

import (
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

func fnBet(cache *Cache) error {
	issue := strconv.Itoa(cache.issue + 1)
	if cache.userBase.BetMode.IsMode() {
		ms := rand.Intn(5000)

		log.Printf("第%s期：使用自定义投注模式，等待%9.2f秒后执行投注 ...\n", strconv.Itoa(cache.issue), float64(ms)/1000)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	surplus, err := hGetGold(cache.userBase)
	if err != nil {
		return err
	}

	log.Printf("第%s期：开奖结果为【%d】，当前账户金额【%d】 >>>>>> \n", strconv.Itoa(cache.issue), cache.result, surplus)

	// Riddle
	bets, aBets, err := fnRiddle(cache, issue)
	if err != nil {
		return err
	}

	mrx := 1.0
	if cache.money < 1<<27 {
		mrx = float64(cache.money) / float64(1<<27) // 134,217,728
	}

	var m1Gold int
	if !cache.userBase.BetMode.IsMode() {
		if surplus < 1<<22 {
			// 4194304
			m1Gold = surplus / 75
		} else if surplus < 1<<23 {
			// 8388608
			m1Gold = surplus / 100
		} else if surplus < 1<<24 {
			// 16777216
			m1Gold = surplus / 125
		} else if surplus < 1<<25 {
			// 33554432
			m1Gold = surplus / 150
		} else if surplus < 1<<26 {
			// 67108864
			m1Gold = surplus / 175
		} else if surplus < 1<<27 {
			// 134217728
			m1Gold = surplus / 200
		} else if surplus < 1<<28 {
			// 268435456
			m1Gold = surplus / 225
		} else {
			m1Gold = surplus / 250
		}
	} else {
		m1Gold, err = hCustomModes(cache.userBase)
		if err != nil {
			return err
		}
	}

	if m1Gold*2 <= 100 {
		log.Printf("第%s期：最小投注额度【%d】小于设定额度【%d】不进行投注 ****** \n", issue, m1Gold*2, 100)
		return nil
	}

	switch cache.userBase.BetMode {
	case BetModeCustom:
		coverage, err := betSingle(cache, issue, mrx, m1Gold, bets)

		log.Printf("第%s期：投注覆盖率【%.3f%%】 !!!!!! \n", issue, coverage/10)
		if err != nil {
			return err
		}
	case BetModeModeAll:
		if err := betMode(cache, issue, m1Gold, bets, false); err != nil {
			return err
		}
	case BetModeModeOnly:
		if err := betMode(cache, issue, m1Gold, bets, true); err != nil {
			return err
		}
	case BetModeHalf:
		if err := betHalf(cache, issue, mrx, m1Gold, aBets); err != nil {
			return err
		}
	}

	return nil
}

func betMode(cache *Cache, issue string, m1Gold int, bets map[int]float64, isOnly bool) error {
	if !isOnly {
		rs := make([]int, 0, len(bets))
		for result := range bets {
			rs = append(rs, result)
		}

		sort.Ints(rs)
		log.Printf("第%s期：竞猜结果 %q \n", issue, fmtIntSlice(rs))
	}

	md := 400
	modeId, modeName := modeFn(bets, md)
	if modeId == 0 {
		log.Printf("第%s期：所有自定义投注模式权重均不超过【%d】，不进行投注 ****** \n", issue, md)
		return nil
	}

	log.Printf("第%s期：使用的自定义投注模式为 %q \n", issue, modeName)
	if err := hModesBetting(issue, modeId, cache.userBase); err != nil {
		return err
	}

	if isOnly {
		return nil
	}

	extras := extraFn(modeId, m1Gold, bets)
	if len(extras) > 0 {
		log.Printf("第%s期：额外投注数字 %q \n", issue, fmtIntSlice(m2sFn(extras)))
	}

	stdBets := []int{200000, 50000, 10000, 5000, 2000, 1000, 500}
	betMaps := make(map[int][]int)

	for _, stdBet := range stdBets {
		betSlice, ok := betMaps[stdBet]
		if !ok {
			betSlice = make([]int, 0)
		}

		for result, betGold := range extras {
			qn := betGold / stdBet
			if qn > 0 {
				for i := 0; i < qn; i++ {
					betSlice = append(betSlice, result)
				}

				extras[result] = betGold - qn*stdBet
			}
		}

		sort.Ints(betSlice)
		betMaps[stdBet] = betSlice
	}

	for _, stdBet := range stdBets {
		if len(betMaps[stdBet]) > 0 {
			log.Printf("第%s期：使用投注额度【%-6d】，投注数字【%s】\n", issue, stdBet, fmtIntSlice(betMaps[stdBet]))
		}

		for _, result := range betMaps[stdBet] {
			if err := hBetting1(issue, stdBet, result, cache.userBase); err != nil {
				return err
			}
		}
	}

	return nil
}

func betSingle(cache *Cache, issue string, mrx float64, m1Gold int, bets map[int]float64) (float64, error) {
	var coverage float64

	for _, result := range SN28 {
		if _, ok := bets[result]; !ok || mrx*bets[result] <= 0.01 {
			continue
		}

		fGold := mrx * bets[result] * float64(2*m1Gold) * float64(STDS1000[result]) / 1000

		var iGold int
		if fGold >= 1<<16 {
			iGold = int(math.Round(fGold/2000.0) * 2000)
		} else if fGold >= 1<<15 {
			iGold = int(math.Round(fGold/1500.0) * 1500)
		} else if fGold >= 1<<14 {
			iGold = int(math.Round(fGold/1000.0) * 1000)
		} else {
			iGold = int(math.Round(fGold/500.0) * 500)
		}

		if err := hBetting1(issue, iGold, result, cache.userBase); err != nil {
			return coverage, err
		}

		coverage = coverage + float64(STDS1000[result])*bets[result]
	}

	return coverage, nil
}

func betHalf(cache *Cache, issue string, mrx float64, m1Gold int, aBets map[int]float64) error {
	type BetResult struct {
		Result int
		Rx     float64
	}

	rs := make([]BetResult, 0, len(aBets))
	for result, rx := range aBets {
		rs = append(rs, BetResult{Result: result, Rx: rx})
	}
	sort.Slice(rs, func(i, j int) bool { return rs[i].Rx > rs[j].Rx })

	var coverage float64

	ns := make([]int, 0)
	for _, s := range rs {
		if coverage > 500 {
			break
		}

		betGold := int(mrx * float64(2*m1Gold) * float64(STDS1000[s.Result]) / 1000)
		if err := hBetting1(issue, betGold, s.Result, cache.userBase); err != nil {
			return err
		}

		ns = append(ns, s.Result)
		coverage = coverage + float64(STDS1000[s.Result])
	}
	log.Printf("第%s期：投注数字为 %s ，覆盖率为 %.2f%% \n", issue, fmtIntSlice(ns), coverage/10)

	return nil
}
