package util

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"jade-mes/config"

	"math/rand"
	"reflect"
	"regexp"
	"runtime"
	"time"

	"jade-mes/app/infrastructure/log"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
	"github.com/sony/sonyflake"
	"github.com/spf13/cast"
)

var _ = proto.Marshal
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
var flake *sonyflake.Sonyflake

// RandString is a helper function to generate random string
func RandString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandNumber ...
func RandNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// NewUUID generates global unique id by github.com/google/uuid
func NewUUID() (string, error) {
	id, err := uuid.NewUUID()

	return id.String(), err
}

// NewUUIDV4 ...
func NewUUIDV4() string {
	return uuid.New().String()
}

// NewGUID generates global unique id by https://github.com/rs/xid
func NewGUID() (string, error) {
	guid := xid.New().String()

	return guid, nil
}

// NewGUN generates global unique number by https://github.com/sony/sonyflake
func NewGUN() (uint64, error) {
	return flake.NextID()
}

// Hash is a helper function to hash string using sha256
func Hash(s string) string {
	h := sha256.New()

	h.Write([]byte(s))

	hash := hex.EncodeToString(h.Sum(nil))

	return hash
}

// Marshal converts a protobuf message to a URL legal string.
func Marshal(message proto.Message) (string, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

// Unmarshal decodes a protobuf message.
func Unmarshal(s string, message proto.Message) error {
	data, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, message)
}

// IsZeroOfUnderlyingType determines whether a value is zero value
func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// CleanZeroValues delete all zero value keys in map
func CleanZeroValues(x map[string]interface{}) {
	for k, v := range x {
		if IsZeroOfUnderlyingType(v) {
			delete(x, k)
		}
	}
}

// ConvertToMap ...
func ConvertToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if err := mapstructure.Decode(data, &result); err != nil {
		return nil, err
	}

	CleanZeroValues(result)

	return result, nil
}

// IncludesString 自制JS的includes方法
func IncludesString(array []string, str string) (contains bool) {
	contains = false
	for _, item := range array {
		if item == str {
			contains = true
			break
		}
	}
	return
}

// IncludesInt 自制JS的includes方法
func IncludesInt(array []int, num int) (contains bool) {
	contains = false
	for _, item := range array {
		if item == num {
			contains = true
			break
		}
	}
	return
}

// IncludesInt64 自制JS的includes方法
func IncludesInt64(array []int64, num int64) (contains bool) {
	contains = false
	for _, item := range array {
		if item == num {
			contains = true
			break
		}
	}
	return
}

// GetFunctionName ...
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// TimeTrack ...
func TimeTrack(start time.Time, logger *log.Logger) {
	if config.GetConfig().GetBool("release_mode") {
		return
	}
	elapsed := time.Since(start)

	// Regex to extract just the function name (and not the module path).
	funcNameFunc := regexp.MustCompile(`^.*\.(.*)$`)
	runtimeFunc := regexp.MustCompile(`^.*\/(.*)$`)
	goroutineFunc := regexp.MustCompile(`created by .*\/(.*)\s*.*\:(\d*)`)

	// Skip this function, and fetch the PC and file for its parent.
	funcPC, _, _, _ := runtime.Caller(1)
	// Retrieve a function object this functions parent.
	funcName := funcNameFunc.ReplaceAllString(runtime.FuncForPC(funcPC).Name(), "$1")

	// get caller info
	callerPC, _, callerLine, _ := runtime.Caller(2)
	callerName := runtimeFunc.ReplaceAllString(runtime.FuncForPC(callerPC).Name(), "$1")
	var parentCallerName string
	var parentCallerLine int

	// 针对在goroutine里执行的函数，通过调用链获取信息
	if callerName == "runtime.goexit" {
		trace := make([]byte, 1<<12)
		runtime.Stack(trace, false)
		stackTrace := string(trace)
		results := goroutineFunc.FindStringSubmatch(stackTrace)
		if len(results) >= 3 {
			parentCallerName = results[1]
			parentCallerLine = cast.ToInt(results[2])
		}
	} else {
		parentCallerPC, _, callerLine, _ := runtime.Caller(3)
		parentCallerLine = callerLine
		parentCallerName = runtimeFunc.ReplaceAllString(runtime.FuncForPC(parentCallerPC).Name(), "$1")
	}

	callerInfo := fmt.Sprintf("%s:%v - %s:%v", parentCallerName, parentCallerLine, callerName, callerLine)

	logger.Info("func execution elapsed",
		log.String("funcName", funcName), log.String("caller", callerInfo), log.String("elapsed", elapsed.String()),
	)
}

//手机号脱敏
func HideStar(str string) (result string) {
	if str == "" {
		return "***"
	}
	if strings.Contains(str, "@") {
		res := strings.Split(str, "@")
		if len(res[0]) < 3 {
			resString := "***"
			result = resString + "@" + res[1]
		} else {
			res2 := Substr2(str, 0, 3)
			resString := res2 + "***"
			result = resString + "@" + res[1]
		}
		return result
	} else {
		reg := `^1[0-9]\d{9}$`
		rgx := regexp.MustCompile(reg)
		mobileMatch := rgx.MatchString(str)
		if mobileMatch {
			result = Substr2(str, 0, 3) + "****" + Substr2(str, 7, 11)
		} else {
			nameRune := []rune(str)
			lens := len(nameRune)

			if lens <= 1 {
				result = "***"
			} else if lens == 2 {
				result = string(nameRune[:1]) + "*"
			} else if lens == 3 {
				result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
			} else if lens == 4 {
				result = string(nameRune[:1]) + "**" + string(nameRune[lens-1:lens])
			} else if lens > 4 {
				result = string(nameRune[:2]) + "***" + string(nameRune[lens-2:lens])
			}
		}
		return
	}
}

func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}
