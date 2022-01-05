package main

import (
	"bytes"

	"github.com/wcharczuk/go-chart/v2"
)

func GetChart(data []chart.Value) []byte {

	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: data,
	}

	var buf = new(bytes.Buffer)
	pie.Render(chart.PNG, buf)

	return buf.Bytes()
}
