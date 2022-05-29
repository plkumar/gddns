package main

import (
	"flag"
	"fmt"
	"strings"

	common "github.com/plkumar/gddns/common"
	config "github.com/plkumar/gddns/config"
	"github.com/plkumar/gddns/ddns"
)

func main() {
	fmt.Println("Google Dynamic DNS Client")
	standalone := flag.Bool("standalone", true, "Run in standalone mode.")
	configFile := flag.String("config", "gddns.yml", "configuration file path.")

	flag.Parse()

	if *standalone {
		y, err := config.GetConfig(*configFile)
		if err == nil {
			gd := ddns.GoogleDDNS{}
			for key, host := range y.Gddns {
				fmt.Println("Updating for: ", key)
				hostParams := host["params"]
				gd.SetHost(&hostParams)

				status, err := gd.UpdateDDNSIp()
				if err != nil {
					fmt.Println(err.Error())
				} else {

					if strings.Contains(status, "success") {
						fmt.Println("DNS Updated successfully.")
					} else if strings.Contains(status, "nochg") {
						fmt.Println("No Change")
					} else {
						// DNS Update failed, log and stop processing current host
						// TODO: Stop further DNS update attempts to ensure google is not blocking the client
						fmt.Println(status, common.DDNSStatusMap[status])
					}
				}
			}
		} else {
			fmt.Println("Error reading configuration :: ", err)
		}
	}
}
