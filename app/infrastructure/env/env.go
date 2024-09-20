package env

import (
	"jade-mes/config"
)

type keys struct {
	JWT           string `mapstructure:"jwt"`
	IDTokenSign   string `mapstructure:"private"`
	IDTokenVerify string `mapstructure:"public"`
	Hash          string `mapstructure:"hash"`
}

// Keys stores all siging key for all algorithm
var Keys keys

// Issuer is for id_token
var Issuer string

// APIConfig ...
type APIConfig map[string]interface{}

// API ...
var API APIConfig

func init() {
	settings := config.GetConfig()

	settings.UnmarshalKey("keys", &Keys)

	Issuer = settings.GetString("issuer")

	API = settings.GetStringMap("api")
}
