package redis

import (
	"cointrade/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	Host     string
	Port     int
	User     string
	Pass     string
	conn     *redis.Client
	Pollsize int
	Minconns int
}
type RedisError struct {
}

func (r *RedisError) Error() string {
	return "redis error"
}
func (r *Redis) SetLinkInfo(Host string, port int, user string, pass string, poolsize int, minconns int) {
	r.Host = Host
	r.Port = port
	r.User = user
	r.Pass = pass
	r.Pollsize = poolsize
	r.Minconns = minconns
}
func (r *Redis) init() *redis.Options {
	options := new(redis.Options)
	options.Addr = r.Host + ":" + strconv.Itoa(r.Port)
	options.Password = r.Pass
	options.MinIdleConns = r.Minconns
	options.DialTimeout = 3 * time.Second
	options.PoolTimeout = 3 * time.Second
	options.PoolSize = r.Pollsize

	return options
}
func (r *Redis) Connect() error {
	c := redis.NewClient(r.init())
	r.conn = c
	return nil
}
func (r *Redis) Set(key string, value interface{}, t time.Duration) *redis.StatusCmd {
	return r.conn.Set(key, value, t)
}
func (r *Redis) Get(key string) string {
	c := r.conn.Get(key)
	if c != nil {
		return c.Val()
	}
	return ""
}
func (r *Redis) Delete(key string) {
	r.conn.Del(key)
}
func (r *Redis) SetValue(hashname string, key string, value interface{}) *redis.BoolCmd {
	//redis存储内容
	jsonstr, _ := json.Marshal(value)
	//utils.Log("json:", string(jsonstr))
	//utils.Log("key:", key)

	b := r.conn.HSet(hashname, key, string(jsonstr))
	_, err := b.Result()
	if err != nil {
		utils.Log("redis error:", err.Error())
	}
	return b
}

func (r *Redis) KeppConnect() {
	if str, _ := r.conn.Ping().Result(); str != "PONG" {
		fmt.Println("重连...")
		r.Connect()
	} else {
		fmt.Println("  链接状态！", str)
	}
}

func (r *Redis) Expire(key string, t time.Duration) {
	r.KeppConnect()
	r.conn.Expire(key, time.Second*t)
}
func (r *Redis) GetValue(hashname string, key string) (interface{}, error) { //获取REDIS存储内容
	r.KeppConnect()
	b := r.conn.HGet(hashname, key)
	str := b.Val()
	var obj interface{}
	fmt.Println(" str 的取值为", key, "===============>", str)
	err := json.Unmarshal([]byte(str), &obj)
	if err != nil {
		return str, err
	}
	return obj, err
}
func (r *Redis) GetObject(hashname string, key string, v interface{}) error {
	b := r.conn.HGet(hashname, key)
	str := b.Val()
	err := json.Unmarshal([]byte(str), v)
	return err
}
func (r *Redis) Del(hashname string, key string) {
	r.conn.HDel(hashname, key)
}

func (r *Redis) PushQueue(hasname string, value interface{}) { //推入队列
	jsonstr, _ := json.Marshal(value)
	cmd := r.conn.RPush(hasname, string(jsonstr))
	_, err := cmd.Result()
	if err != nil {
		utils.Log("redis error:", err.Error())
	}
}

func (r *Redis) PopQueue(hashname string) []string { //取出队列
	/*rs := r.conn.LRange(hashname, 0, int64(n-1))

	ln := len(rs.Val())
	//println(rs.Val())
	if ln > 0 {
		utils.Log(rs.Val())
	}
	r.conn.LTrim(hashname, int64(ln), -1)
	return rs.Val()*/
	rs := r.conn.BLPop(time.Duration(3*time.Second), hashname).Val()
	if len(rs) > 1 {
		return rs[1:]
	}
	return rs
}
func (r *Redis) Pop(hashname string) string {
	rs := r.conn.LPop(hashname)
	return rs.Val()
}
func (r *Redis) GetAll(hashname string) map[string]string {
	rs := r.conn.HGetAll(hashname)
	return rs.Val()
}
