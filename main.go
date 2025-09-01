package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	yuril := "https://yxtwitter.com"

	resp, err := http.Get(yuril)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	imgSrc := ""
	doc.Find("div.pb-28 img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			imgSrc = src
			return
		}
	})
	if imgSrc == "" {
		fmt.Println("No image found in div.pb-28")
		return
	}

	imgResp, err := http.Get(imgSrc)
	if err != nil {
		panic(err)
	}
	defer imgResp.Body.Close()

	tmpFile, err := os.CreateTemp("", "bokita-*.jpg")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, imgResp.Body)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("viu", tmpFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running viu:", err)
	}
}
