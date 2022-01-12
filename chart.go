package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/v2"
)

func GetChart(data []chart.Value) []byte {

	var dataWithFont []chart.Value
	ZHFont := getZHFont()

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

	log.Println("Chart bytes", buf.Bytes())

	return buf.Bytes()
}

func getZHFont() *truetype.Font {
	fontFile := "NotoSansTC-Regular.otf"

	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Println(err)
		return nil
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return nil
	}
	return font
}
