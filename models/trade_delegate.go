package models

import (
	"cointrade/utils"
	"fmt"
	"math/rand"
	"time"
)

func (m *TradeModel) MakeSn(uid int, t int) string {
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	switch t {
	case DELEGATE_TYPE_BUY:
		return fmt.Sprintf("%s%s%s%d", TRADE_BUY_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	case DELEGATE_TYPE_SELL:
		return fmt.Sprintf("%s%s%s%d", TRADE_SELL_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	}

	return fmt.Sprintf("%s%s%s%d", TRADE_SELL_PREFIX, timestr, uidstr, 10+rand.Intn(89))
}

func (m *TradeModel) DelegateTrade(uid int, rq *TradeDelegateRequest) *BaseResponse {
	return tradeSvc.DelegateTrade(uid, rq)
}
