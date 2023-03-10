package plugins

import (
	"encoding/json"
	"io"
	"net/http"
)

func (P *Plugin) RunPlugin() []string {
	contents := []string{}
	for _, space := range P.getSpaces() {
		for _, page := range P.getPages(space) {
			contents = append(contents, P.getContent(page))
		}
	}

	return contents
}

func (P *Plugin) getSpaces() []Space_Result {
	resp := HttpRequest("GET", P.url+"rest/api/space", P.email, P.token)

	data_obj := Space_Response{}
	jsonErr := json.Unmarshal(resp, &data_obj)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return data_obj.Results
}

func (P *Plugin) getPages(space Space_Result) []Page {
	resp := HttpRequest("GET", P.url+"rest/api/space/"+space.Key+"/content?limit=90", P.email, P.token)

	pages_obj := Page_Response{}
	jsonErr := json.Unmarshal(resp, &pages_obj)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return pages_obj.Results.Pages
}

func (P *Plugin) getContent(page Page) string {
	return string(HttpRequest("GET", P.url+"rest/api/content/"+page.ID+"?expand=body.storage", P.email, P.token))
}

// Utils
func HttpRequest(httpmethod string, uri string, email string, token string) []byte {
	var resp *http.Response
	var err error
	if email != "" && token != "" {
		resp, err = NonAuthenticatedHttpRequest(httpmethod, uri)
	} else {
		resp, err = AuthenticatedHttpRequest(httpmethod, uri, email, token)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func NonAuthenticatedHttpRequest(httpmethod string, uri string) (*http.Response, error) {
	var resp *http.Response
	var err error
	if httpmethod == "GET" {
		resp, err = http.Get(uri)
		if err != nil {
			panic(err)
		}
	}

	return resp, err
}

func AuthenticatedHttpRequest(httpmethod string, uri string, email string, key string) (*http.Response, error) {
	req, err := http.NewRequest(httpmethod, uri, nil)

	req.SetBasicAuth(email, key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp, err
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
