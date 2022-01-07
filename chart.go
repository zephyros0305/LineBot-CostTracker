package main

import (
	"bytes"
	"log"

	"github.com/wcharczuk/go-chart/v2"
)

func GetChart(data []chart.Value) []byte {

	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: data,
	}

	log.Println("Chart data", data)

	var buf = new(bytes.Buffer)
	pie.Render(chart.PNG, buf)

	log.Println("Chart bytes", buf.Bytes())

	return buf.Bytes()
}
