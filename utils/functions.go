package util

import (
	"check_vpn/mylog"
	"context"
	"errors"
	"fmt"
	"github.com/commander-cli/cmd"
	jsoniter "github.com/json-iterator/go"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var Ctx context.Context

func GetCtx() context.Context {
	return Ctx
}
func SetCtx(ctx context.Context) {
	Ctx = ctx
}
func Ping(host string) *probing.Statistics {
	pinger, err := probing.NewPinger(host)
	mylog.Der(err)
	pinger.Count = 3
	pinger.Timeout = time.Second
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		mylog.Logf("pingerr:%v,\nip:%v", err, host)
		panic(err)
	}
	return pinger.Statistics()
}
func LocalRunCmd(cm string) (string, error) {
	//mylog.Logf("local run:%v", cm)
	c := cmd.NewCommand(cm)

	err := c.Execute()
	if err != nil {
		return "", err
	}
	if c.Stderr() != "" {
		return "", &MyError{c.Stderr()}
	}
	return c.Stdout(), nil
}
func DelWgInterface(ret string) {
	rt := strings.Split(ret, "\n")
	for _, i := range rt {
		if strings.Index(i, "interface:") != -1 {
			//fmt.Println(i)
			facename := strings.Trim(strings.Replace(i, "interface:", "", -1), " \n\r")
			cmd := "wg-quick down " + facename
			LocalRunCmd("touch /root/" + facename + ".conf")
			LocalRunCmd(cmd)
			//fmt.Println(err)
			//fmt.Println(ret)

		}
	}
}
func CheckWg(uip, upk, spk, aip, eip string) (ret_long float32) {
	//flag.Parse()
	ip, err := LocalRunCmd("curl -s 'https://ip-score.com/ip'")
	mylog.Logf("pre-ip:%v", ip)
	wgret, err := LocalRunCmd("wg")
	if wgret != "" {
		DelWgInterface(wgret)
	}
	confg := "echo \"[Interface]\nAddress = %v\nSaveConfig = true\nListenPort = %v\nFwMark = 0xca6c\nPrivateKey = %v\n\n[Peer]\nPublicKey = %v\nAllowedIPs = %v\nEndpoint = %v\">/etc/wireguard/wg0.conf"
	_, err = LocalRunCmd(fmt.Sprintf(confg, uip, "7867", upk, spk, aip, eip))
	if err != nil {
		log.Fatalf("failed LocalRunCmd: %v", err)
	}
	_, err = LocalRunCmd("wg-quick up wg0")
	if err != nil {
		//log.Fatalf("failed up wg0: %v", err)
	}
	start := time.Now().UnixMilli()
	ip, err = LocalRunCmd("curl -s 'https://ip-score.com/ip'")
	wgret, err = LocalRunCmd("wg")
	ret_long = float32(time.Now().UnixMilli()-start) / float32(1000.0000000000000)

	if wgret != "" {
		DelWgInterface(wgret)
	}
	mylog.Logf("new-ip:%v", ip)
	return
}
func CheckSs(addr, cipher, password string) float32 {

	// 配置SS代理
	var key []byte

	var err error

	udpAddr := addr

	ciph, err := core.PickCipher(cipher, key, password)
	if err != nil {
		log.Fatal(err)
	}
	go socksLocal("0.0.0.0:777", udpAddr, ciph.StreamConn)
	ip, err := LocalRunCmd("curl -s 'https://ip-score.com/ip'")
	mylog.Logf("pre-ip:%v", ip)
	<-time.After(2 * time.Second)

	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:777", nil, proxy.Direct)
	if err != nil {
		fmt.Printf("Failed to create SOCKS dialer: %v\n", err)
		return 0
	}

	// 创建HTTP客户端并设置Transport为SOCKS拨号器
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}
	mylog.Logf("start:time:%v", time.Now().UnixMilli())
	// 发送HTTP GET请求并打印响应
	start := time.Now().UnixMilli()

	resp, err := httpClient.Get("https://ip-score.com/ip")
	mylog.Logf("stop:time:%v", time.Now().UnixMilli())
	if err != nil {
		fmt.Printf("Failed to send HTTP request via SOCKS proxy: %v\n", err)
		return 0
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	ret_long := float32(time.Now().UnixMilli()-start) / float32(1000.0000000000000)

	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return 0
	}
	//fmt.Printf("Response body: %s\n", string(body))
	return ret_long
}

// Create a SOCKS server listening on addr and proxy to server.
func socksLocal(addr, server string, shadow func(net.Conn) net.Conn) {
	mylog.Logf("SOCKS proxy %s <-> %s", addr, server)
	tcpLocal(addr, server, shadow, func(c net.Conn) (socks.Addr, error) { return socks.Handshake(c) })
}
func tcpLocal(addr, server string, shadow func(net.Conn) net.Conn, getAddr func(net.Conn) (socks.Addr, error)) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		mylog.Logf("failed to listen on %s: %v", addr, err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			mylog.Logf("failed to accept: %s", err)
			continue
		}

		go func() {
			defer c.Close()
			tgt, err := getAddr(c)
			if err != nil {

				// UDP: keep the connection until disconnect then free the UDP socket
				if err == socks.InfoUDPAssociate {
					buf := make([]byte, 1)
					// block here
					for {
						_, err := c.Read(buf)
						if err, ok := err.(net.Error); ok && err.Timeout() {
							continue
						}
						mylog.Logf("UDP Associate End.")
						return
					}
				}

				mylog.Logf("failed to get target address: %v", err)
				return
			}

			rc, err := net.Dial("tcp", server)
			if err != nil {
				mylog.Logf("failed to connect to server %v: %v", server, err)
				return
			}
			defer rc.Close()
			//if config.TCPCork {
			//	rc = timedCork(rc, 10*time.Millisecond, 1280)
			//}
			rc = shadow(rc)

			if _, err = rc.Write(tgt); err != nil {
				mylog.Logf("failed to send target address: %v", err)
				return
			}

			mylog.Logf("proxy %s <-> %s <-> %s", c.RemoteAddr(), server, tgt)
			if err = relay(rc, c); err != nil {
				mylog.Logf("relay error: %v", err)
			}
		}()
	}
}
func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()
	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}

var Myjson = jsoniter.ConfigCompatibleWithStandardLibrary

type MyError struct {
	Name string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("%v", e.Name)
}
