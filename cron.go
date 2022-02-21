package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/robfig/cron/v3"
	"github.com/wcharczuk/go-chart/v2"
)

var month time.Time

/*
	Cron for pushing statistics chart to each user monthly.
*/
func PushChartMonthlyCron() {
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	//Every first day of month at 9:00
	c.AddFunc("05 18 21 */1 *", func() {
		now := time.Now()
		month = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

		users := GetMonthPushUsers(month)

		var wg sync.WaitGroup
		var m sync.Mutex

		//create 10 goroutines to psuh chart to users
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go pushChartWorker(&m, &wg, &users)
		}

		wg.Wait()
	})
	c.Run()

	defer c.Stop()
}

func pushChartWorker(m *sync.Mutex, wg *sync.WaitGroup, users *[]string) {
	for len(*users) > 0 {
		m.Lock()
		user := (*users)[0]
		(*users) = (*users)[1:]
		m.Unlock()

		pushMonthlyStatChartToUser(user)
	}

	wg.Done()
}

func pushMonthlyStatChartToUser(userId string) {
	stats := GetMonthStatDataByUser(month, userId)
	var chartData []chart.Value

	// format statistic data to chart data
	for _, v := range stats.Data {
		chartData = append(chartData, chart.Value{Label: fmt.Sprintf("%s $%d (%.2f%%)", v.Class, v.Sum, float64(v.Sum/stats.Total)), Value: float64(v.Sum)})
	}

	chart := GetChart(chartData)
	link := UploadToImgur(chart, os.Getenv("ImgurAccessToken"))
	if link != "" {
		if _, err := bot.PushMessage(userId, linebot.NewImageMessage(link, link)).Do(); err != nil {
			log.Println("Line reply error=", err)
		}
	}
}
