package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"time"

	config "github.com/plkumar/gddns/config"
)

func getIp() (string, error) {
	resp, err := http.Get("https://domains.google.com/checkip")
	if err == nil {
		scanner := bufio.NewScanner(resp.Body)
		if scanner.Scan() {
			return scanner.Text(), nil
		}
	}

	return "", errors.New("error fetching ip address")
}

func updateDnsIp(cfg config.Params) (string, error) {

	status := ""
	ip, err := getIp()
	if err != nil {
		fmt.Println(err.Error())
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", "https://domains.google.com/nic/update", nil)
	if err != nil {
		//fmt.Print("Got error %s", err.Error())
		return "", err
	}

	req.Header.Set("User-Agent", "Chrome/41.0 kumar.lakshman@gmail.com")
	req.SetBasicAuth(cfg.Username, cfg.Password)

	q := req.URL.Query()
	q.Add("hostname", cfg.Hostname)
	q.Add("myip", ip)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err == nil {
		defer resp.Body.Close()
		defer client.CloseIdleConnections()
		scanner := bufio.NewScanner(resp.Body)
		if scanner.Scan() {
			status = scanner.Text()
			fmt.Println(status)
		}
	}

	return status, nil
}

func main() {
	fmt.Println("Google Dynamic DNS Client")

	y, err := config.GetConfig()
	if err == nil {
		for key, host := range y.Gddns {
			fmt.Println(key)

			err := updateDnsIp(host[key])
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
