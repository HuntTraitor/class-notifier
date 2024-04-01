package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	formUrl := "https://pisa.ucsc.edu/class_search/index.php"

	formData := url.Values{
		"binds[:term]":       {"2244"},
		"binds[:reg_status]": {"all"},
	}

	// encodedFormData := formData.Encode()

	client := &http.Client{}

	req, err := http.NewRequest("POST", formUrl, strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Golang_Super_Bot/0.1")
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

}
