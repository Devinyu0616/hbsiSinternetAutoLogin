package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func CheckServer() bool {
	timeout := time.Duration(5 * time.Second)
	t1 := time.Now()
	_, err := net.DialTimeout("tcp", "10.255.252.1:9090", timeout)
	fmt.Println("waist time :", time.Now().Sub(t1))
	if err != nil {
		fmt.Println("Site unreachable, error: ", err)
		return false
	}
	return true
}
func readConfig() url.Values {
	postData := url.Values{}
	f, err := os.Open("./data.conf")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	//建立缓冲区，把文件内容放到缓冲区中
	buf := bufio.NewReader(f)
	for {
		//遇到\n结束读取
		data, errR := buf.ReadBytes('\n')
		if errR != nil {
			if errR == io.EOF {
				break
			}
			fmt.Println(errR.Error())
		}
		dataSlice := strings.Split(string(data), ": ")

		//fmt.Println(dataSlice[0],dataSlice[1])

		str := strings.Replace(dataSlice[1], "\r\n", "", -1)

		postData.Add(dataSlice[0], str)
	}
	return postData
}
func httpDo(postData url.Values) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://10.255.252.1:9090/zportal/login/do",
		strings.NewReader(postData.Encode()))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")
	req.Header.Set("Host", "10.255.252.1:9090")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1")
	req.Header.Set("Referer", "http://10.255.252.1:9090/zportal/loginForWeb?wlanuserip=472fd2f123c4e6fd0e3fe16fd01d0744&wlanacname=8eca7b933ea6d97f271792fbc8bda5e9&ssid=4c38bf9aeaa52fc3&nasip=a7c9624b2f5351616880ead92b8afb1b&snmpagentip=&mac=2bd04c8d46ce13d58c8f39b11a6df8bb&t=wireless-v2&url=ac8f72321487d465d5ee224760897194&apmac=&nasid=8eca7b933ea6d97f271792fbc8bda5e9&vid=ed35b497bb97eef4&port=c148b2bc27d3020a&nasportid=74a80b37dc464ebf2d5ccc4a69bac5b1affa983c38ac302ae64ba4a265879dd7")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
func main() {
	var accountClock int = 0
	//先检查是否可以链接认证服务器，不能链接一直进行验证直到成功
	for {
		if accountClock == 101 {
			fmt.Println("connect failed maybe can't wifi")
			return
		}
		if CheckServer() {
			break
		} else {
			accountClock++
			time.Sleep(time.Second * 3)
		}
	}
	//读取配置文件
	data := readConfig()
	httpDo(data)
	//使用chan阻塞窗口
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("start")
	<-done
}
