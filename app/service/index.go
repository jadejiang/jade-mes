package service

import (
	"time"

	"github.com/go-resty/resty/v2"
)

var timeout = 10 * time.Second
var client *resty.Client

const (
	// HeaderContentType ...
	HeaderContentType = "Content-Type"
	// ContentTypeJSON ...
	ContentTypeJSON = "application/json"
)

func init() {
	println("initing service...")

	client = resty.New()
	client.SetRetryCount(3)

	// ============================================================================
	// MQ Consumer
	// ============================================================================
}
