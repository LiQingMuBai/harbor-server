package utils

import (
	"bufio"
	"cointrade/lib/db"

	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wxnacy/wgo/arrays"
)

var (
	DEBUGMODE       = 1
	PAGE_LIMIT_SIZE = 20

	utilsEnvOnce sync.Once
)

func init() {
	loadRuntimeEnv()
}

func loadRuntimeEnv() {
	utilsEnvOnce.Do(func() {
		loadDotEnvFile()
		DEBUGMODE = getEnvInt("DEBUGMODE", DEBUGMODE)
		PAGE_LIMIT_SIZE = getEnvInt("PAGE_LIMIT_SIZE", PAGE_LIMIT_SIZE)
	})
}

func loadDotEnvFile() {
	candidates := []string{
		".env",
		"../.env",
		"../../.env",
	}
	if envPath := strings.TrimSpace(os.Getenv("APP_ENV_FILE")); envPath != "" {
		candidates = append([]string{envPath}, candidates...)
	}
	for _, path := range candidates {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		parseDotEnvLines(content)
		return
	}
}

func parseDotEnvLines(content []byte) {
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		index := strings.Index(line, "=")
		if index <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:index])
		value := strings.Trim(strings.TrimSpace(line[index+1:]), `"'`)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, value)
	}
}

func getEnvInt(key string, defaultValue int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return number
}

func Log(v ...interface{}) {
	loadRuntimeEnv()
	if DEBUGMODE != 1 {
		return
	}
	ServiceInfo(v...)
}
func FormatFloatA(v float64, n int) float64 {
	s := fmt.Sprintf("%."+strconv.Itoa(n)+"f", v)
	return GetFloat(s)
}
func WriteLog(filePath string, text interface{}) {

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		ServiceError("open log file failed:", err)
		return
	}
	//及时关闭file句柄
	defer file.Close()
	write := bufio.NewWriter(file)

	write.WriteString(GetJsonValue(text) + "\n")
	write.Flush()
}
func FormatFloat(v float64) float64 {
	s := fmt.Sprintf("%.2f", v)
	return GetFloat(s)
}
func CheckUserName(username string) bool {
	r := regexp.MustCompile(`[^0-9a-zA-Z_]`)
	return !r.Match([]byte(username))
}
func CheckPhone(phone string) bool {
	r := regexp.MustCompile(`[^0-9]`)
	return !r.Match([]byte(phone))
}
func CheckEmail(email string) bool {
	r := regexp.MustCompile(`^[0-9a-zA-Z_\-\.]+@[0-9a-zA-Z_\-]+\.[a-zA-Z\.]{2,9}$`)
	return r.Match([]byte(email))
	//return false
}
func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func Ip2Long(ip string) int { //ip转为整数
	t_arr := strings.Split(ip, ".")
	if len(t_arr) < 4 {
		return 0
	}
	rs := 0
	for i := 0; i < 4; i++ {
		//rs = rs + GetInt(t_arr[i])*int(math.Pow(float64(255), float64(3-i)))
		rs = rs | (GetInt(t_arr[i]) << ((3 - i) * 8))
	}
	return rs
}
func Long2Ip(n int) string { //整数转换为IP

	return fmt.Sprintf("%d.%d.%d.%d", (n>>24)&0xFF, (n>>16)&0xFF, (n>>8)&0xFF, n&0xFF)

}
func ConvertObjectByJson(from interface{}, to interface{}) error {
	fromJson, err := json.Marshal(from)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fromJson, to)
	return err
}
func GetJsonValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch v.(type) {
	case int:
		return fmt.Sprintf("%d", v.(int))
	case int64:
		return fmt.Sprintf("%d", v.(int64))
	case float32:
		return fmt.Sprintf("%.8f", v.(float64))
	case float64:
		return fmt.Sprintf("%.8f", v.(float64))
	case string:
		return v.(string)
	case []byte:
		return string(v.([]byte))

	default:
		s, err := json.Marshal(v)
		if err == nil {
			return string(s)
		}
		return ""
	}
}

