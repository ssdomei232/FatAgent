package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.mmeiblog.cn/mei/FatAgent/pkg"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()

	c.AddFunc("@every 1m", func() {
		pkg.NewACController("/dev/tty1")
	})

	c.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig // 阻塞直到收到信号
	fmt.Printf("Hava a good Day!")
	c.Stop() // 关闭 cron
}
