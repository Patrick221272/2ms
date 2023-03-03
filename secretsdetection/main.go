package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const API_KEY = "<your_api_key>"
const USER = "<your_user_email>"
const BASE_URL = "https://checkmarx.atlassian.net/wiki/"

func main() {

	for _, space := range getSpaces() {
		for _, page := range getPages(space) {
			findSecrets(getContent(page))
		}
	}
}

func findSecrets(content string) {
	RegExp := []string{"PRIVATE KEY", "ASIA[0-9A-Z]{16}", "AKIA[0-9A-Z]{16}", "[0-9]+-[0-9A-Za-z_]{32}.apps.googleusercontent.com"}
	for _, value := range RegExp {
		pattern := regexp.MustCompile(value)
		matches := pattern.FindAllString(content, -1)

		if len(matches) > 0 {
			fmt.Println(matches)
		}
	}

}

func getSpaces() []Space_Result {
	resp := HttpRequest("GET", BASE_URL+"rest/api/space")

	data_obj := Space_Response{}
	jsonErr := json.Unmarshal(resp, &data_obj)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return data_obj.Results
}

func getPages(space Space_Result) []Page {
	resp := HttpRequest("GET", BASE_URL+"rest/api/space/"+space.Key+"/content")

	pages_obj := Page_Response{}
	jsonErr := json.Unmarshal(resp, &pages_obj)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return pages_obj.Results.Pages
}

func getContent(page Page) string {
	return string(HttpRequest("GET", BASE_URL+"rest/api/content/"+page.ID+"?expand=body.storage"))
}

// Utils
func HttpRequest(httpmethod string, uri string) []byte {
	req, err := http.NewRequest(httpmethod, uri, nil)

	req.SetBasicAuth(USER, API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

// Responses Structs
type Space_Result struct {
	ID    int               `json:"id"`
	Key   string            `json:"key"`
	Name  string            `json:"name"`
	Links map[string]string `json:"_links"`
}

type Space_Response struct {
	Results []Space_Result `json:`
}

type Page struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type Page_Result struct {
	Pages []Page `json:"results"`
}

type Page_Response struct {
	Results Page_Result `json:"page"`
}
