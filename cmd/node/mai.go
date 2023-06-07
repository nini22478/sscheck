package main

import (
	"check_vpn/dbs"
	"check_vpn/dbs/model"
	"check_vpn/mylog"
	util "check_vpn/utils"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	util.SetCtx(context.Background())
	//TestGet()

	flag.Int64Var(&util.DoId, "nid", 0, "节点id")
	flag.Parse()
	if util.DoId == 0 {
		flag.Usage()
		os.Exit(1)
	}
	mods, _ := dbs.GetAllNode()
	for _, mod := range mods {
		go DoCheckMi(*mod)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
func TestGet() {
	util.GetWgRet("http://127.0.0.1:21231/www/test_lines", "aes-128-ecb", "d7eff65a678fe32a")
}
func DoCheckMi(node model.CheckNode) {
	interval := 1 * time.Second

	// 创建一个通道用于接收停止信号
	done := make(chan bool)

	// 在一个 goroutine 中执行定时任务
	go func() {
		for {
			select {
			case <-done:
				// 接收到停止信号，退出定时任务
				return
			default:
				node, err := dbs.GetOneNode(node.ID)
				if err != nil {
					mylog.Logf("%v", err)
					//done <- true
					interval := 1 * time.Second
					time.Sleep(interval)

					continue
				}
				// 执行定时任务的逻辑
				fmt.Println("定时任务执行时间:", time.Now())
				interval = time.Duration(*node.LimitWait) * time.Minute
				fmt.Println("定时任务jiange:", interval)

				dbs.DoCheck(*node)
				// 等待一段时间后执行下一次任务
				time.Sleep(interval)
			}
		}
	}()

	// 发送停止信号，通知定时任务停止
	//done <- true

}
