package pollers

import (
	"encoding/json"
	"io"
	"log"
	"main/core"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	jabbacore "github.com/johnjones4/Jabba/core"
)

const (
	ipV4Endpoint = "https://api.ipify.org?format=json"
	ipV6Endpoint = "https://api64.ipify.org?format=json"
)

type INetPoller struct {
	lastIPv4Address string
	lastIPv6Address string
	lastSuccess     time.Time
}

type ipResponse struct {
	IP string `json:"ip"`
}

func NewINetPoller() *INetPoller {
	return &INetPoller{}
}

func makeCall(endpoint string) (string, error) {
	res, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var resp ipResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	return resp.IP, nil
}

func isTCPError(err error) bool {
	if _, ok := err.(net.Error); ok {
		log.Println("Error connecting")
		return true
	}
	return false
}

func logINetDown(u jabbacore.Upstream) error {
	log.Println("Logging internet as down")
	jEvent := jabbacore.Event{
		EventVendorType: "inet",
		EventVendorID:   uuid.NewString(),
		VendorInfo: map[string]interface{}{
			"ipv4": "",
			"ipv6": "",
		},
		Created:  time.Now().UTC(),
		IsNormal: false,
	}
	err := u.LogEvent(&jEvent)
	if err != nil {
		return err
	}
	return nil
}

func (p *INetPoller) runALoop(u jabbacore.Upstream) error {
	var (
		ipv4 string
		ipv6 string
		err  error
	)
	log.Println("Checking internet status ...")
	for attempts := 0; attempts < 10 && ipv4 == "" && ipv6 == ""; attempts++ {
		ipv4, err = makeCall(ipV4Endpoint)
		if err == nil {
			ipv6, err = makeCall(ipV6Endpoint)
		}
		if err != nil {
			time.Sleep(time.Second * time.Duration(attempts+1))
		}
	}

	if err != nil {
		if isTCPError(err) {
			p.lastSuccess = time.Time{}
			p.lastIPv4Address = ""
			p.lastIPv6Address = ""
			return logINetDown(u)
		} else {
			return err
		}
	}

	log.Printf("Received addresses: %s, %s\n", ipv4, ipv6)

	if (p.lastIPv4Address == ipv4 && p.lastIPv6Address == ipv6) && p.lastSuccess.After(time.Now().UTC().Add(-core.Day)) {
		log.Println("Internet info unchanged from last loop")
		return nil
	}

	log.Println("Logging change in internet status")

	jEvent := jabbacore.Event{
		EventVendorType: "inet",
		EventVendorID:   uuid.NewString(),
		VendorInfo: map[string]interface{}{
			"ipv4": ipv4,
			"ipv6": ipv6,
		},
		Created:  time.Now().UTC(),
		IsNormal: true,
	}
	err = u.LogEvent(&jEvent)
	if err != nil {
		return err
	}

	p.lastIPv4Address = ipv4
	p.lastIPv6Address = ipv6
	p.lastSuccess = time.Now().UTC()

	return nil
}

func (p *INetPoller) Setup() error {
	return nil
}

func (p *INetPoller) Poll(u jabbacore.Upstream) {
	for {
		err := p.runALoop(u)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Minute)
	}
}
