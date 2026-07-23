package lib

import (
	"cointrade/utils"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	BET_LOCK_TIME   = 10 //锁单时间 单位秒
	PER_MINUTE      = 2  //每一期的时间
	GAMEMODE_SINGLE = 0  //单选
	GAMEMODE_AB     = 1  //前2
	GAMEMODE_ABC    = 2  //前3
	GAMEMODE_ABCD   = 3  //前4
	GAMEMODE_ABCDE  = 4  //前5
)

type LuckSN struct {
	Presn    string `json:"presn"`    //上一期的期号
	Sn       string `json:"sn"`       //当前期号
	LockTime int    `json:"locktime"` //锁单时间
	LeftTime int    `json:"lefttime"` //剩余时间 开奖剩余时间
	DiffTime int    `json:"difftime"` //锁单时间间隔
}

func CreateSn() *LuckSN { //获取当前期号信息
	rs := new(LuckSN)
	now := time.Now()
	//s := now.Format("200601021504")
	sn_prefix := now.Format("20060102")
	diff_minute := now.Hour()*60 + now.Minute() + 1
	n := math.Ceil(float64(diff_minute) / float64(PER_MINUTE))
	pre_n := n - 1
	if pre_n == 0 {
		//上一期应该是昨天的最后一期
		pre_now := now.AddDate(0, 0, -1)
		pre_prefix := pre_now.Format("20060102")
		rs.Presn = pre_prefix + PaddingSN(strconv.Itoa(24*60/PER_MINUTE))
	} else {
		rs.Presn = sn_prefix + PaddingSN(strconv.Itoa(int(n)-1))
	}
	rs.Sn = sn_prefix + PaddingSN(strconv.Itoa(int(n)))
	rs.LockTime = GetLockTime(now, int(n))
	rs.LeftTime = rs.LockTime - int(now.Unix()) + BET_LOCK_TIME
	rs.DiffTime = BET_LOCK_TIME
	return rs
}
func PaddingSN(sn string) string {
	rs := sn
	n := len(sn)
	diff_n := 4 - n
	for i := 0; i < diff_n; i++ {
		rs = "0" + rs
	}
	return sn
}
func GetLockTime(now time.Time, n int) int {
	//根据期号获得锁单时间
	n1 := now.Unix()
	n1 = n1 - int64(now.Hour())*60*60 - int64(now.Minute())*60 - int64(now.Second())
	n1 = n1 + int64(n*PER_MINUTE*60)
	return int(n1) - BET_LOCK_TIME
}
func CheckBeting(betcontent string) int {
	//检测投注的合法性 反回为总注数 为0则投注不合法
	rs := 0
	s_arr := strings.Split(betcontent, ",")
	if len(s_arr) < 5 {
		return rs
	}
	for _, v := range s_arr {
		v = strings.Replace(v, " ", "", -1) //去除所有空格
		_, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			continue
		}
		rs = rs + len(v)
	}
	return rs
}

func GetBetCount(betcount string, mode int) int {
	return 0
}
func CreateResult() string {
	//返回投注结果
	s := ""
	for i := 0; i < 5; i++ {
		s = s + strconv.Itoa(rand.Intn(9))
	}
	return s
}
func CheckWin(betcontent string, result string) int { //返回中奖总注数
	rs := 0
	utils.ServiceInfo("lucky result:", result)
	bet_arr := strings.Split(betcontent, ",")
	if len(bet_arr) < 5 {
		return rs
	}
	for k, v := range bet_arr {
		if k >= 5 { //防止恶意数据导致崩溃
			break
		}
		if strings.Index(v, result[k:k+1]) >= 0 {
			rs = rs + 1
		}
	}
	return rs
}
