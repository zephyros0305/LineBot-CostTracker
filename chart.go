package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"sync"

	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/v2"
)

var (
	_defaultZHFontLock sync.Mutex
	_defaultZHFont     *truetype.Font
)

func GetChart(data []chart.Value) []byte {

	var dataWithFont []chart.Value
	ZHFont, err := getZHFont()

	if err != nil {
		log.Panicln("ZHFont load error=", err)
	}

	for _, d := range data {
		d.Style = chart.Style{Font: ZHFont}
		dataWithFont = append(dataWithFont, d)
	}

	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: dataWithFont,
	}

	log.Println("Chart data", data)

	var buf = new(bytes.Buffer)
	pie.Render(chart.PNG, buf)

	return buf.Bytes()
}

func getZHFont() (*truetype.Font, error) {

	if _defaultZHFont == nil {
		_defaultZHFontLock.Lock()
		defer _defaultZHFontLock.Unlock()
		if _defaultZHFont == nil {
			fontFile := ".fonts/TaipeiSansTCBeta-Regular.ttf"

			fontBytes, err := ioutil.ReadFile(fontFile)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			font, err := truetype.Parse(fontBytes)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			_defaultZHFont = font
		}
	}

	return _defaultZHFont, nil
}
