package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type UploadResponse struct {
	Data    ImageData `json:"data"`
	Status  int       `json:"status"`
	Success bool      `json:"success"`
}

type ImageData struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Datetime      int    `json:"datetime"`
	Type          string `json:"type"`
	Animated      bool   `json:"animated"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	Size          int    `json:"size"`
	Views         int    `json:"views"`
	Bandwidth     int    `json:"bandwidth"`
	Favorite      bool   `json:"favorite"`
	Account_url   string `json:"account_url"`
	Account_id    int    `json:"account_id"`
	Is_ad         bool   `json:"is_ad"`
	In_most_viral bool   `json:"in_most_viral"`
	Ad_type       int    `json:"ad_type"`
	Ad_url        string `json:"ad_url"`
	In_gallery    bool   `json:"in_gallery"`
	Deletehash    string `json:"deletehash"`
	Name          string `json:"name"`
	Link          string `json:"link"`
}

func UploadToImgur(byteImage []byte, clientId string) string {
	var imgBuf = bytes.NewBuffer(byteImage)
	buf := new(bytes.Buffer)
	multiWriter := multipart.NewWriter(buf)

	now := time.Now()

	part, _ := multiWriter.CreateFormFile("image", fmt.Sprintf("chart%d", now.Unix()))
	io.Copy(part, imgBuf)
	multiWriter.Close()

	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/upload", buf)
	req.Header.Set("Content-Type", multiWriter.FormDataContentType())
	req.Header.Set("Authorization", "Client-ID "+clientId)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return ""
	}
	defer res.Body.Close()
	jsonResponse, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(jsonResponse))

	var response UploadResponse
	err = json.Unmarshal(jsonResponse, &response)
	if err == nil {
		return response.Data.Link
	} else {
		log.Println(err)
		return ""
	}
}
