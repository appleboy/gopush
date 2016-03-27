package main

import (
	"flag"
	"github.com/appleboy/gopush/gopush"
	"log"
)

func main() {
	version := flag.Bool("v", false, "gopush version")
	confPath := flag.String("c", "", "yaml configuration file path for gopush")
	certificateKeyPath := flag.String("i", "", "iOS certificate key file path for gopush")
	apiKey := flag.String("k", "", "Android api key configuration for gopush")
	port := flag.String("p", "", "port number for gopush")

	flag.Parse()

	if *version {
		gopush.PrintGoPushVersion()
		return
	}

	var err error

	// set default parameters.
	gopush.PushConf = gopush.BuildDefaultPushConf()

	// load user define config.
	if *confPath != "" {
		gopush.PushConf, err = gopush.LoadConfYaml(*confPath)

		if err != nil {
			log.Printf("Unable to load yaml config file: '%v'", err)

			return
		}
	}

	if !gopush.PushConf.Ios.Enabled && !gopush.PushConf.Android.Enabled {
		log.Printf("Note: Please enable iOS or Android config in yaml config")

		return
	}

	if gopush.PushConf.Ios.Enabled {

		if *certificateKeyPath != "" {
			gopush.PushConf.Ios.PemKeyPath = *certificateKeyPath
		}

		if gopush.PushConf.Ios.PemKeyPath == "" {
			log.Println("iOS certificate path not define")

			return
		}
	}

	// check andorid api key exist
	if gopush.PushConf.Android.Enabled {

		if *apiKey != "" {
			gopush.PushConf.Android.ApiKey = *apiKey
		}

		if gopush.PushConf.Android.ApiKey == "" {
			log.Println("Android API Key not define")

			return
		}
	}

	// overwrite server port
	if *port != "" {
		gopush.PushConf.Core.Port = *port
	}

	gopush.InitAPNSClient()
	gopush.RunHTTPServer()
}