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
	"log"
	"os"
	"sync"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		PushChartMonthlyCron()
		log.Fatalln("PushChartMonthlyCron stop")
		wg.Done()
	}()

	go func() {
		err := ServeLineBot()
		log.Fatalln("LineBotServer error! err=", err)
		wg.Done()
	}()

	wg.Wait()
	log.Fatalln("All goroutine stop!")
}
