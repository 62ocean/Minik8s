package HTTPClient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Host string // ip and port
	http *http.Client
}

func CreateHTTPClient(serverHost string) *Client {
	ret := Client{
		Host: serverHost,
		// timeout: 5s
		http: &http.Client{Timeout: 5 * time.Second},
	}
	return &ret
}

func (c *Client) Get(url string) string {
	resp, err := c.http.Get(c.Host + url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}

func (c *Client) Post(url string, requestBody string) string {
	res, err := c.http.Post(c.Host+url, "application/json", strings.NewReader(string(requestBody)))

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return string(body)
}

func (c *Client) Del(url string) string {
	req, _ := http.NewRequest("DELETE", c.Host+url, nil)
	res, err1 := c.http.Do(req)
	if err1 != nil {
		fmt.Println(err1.Error())
		return ""
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(body)
}
