package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	req, err := http.NewRequest("GET", "https://streamtape.com/e/w72PqjR9DDFJbb9/2082_J", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("referer", "https://momoiroadult.com/archives/93315")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Process the response (for example, print the status code)
	println("Response Status:", res.Status)

	f, err := os.Create("output.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		panic(err)
	}
}
