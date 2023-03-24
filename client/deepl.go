package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func ProvideAPIClient() ApiClient {
	return ApiClient{}
}

type ApiClient struct{}

func (c ApiClient) Request(inputText string) (string, error) {
	key := os.Getenv("API_KEY")
	sourceLang := "EN"
	targetLang := "JA"

	apiUrl := "https://api-free.deepl.com/v2/translate"
	values := url.Values{}
	values.Set("text", inputText)
	values.Set("source_lang", sourceLang)
	values.Set("target_lang", targetLang)

	req, _ := http.NewRequest("POST", apiUrl, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", key))

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(resp.Body)

	var result struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Translations[0].Text, nil
}
