package main

import (
	"flag"
	"log"
	"net"

	"github.com/linkernetworks/oauth/src/app"
	"github.com/linkernetworks/oauth/src/app/config"
)

func main() {
	var configPath string
	var host string
	var port string
	var noSsl bool
	flag.StringVar(&host, "h", "", "hostname")
	flag.StringVar(&port, "p", "9096", "port")
	flag.StringVar(&configPath, "config", "config/default.json", "config file path")
	flag.BoolVar(&noSsl, "no-ssl", false, "disable https")
	flag.Parse()

	appConfig := config.Read(configPath)
	appService := app.NewServiceProviderFromConfig(*appConfig)
	bind := net.JoinHostPort(host, port)

	if noSsl {
		log.Println("SSL is disabled")
		log.Println(app.Start(bind, appService))
	} else {
		log.Println("SSL is able")
		log.Println(app.StartSsl(bind, appService))
	}
}
