package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type KRequest struct {
	Url      string
	PingTime float64
}

var defaultURL = "http://vnexpress.net/tin-tuc/the-gioi/obama-se-cho-trung-quoc-thay-hau-qua-tren-bien-dong-3462429.html"

func parseDataToArray(data []byte) []string {
	return strings.Split(
		string(data),
		"\n",
	)
}

func parseUrlToIpAndPort(url string) (string, string) {
	var tmp = strings.SplitN(url, string(' '), 2)
	return tmp[0], tmp[1]
}

func ping(ip string, port string, c chan KRequest) {
	var timeout = time.Duration(30 * time.Second)

	startAt := time.Now()
	host := fmt.Sprintf("%s:%s", ip, port)

	proxyUrl := &url.URL{Host: host}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: timeout,
	}

	response, err := client.Get(defaultURL)

	if err != nil {
		c <- KRequest{
			Url:      host,
			PingTime: float64(-1),
		}
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	delta := time.Now().UnixNano() - startAt.UnixNano()

	if strings.Contains(string(body), "Obama") {
		c <- KRequest{
			Url:      host,
			PingTime: float64(delta) / 1e9,
		}

	} else {
		c <- KRequest{
			Url:      host,
			PingTime: float64(-1),
		}
	}

}

func main() {
	// Read all adress from ip.txt
	data, err := ioutil.ReadFile("ip.txt")

	if err != nil {
		fmt.Println("Can't read file")
		return
	}

	arrOfUrls := parseDataToArray(data)
	queue := make(chan KRequest, 100)

	for _, url := range arrOfUrls {
		ip, port := parseUrlToIpAndPort(url)
		go ping(ip, port, queue)
	}

	for _, _ = range arrOfUrls {
		r := <-queue
		if r.PingTime > 1e-9 {
			fmt.Printf("%s %v\n", r.Url, r.PingTime)
		}
	}
}
