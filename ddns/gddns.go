package ddns

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/plkumar/gddns/common"
	"github.com/plkumar/gddns/config"
)

type GoogleDDNS struct {
	HostConfig config.Params
	tempFile   string
}

func (g *GoogleDDNS) SetHost(cfg *config.Params) {
	g.HostConfig = *cfg
}

func (g *GoogleDDNS) writeCurrentIP(tmpFile string, ip string) error {
	f, err := os.CreateTemp("", tmpFile)
	if err == nil {
		defer f.Close()
		f.WriteString(ip)
		g.tempFile = f.Name()
		return nil
	} else {
		return err
	}
}

func (g *GoogleDDNS) readIPFromTmpFile() (string, error) {
	data, err := os.ReadFile(g.tempFile)
	if err == nil {
		return string(data), nil
	} else {
		return "", err
	}
}

func (g *GoogleDDNS) checkIPChanged(tmpFile string, ip string) bool {
	lastIPAddress, err := g.readIPFromTmpFile()
	if err == nil {
		if lastIPAddress == ip {
			return false
		}
	}
	return true
}

func (g *GoogleDDNS) GetIP() (string, error) {

	resp, err := http.Get(common.GOOGLE_URL_IP_CHECK)
	if err == nil {
		scanner := bufio.NewScanner(resp.Body)
		if scanner.Scan() {
			return scanner.Text(), nil
		}
	}

	return "", errors.New("error fetching ip address")
}

func (g *GoogleDDNS) UpdateDDNSIp() (string, error) {

	status := ""
	ip, err := g.GetIP()
	if err != nil {
		fmt.Println(err.Error())
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", common.GOOGLE_URL_DDNS_UPDATE, nil)
	if err != nil {
		//fmt.Print("Got error %s", err.Error())
		return "", err
	}

	req.Header.Set("User-Agent", "Chrome/41.0 kumar.lakshman@gmail.com")
	req.SetBasicAuth(g.HostConfig.Username, g.HostConfig.Password)

	q := req.URL.Query()
	q.Add("hostname", g.HostConfig.Hostname)
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
