package client

import (
	"context"
	"github.com/romanornr/Bitmex-referral-analyzer/config"
	"github.com/spf13/viper"
	"github.com/zmxv/bitmexgo"
	"log"
)

var c config.Conf
var apiKey string
var apiSecret string

var  instance *bitmexgo.APIClient

// Using the singleton design pattern to check if an instance already exist
// if not, only then create a new one
func GetInstance() (context.Context, *bitmexgo.APIClient) {
	if instance != nil {
		return nil, instance
	}

	config.GetViperConfig()

	// Get your API key/secret pair at https://www.bitmex.com/app/apiKeys
	apiKey = viper.GetString("api_key")
	apiSecret = viper.GetString("api_secret")

	var err error
	// Create an authentication context
	auth := bitmexgo.NewAPIKeyContext(apiKey, apiSecret)

	// Create a shareable API client instance
	instance = bitmexgo.NewAPIClient(bitmexgo.NewConfiguration())
	if err != nil {
		log.Fatal(err)
	}

	return auth, instance
}