func GetNow() int {
	return int(time.Now().Unix())
}
func GetInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	i := strings.Index(s, ".")
	if i > 0 {
		//如果被转换成了浮点则要去除浮点
		s = s[0:i]
	}
	i, err := strconv.Atoi(s)

	if err != nil {
		Log("convert error", err.Error())
		return 0
	}
	return i
}
func GetFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}
	return 0
}
func GetFilename(filename string) string {
	index := strings.Index(filename, ".")
	return filename[0:index]
}

func CreateCacheId(member_id int) string {
	return strconv.Itoa(member_id)
}
func GetDay() int {
	//返回当天的UNIX时间
	now := time.Now()
	y := now.Year()
	m := now.Month()
	d := now.Day()
	n := time.Date(y, m, d, 0, 0, 0, 0, time.Local).Unix()
	return int(n)
}

func IntTimeToString(unix int64) string {
	return time.Unix(unix, 0).Format("2006-01-02 15:04:05")
}
func TimeToint64(timeString string) int64 {
	t := strings.Split(timeString, " ")
	if len(t) < 2 {
		timeString += " 00:00:00"
	}
	loc, _ := time.LoadLocation("Local") //获取当地时区
	location, err := time.ParseInLocation("2006-01-02 15:04:05", timeString, loc)
	if err != nil {
		Log("TimeToint64", err, "\r")
	}
	return location.Unix()
}

//对长度不足n的数字前面补0
func Sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

func RandName() string {
	rnum := strconv.Itoa(rand.Intn(10000000)) + strconv.Itoa(GetNow())
	return Md5(rnum)
}
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	t := time.Now()
	hs := strings.Split(host, ":")
	year, month, day := t.Date()
	curtime := fmt.Sprintf("%d-%d-%d %d:%d", year, month, day, t.Hour(), t.Minute())
	auth := smtp.PlainAuth("", user, password, hs[0])
	tos := strings.Split(to, ",")
	header := make(map[string]string)
	header["From"] = user
	header["To"] = tos[0]
	header["Date"] = curtime
	header["Subject"] = subject
	header["Content-Type"] = "text/html;charset=UTF-8"

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	msg += "\r\n" + body
	conn, err := tls.Dial("tcp", host, &tls.Config{})
	if err != nil {
		ServiceError("smtp tls dial failed:", err)
		return err
	}
	co, err := smtp.NewClient(conn, hs[0])
	if err != nil {
		ServiceError("smtp client create failed:", err)
		return err
	}
	defer co.Close()
	if auth != nil {
		if ok, _ := co.Extension("AUTH"); ok {
			if err = co.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = co.Mail(user); err != nil {
		return err
	}
	for _, addr := range tos {
		if err = co.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := co.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	co.Quit()
	return nil
}
func Keys(obj map[string]interface{}) []string {
	rs := make([]string, 0)
	for k, _ := range obj {
		rs = append(rs, k)
	}
	return rs
}

func CalPageParam(total int, size int) (totalpage int, pagesize int) {
	if !(size > 0) {
		size = 30
	}
	totalPage := math.Ceil(float64(total) / float64(size))
	return int(totalPage), size
}

func CalNowOffset(page int, offset int) (limitnum int, offsetnum int) {
	if !(offset > 0) {
		offset = 30
	}
	if !(page > 0) {
		return 0, offset
	}
	limit := (page - 1) * offset
	return limit, offset
}

func WriteFile(path string, base64_image_content string) (url string, status bool) {
	b, _ := regexp.MatchString(`^data:\s*image\/(\w+);base64,`, base64_image_content)

	if !b {
		return "", false
	}
	re, _ := regexp.Compile(`^data:\s*image\/(\w+);base64,`)
	allData := re.FindAllSubmatch([]byte(base64_image_content), 2)
	fileType := string(allData[0][1]) //png ，jpeg 后缀获取

	base64Str := re.ReplaceAllString(base64_image_content, "")

	date := time.Now().Format("20060102")
	if ok := IsFileExist(path + "/" + date); !ok {
		pathArr := strings.Split(path+"/"+date, "/")
		var b string
		for _, p := range pathArr {
			b += p + "/"
			Log(b)
			if ok = IsFileExist(b); !ok {
				os.Mkdir(b, 06666)
			}
		}
	}

	var file string = path + "/" + date + "/" + strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(999999-100000)+100000) + "." + fileType

	byte, _ := base64.StdEncoding.DecodeString(base64Str)

	err := ioutil.WriteFile(file, byte, 0666)
	if err != nil {
		Log(err)
		return "", false
	}
	if file[0:9] == "../static" {
		file = file[len("../static"):]
	}
	return file, true
}

func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CheckPage(page interface{}, totalPage int) *db.DBValue {
	var value *db.DBValue
	var ok bool
	if value, ok = page.(*db.DBValue); !ok {
		return &db.DBValue{Value: 1}
	}
	if value.ToInt() > totalPage {
		return &db.DBValue{Value: totalPage}
	}
	return value
}
func HidCardNum(num string) string {
	if len(num) < 8 {
		return num
	}
	b := []byte(num)
	for i := 0; i < 4; i++ {
		b[i] = '*'
	}
	for i := len(num) - 4; i < len(num); i++ {
		b[i] = '*'
	}
	return string(b)
}
func HttpJsonPost(url string, headers map[string][]string, params map[string]interface{}) string {
	/*if params == nil {
		return ""
	}*/
	postdata := GetJsonValue(params)

	rq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(postdata))
	if err != nil {
		return ""
	}

	defer rq.Body.Close()

	rq.Header = headers
	//rq.Header.Add("Connection", "keep-alive")
	//rq.Header.Add("Content-Type","application/json;charset=UTF-8")
	//utils.Log(ioutil.ReadAll(rq.Body))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := new(http.Client)

	c.Transport = tr
	rp, err := c.Do(rq)
	if err != nil {
		Log(err.Error())
		return ""
	}
	defer rp.Body.Close()
	b, err := ioutil.ReadAll(rp.Body)
	if err != nil {
		return ""
	}
	return string(b)
	//rp, err := c.Post(url, "application/json;charset=UTF-8", strings.NewReader(postdata))
}

