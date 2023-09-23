package main

import (
	"github.com/stellar-payment/sp-gateway/cmd/webservice"
	"github.com/stellar-payment/sp-gateway/internal/component"
	"github.com/stellar-payment/sp-gateway/internal/config"
	"github.com/stellar-payment/sp-gateway/pkg/initutil"
)

var (
	buildVer  string = "unknown"
	buildTime string = "unknown"
)

func main() {
	config.Init(buildTime, buildVer)
	conf := config.Get()

	initutil.InitDirectory()

	logger := component.NewLogger(component.NewLoggerParams{
		ServiceName: conf.ServiceName,
		PrettyPrint: true,
	})

	webservice.Start(conf, logger)
}
