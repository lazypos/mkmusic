package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	url       = `http://wangguan.qianyifu.com:8881/gateway/pay.asp?userid=49497&orderid=aaa110086&money=1&hrefurl=http://lazypos.pw:51888/pay_verify&url=http://lazypos.pw:51888/pay_verify&bankid=weixin&sign=%s&ext=`
	sig       = `userid=49497&orderid=aaa110086&bankid=weixin&keyvalue=nGX1MqFtet0sVAwzj7RYt5Jph4Mu5Kh1d6D0EuQx`
	urlreturn = `lazypos.pw:51888/pay_verify?returncode=1&orderid=10086&money=5&sign=123`
)

func main() {
	md5sig := fmt.Sprintf(`%x`, md5.Sum([]byte(sig)))
	log.Println(fmt.Sprintf(url, md5sig))

	return

	resp, err := http.Get(fmt.Sprintf(url, md5sig))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(content))
}