func CsvList(data interface{}) []map[string]string {
	rs := make([]map[string]string, 0)
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Array {
		return rs
	}
	code, err := json.Marshal(data)
	if err != nil {
		return rs
	}
	json.Unmarshal(code, &rs)
	return rs
}

func In(target string, str_array []string) bool {

	sort.Strings(str_array)

	index := sort.SearchStrings(str_array, target)

	if index < len(str_array) && str_array[index] == target {

		return true

	}

	return false

}

func Limit(p ...int) string {
	loadRuntimeEnv()
	page := p[0]
	limit := PAGE_LIMIT_SIZE
	if len(p) > 1 && p[1] > 0 {
		limit = p[1]
	}
	if page == 0 {
		page = 1
	}
	return fmt.Sprintf(" LIMIT %d, %d", (page-1)*limit, limit)
}

func Order(order string) string {
	if order == "" {
		return " "
	}
	if strings.Contains(order, "-") {
		order = strings.Replace(order, "-", "  ", 1)
	}

	return fmt.Sprintf(" ORDER BY %s", order)
}

/**
 *	检查是否存在于数组
 */
func Exists(list []string, find string) bool {
	index := arrays.ContainsString(list, find)
	return index > -1
}
func ClearSchar(originStr string) string {
	srcRunes := []rune(originStr)
	dstRunes := make([]rune, 0, len(srcRunes))
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 {
			continue
		}
		if c == 127 {
			continue
		}
		dstRunes = append(dstRunes, c)
	}
	return string(dstRunes)
}
func ReadAllFile(filepath string) string {
	f, e := os.OpenFile(filepath, os.O_RDONLY, os.ModeExclusive)
	if e != nil {
		return ""
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}
	return string(bs)
}
