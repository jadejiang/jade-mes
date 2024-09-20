package util

import "github.com/afex/hystrix-go/hystrix"

const (
	circuitBreakerNamePostData           = "PostData"
	circuitBreakerNameGetData            = "GetData"
	circuitBreakerNameGetDataWithHeaders = "GetDataWithHeaders"
)

// CircuitBreakerDo ...
func CircuitBreakerDo(name string, run func() error, fallback func(error) error) error {
	return hystrix.Do(name, run, fallback)
}
