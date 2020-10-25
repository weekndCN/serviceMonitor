package alert

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/serviceMonitor/m/dingtalk"
	"github.com/serviceMonitor/m/k8s"
)

// Exception hanlder
func Exception(webhook string, rchan chan k8s.ExceptMsg) {
	robot := dingtalk.NewRobot(webhook)
	for {
		r := <-rchan
		limiter := getService(r.Service)
		if !limiter.Allow() {
			log.Println("limiter hijacked")
		} else {
			title := r.Service
			fmt.Println(r.Message)
			text := fmt.Sprintf("**<font color='#f60'>%s异常</font>** \n\n> **开始时间:** %s\n\n> **耗时:** %s\n\n> **错误信息:** %s\n\n", r.Service, r.Start.Format(time.UnixDate), r.End.Sub(r.Start), r.Message)
			atMobiles := []string{"13524353067"}
			isAtAll := false
			err := robot.SendMarkdown(title, text, atMobiles, isAtAll)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func exceptSet() int {
	interval := 60
	if os.Getenv("alert_interval") == "" {
		interval = 60
		log.Println("default change to 60")
	}

	interval, err := strconv.Atoi(os.Getenv("alert_interval"))
	if err != nil {
		log.Println("Set incorrect value.need int number,default change to 60")
	}

	return interval
}
