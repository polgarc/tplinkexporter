package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

type TPLINKSwitchClient interface {
	GetPortStats() ([]portStats, error)
	GetHost() string
	// Login()
}

type TPLINKSwitch struct {
	host       string
	username   string
	password   string
	httpClient *http.Client
	loggedIn   bool
}

type portStats struct {
	State      int
	LinkStatus int
	PktCount   map[string]int
}

func (client *TPLINKSwitch) GetHost() string {
	return client.host
}

func (client *TPLINKSwitch) Login() error {
	requestData := url.Values{"username": {client.username}, "password": {client.password}, "logon": {"Login"}}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/logon.cgi", client.host), strings.NewReader(requestData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(body), "var logonInfo = new Array(\n0,") {
		client.loggedIn = true
		log.Printf("Logged in to %s as %s\n", client.host, client.username)
		return nil
	} else {
		client.loggedIn = false
		return fmt.Errorf("error logging in to %s as %s", client.host, client.username)
	}
}

func (client *TPLINKSwitch) fetchPortStats() (string, error) {
	resp, err := client.httpClient.Get(fmt.Sprintf("http://%s/PortStatisticsRpm.htm", client.host))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (client *TPLINKSwitch) GetPortStats() ([]portStats, error) {
	type allInfo struct {
		State      []int
		LinkStatus []int
		Pkts       []int
	}

	if !client.loggedIn {
		err := client.Login()
		if err != nil {
			return nil, err
		}
	}

	body, err := client.fetchPortStats()
	if err != nil {
		return nil, err
	} else if strings.Contains(string(body), "var logonInfo = new Array(\n0,") {
		// logged out or session expired, try to login again
		client.Login()
		body, err = client.fetchPortStats()
		if err != nil {
			return nil, err
		}
	}

	// fmt.Println(string(body))
	var jbody string = strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				string(body), "link_status", `"linkStatus"`),
			"state", `"State"`),
		"pkts", `"Pkts"`)
	// fmt.Println(string(jbody))
	res := regexp.MustCompile(`all_info = ({[^;]*});`).FindStringSubmatch(jbody)
	if res == nil {
		// fmt.Println(jbody)
		return nil, errors.New("unexpected response for port statistics http call: " + jbody)
	}
	// fmt.Println(res[1])
	var jparsed allInfo
	json.Unmarshal([]byte(res[1]), &jparsed)
	// fmt.Println(jparsed)
	var portsInfos []portStats
	portcount := len(jparsed.State)
	for i := 0; i < portcount; i++ {
		var portInfo portStats
		portInfo.State = jparsed.State[i]
		portInfo.LinkStatus = jparsed.LinkStatus[i]
		if portInfo.State == 1 {
			portInfo.PktCount = make(map[string]int)
			portInfo.PktCount["TxGoodPkt"] = jparsed.Pkts[4*i]
			portInfo.PktCount["TxBadPkt"] = jparsed.Pkts[4*i+1]
			portInfo.PktCount["RxGoodPkt"] = jparsed.Pkts[4*i+2]
			portInfo.PktCount["RxBadPkt"] = jparsed.Pkts[4*i+3]
		}
		portsInfos = append(portsInfos, portInfo)
	}
	// fmt.Println(portsInfos)
	log.Printf("Fetched port statistics")
	return portsInfos, nil
}

/*
sample output of PortStatisticsRpm.htm call:
<script>
var max_port_num = 8;
var port_middle_num  = 16;
var all_info = {
state:[1,1,1,1,1,1,1,1,0,0],
link_status:[6,6,0,6,0,0,0,5,0,0],
pkts:[1901830310,0,1338131260,33254,4291149014,0,2311488878,564,0,0,0,0,1814018004,0,33552310,0,0,0,0,0,0,0,0,0,0,0,0,0,1678459124,0,1866169392,0,0,0]
};
var tip = "";
</script>
*/

func NewTPLinkSwitch(host string, username string, password string) (*TPLINKSwitch, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Jar: jar,
	}
	return &TPLINKSwitch{
		host:       host,
		username:   username,
		password:   password,
		httpClient: httpClient,
	}, nil
}
