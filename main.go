package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	common "github.com/plkumar/gddns/common"
	config "github.com/plkumar/gddns/config"
	"github.com/plkumar/gddns/ddns"
	"github.com/takama/daemon"
)

const (
	name        = "gddns"
	description = "Google Dynamic DNS Client Daemon"
)

var dependencies = []string{"dummy.service"}

var stdlog, errlog *log.Logger

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage(configFile *string) (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	timer1 := time.NewTimer(5 * time.Minute)
	for {
		select {
		case <-timer1.C:
			updateIP(configFile)
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)

			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

func updateIP(configFile *string) {

	y, err := config.GetConfig(*configFile)
	if err == nil {
		gd := ddns.GoogleDDNS{}
		for key, host := range y.Gddns {
			stdlog.Println("Updating for: ", key)
			hostParams := host["params"]
			gd.SetHost(&hostParams)

			status, err := gd.UpdateDDNSIp()
			if err != nil {
				errlog.Println(err.Error())
			} else {

				if strings.Contains(status, "success") {
					stdlog.Println("DNS Updated successfully.")
				} else if strings.Contains(status, "nochg") {
					stdlog.Println("No Change")
				} else {
					// DNS Update failed, log and stop processing current host
					// TODO: Stop further DNS update attempts to ensure google is not blocking the client
					stdlog.Println(status, common.DDNSStatusMap[status])
				}
			}
		}
	} else {
		errlog.Println("Error reading configuration :: ", err)
	}

}

func main() {
	//stdlog.Println("Google Dynamic DNS Client")

	standalone := flag.Bool("standalone", false, "Run in standalone mode.")
	configFile := flag.String("config", "gddns.yml", "configuration file path.")

	flag.Parse()

	if *standalone {
		updateIP(configFile)
	} else {

		srv, err := daemon.New(name, description, daemon.SystemDaemon, dependencies...)
		if err != nil {
			errlog.Println("Error: ", err)
			os.Exit(1)
		}
		service := &Service{srv}
		status, err := service.Manage(configFile)
		if err != nil {
			errlog.Println(status, "\nError: ", err)
			os.Exit(1)
		}
		stdlog.Println(status)
	}
}
