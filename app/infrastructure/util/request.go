package util

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"jade-mes/app/infrastructure/log"
)

var httpClient *resty.Client

func init() {
	httpClient = resty.New()
	httpClient.SetRetryCount(3)

	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	t := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          800,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   400,
	}

	httpClient.SetTransport(otelhttp.NewTransport(t))
}

// PostData ...
func PostData(ctx context.Context, url string, body map[string]interface{}, isForm bool) (string, error) {
	request := httpClient.R()
	request.SetContext(ctx)
	if isForm {
		formData := make(map[string]string, len(body))

		for k, v := range body {
			formData[k] = fmt.Sprintf("%v", v)
		}
		request = request.SetFormData(formData)
	} else {
		request = request.SetBody(body)
	}

	var response *resty.Response
	err := CircuitBreakerDo(circuitBreakerNamePostData, func() error {
		res, err := request.Post(url)
		response = res

		return err
	}, nil)

	var rawResponse string
	var latency int64
	if response != nil {
		rawResponse = response.String()
		latency = response.Time().Milliseconds()
	}

	logger := log.With()
	if deviceID, exists := body["deviceId"]; exists {
		logger = logger.With(log.Any("deviceId", deviceID))
	} else if deviceID, exists := body["virtualId"]; exists {
		logger = logger.With(log.Any("deviceId", deviceID))
	}
	if orderID, exists := body["orderId"]; exists {
		logger = logger.With(log.Any("orderId", orderID))
	}
	if orderNo, exists := body["orderNo"]; exists {
		logger = logger.With(log.Any("orderNo", orderNo))
	}

	logger.Debug("postData response",
		log.String("url", url),
		log.Reflect("postData", body),
		log.String("rawResponse", rawResponse),
		log.Int64("latency", latency),
		log.Err(err),
	)

	return rawResponse, err
}

// GetData ...
func GetData(ctx context.Context, url string, data map[string]interface{}) (string, error) {
	query := make(map[string]string)
	for key, value := range data {
		query[key] = fmt.Sprintf("%v", value)
	}

	var response *resty.Response
	err := CircuitBreakerDo(circuitBreakerNameGetData, func() error {
		request := httpClient.R()
		res, err := request.SetContext(ctx).SetQueryParams(query).Get(url)
		response = res

		return err
	}, nil)

	var rawResponse string
	var latency int64
	if response != nil {
		rawResponse = response.String()
		latency = response.Time().Milliseconds()
	}

	logger := log.With()
	if deviceID, exists := data["deviceId"]; exists {
		logger = logger.With(log.Any("deviceId", deviceID))
	} else if deviceID, exists := data["virtualId"]; exists {
		logger = logger.With(log.Any("deviceId", deviceID))
	}
	if orderID, exists := data["orderId"]; exists {
		logger = logger.With(log.Any("orderId", orderID))
	}
	if orderNo, exists := data["orderNo"]; exists {
		logger = logger.With(log.Any("orderNo", orderNo))
	}
	logger.Debug("getData response",
		log.String("url", url),
		log.Reflect("requestQuery", query),
		log.String("rawResponse", rawResponse),
		log.Int64("latency", latency),
		log.Err(err),
	)

	return rawResponse, err
}

// GetDataWithHeaders ...
func GetDataWithHeaders(url string, data map[string]interface{}, headers map[string]string) (string, error) {
	query := make(map[string]string)
	for key, value := range data {
		query[key] = fmt.Sprintf("%v", value)
	}

	var response *resty.Response
	err := CircuitBreakerDo(circuitBreakerNameGetDataWithHeaders, func() error {
		res, err := httpClient.R().SetQueryParams(query).SetHeaders(headers).Get(url)
		response = res

		return err
	}, nil)

	var rawResponse string
	var latency int64
	if response != nil {
		rawResponse = response.String()
		latency = response.Time().Milliseconds()
	}

	logger := log.With()
	if deviceID, exists := data["deviceId"]; exists {
		logger = logger.With(log.Any("deviceId", deviceID))
	} else if deviceID, exists := data["virtualId"]; exists {
		logger = logger.With(log.Any("deviceId", deviceID))
	}
	if orderID, exists := data["orderId"]; exists {
		logger = logger.With(log.Any("orderId", orderID))
	}
	if orderNo, exists := data["orderNo"]; exists {
		logger = logger.With(log.Any("orderNo", orderNo))
	}

	logger.Debug("GetDataWithHeaders response",
		log.String("url", url),
		log.Reflect("requestHeaders", headers),
		log.Reflect("requestQuery", query),
		log.String("rawResponse", rawResponse),
		log.Int64("latency", latency),
		log.Err(err),
	)

	return rawResponse, err
}
