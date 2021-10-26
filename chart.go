package main

import (
	"bytes"
	"mime/multipart"

	"github.com/wcharczuk/go-chart/v2"
)

func GetChart() *bytes.Buffer {

	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: []chart.Value{
			{Value: 5, Label: "Blue"},
			{Value: 5, Label: "Green"},
			{Value: 4, Label: "Gray"},
			{Value: 4, Label: "Orange"},
			{Value: 3, Label: "Brown"},
			{Value: 3, Label: "??"},
			{Value: 2, Label: "Deep Blue"},
			{Value: 1, Label: "!!"},
		},
	}

	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	part, _ := writer.CreateFormFile("image", "dont care about name")
	pie.Render(chart.PNG, part)
	writer.Close()

	return buf
}
