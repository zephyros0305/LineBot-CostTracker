package main

import (
	"bytes"
	"mime/multipart"

	"github.com/wcharczuk/go-chart/v2"
)

func GetChart(data []chart.Value) *bytes.Buffer {

	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: data,
	}

	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	part, _ := writer.CreateFormFile("image", "dont care about name")
	pie.Render(chart.PNG, part)
	writer.Close()

	return buf
}
