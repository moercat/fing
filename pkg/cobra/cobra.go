package cobra

import (
	"fmt"
	"time"
)

func Cobra() {

	go Ticker()

}

func Ticker() {
	ticker := time.NewTicker(time.Hour * 10)
	for range ticker.C {
		fmt.Println("这是一个定时任务")
	}
}
