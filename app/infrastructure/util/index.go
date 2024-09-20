package util

import (
	"math/rand"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	jsoniter "github.com/json-iterator/go"
	"github.com/sony/sonyflake"

	"jade-mes/app/infrastructure/log"
)

var mqHystrixConfig = hystrix.CommandConfig{
	Timeout:               5000,
	SleepWindow:           10000,
	MaxConcurrentRequests: 10000,
}

var httpHystrixConfig = hystrix.CommandConfig{
	Timeout:               15000,
	SleepWindow:           20000,
	MaxConcurrentRequests: 60000,
}

func init() {
	println("initing util...")
	rand.Seed(time.Now().UnixNano())

	flake = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Unix(0, 0),
	})

	JSON = jsoniter.Config{UseNumber: true}.Froze()

	Cache = newCacheClient("")

	hystrix.SetLogger(log.GetLogger())
	hystrix.Configure(map[string]hystrix.CommandConfig{
		circuitBreakerNamePostData:           httpHystrixConfig,
		circuitBreakerNameGetData:            httpHystrixConfig,
		circuitBreakerNameGetDataWithHeaders: httpHystrixConfig,
	})
}
