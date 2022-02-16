package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/wcharczuk/go-chart/v2"
)

func ServeLineBot() error {
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	return http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				log.Println(message.Text)

				if err != nil {
					log.Println("Quota err:", err)
				}

				switch operation, operationData := DetermineOperation(message.Text); operation {
				case Error:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("錯誤的輸入格式，請重新輸入！")).Do(); err != nil {
						log.Println("Line reply error=", err)
					}
				case KeepRecord:
					record := ConvertToRecord(*operationData)
					record.UserID = event.Source.UserID
					if record.Save() {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("成功儲存！\n%s-%d-%s", record.Class, record.Cost, record.Memo))).Do(); err != nil {
							log.Println("Line reply error=", err)
						}
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("儲存資料錯誤，請稍後再試！")).Do(); err != nil {
							log.Println("Line reply error=", err)
						}
					}
				case GetRecord:
					records := GetLastRecords(uint(operationData.Number))
					responseContainer := GetListRecordResponse(records)

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("Record list", responseContainer)).Do(); err != nil {
						log.Println("Line reply error=", err)
					}
				case GetStatistic:
					stats := GetStatData()
					var chartData []chart.Value

					for _, v := range stats {
						chartData = append(chartData, chart.Value{Label: v.Class, Value: float64(v.Sum)})
					}

					chart := GetChart(chartData)
					link := UploadToImgur(chart, os.Getenv("ImgurAccessToken"))
					if link != "" {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(link, link)).Do(); err != nil {
							log.Println("Line reply error=", err)
						}
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("統計圖表產出錯誤，請聯絡開發人員！")).Do(); err != nil {
							log.Println("Line reply error=", err)
						}
					}
				case GetUserMonthStatistic:
					now := time.Now()
					thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
					stats := GetMonthStatDataByUser(thisMonth, event.Source.UserID)

					var chartData []chart.Value

					for _, v := range stats.Data {
						chartData = append(chartData, chart.Value{Label: fmt.Sprintf("%s $%d", v.Class, v.Sum), Value: float64(v.Sum)})
					}

					chart := GetChart(chartData)
					link := UploadToImgur(chart, os.Getenv("ImgurAccessToken"))
					if link != "" {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(link, link)).Do(); err != nil {
							log.Println("Line reply error=", err)
						}
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("統計圖表產出錯誤，請聯絡開發人員！")).Do(); err != nil {
							log.Println("Line reply error=", err)
						}
					}
				}
			}
		}
	}
}
