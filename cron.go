package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/robfig/cron/v3"
	"github.com/wcharczuk/go-chart/v2"
)

func PushChartMonthlyCron() {
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	c.AddFunc("*/5 * * * *", func() {
		userId := os.Getenv("TestUserId")
		now := time.Now()
		thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		stats := GetMonthStatDataByUser(thisMonth, userId)

		var chartData []chart.Value

		for _, v := range stats {
			chartData = append(chartData, chart.Value{Label: fmt.Sprintf("%s $%d", v.Class, v.Total), Value: float64(v.Total)})
		}

		chart := GetChart(chartData)
		link := UploadToImgur(chart, os.Getenv("ImgurAccessToken"))
		if link != "" {
			if _, err := bot.PushMessage(userId, linebot.NewImageMessage(link, link)).Do(); err != nil {
				log.Println("Line reply error=", err)
			}
		}
	})
	c.Run()

	defer c.Stop()
}
