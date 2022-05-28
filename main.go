package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

func updateDnsIp() error {

	ip, err := getIp()
	if err == nil {
		fmt.Println("IP :", ip)
	}

	data := url.Values{
		"hostname": {"home.peethani.me"},
		"myip":     {},
	}

	resp, err := http.PostForm("", data)
}

func main() {
	fmt.Println("Google Dynamic DNS Client")

}
