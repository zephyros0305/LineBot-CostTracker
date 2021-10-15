// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
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

				if err != nil {
					log.Println("Quota err:", err)
				}

				switch operation, operationData := DetermineOperation(message.Text); operation {
				case Error:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("錯誤的輸入格式，請重新輸入！")).Do(); err != nil {
						log.Print(err)
					}
				case KeepRecord:
					record := ConvertToRecord(*operationData)
					record.UserID = event.Source.UserID
					if record.Save() {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("成功儲存！\n%s-%d-%s", record.Class, record.Cost, record.Memo))).Do(); err != nil {
							log.Print(err)
						}
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("儲存資料錯誤，請晚點稍後再試！")).Do(); err != nil {
							log.Print(err)
						}
					}
				case GetRecord:
					records := GetLastRecords(uint(operationData.Number))
					responseContainer := GetListRecordResponse(records)

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("Record list", responseContainer)).Do(); err != nil {
						log.Print(err)
					}
				case GetStatistic:

				}

				log.Println(message.Text)
			}
		}
	}
}
