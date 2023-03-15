package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
)

func (P *Plugin) RunPlugin() []Content {
	contents := []Content{}

	contentChan := make(chan Content)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, space := range P.getTotalSpaces() {
		for _, page := range P.getTotalPages(space).Pages {
			wg.Add(1)
			go P.getContent(page, space, contentChan)
			//contents = append(contents, P.getContent(page, space))
		}
	}

	for range P.getTotalSpaces() {
		content := <-contentChan

		mutex.Lock()
		contents = append(contents, content)
		mutex.Unlock()
		wg.Done()
	}
	wg.Wait()

	log.Info().Msgf("Confluence plugin completed successfully. Total of %d items detected", len(contents))

	return contents
}

func (P *Plugin) getTotalSpaces() []Space_Result {
	totalSpaces := P.getSpaces(0)
	actualSize := totalSpaces.Size

	for actualSize != 0 {
		moreSpaces := P.getSpaces(totalSpaces.Size)
		totalSpaces.Results = append(totalSpaces.Results, moreSpaces.Results...)
		totalSpaces.Size += moreSpaces.Size
		actualSize = moreSpaces.Size
	}

	log.Info().Msgf(" Total of %d Spaces detected", len(totalSpaces.Results))

	return totalSpaces.Results
}

func (P *Plugin) getSpaces(start int) Space_Response {
	resp := HttpRequest("GET", fmt.Sprintf("%srest/api/space?start=%d", P.url, start), P.email, P.token)

	data_obj := Space_Response{}
	jsonErr := json.Unmarshal(resp, &data_obj)
	if jsonErr != nil {
		log.Fatal().Msg("Unauthorized!")
	}

	return data_obj
}

func (P *Plugin) getTotalPages(space Space_Result) Page_Result {
	totalPages := P.getPages(space, 0)
	actualSize := len(totalPages.Pages)

	for actualSize != 0 {
		morePages := P.getPages(space, len(totalPages.Pages))
		totalPages.Pages = append(totalPages.Pages, morePages.Pages...)
		actualSize = len(morePages.Pages)
	}

	log.Info().Msgf(" Space - %s have %d pages", space.Name, len(totalPages.Pages))

	return totalPages
}

func (P *Plugin) getPages(space Space_Result, start int) Page_Result {
	resp := HttpRequest("GET", fmt.Sprintf("%srest/api/space/%s/content?start=%d", P.url, space.Key, start), P.email, P.token)

	pages_obj := Page_Response{}
	jsonErr := json.Unmarshal(resp, &pages_obj)
	if jsonErr != nil {
		log.Fatal().Msg("Unauthorized!")
	}

	return pages_obj.Results
}

func (P *Plugin) getContent(page Page, space Space_Result, contentChan chan Content) Content {
	source := P.url + "rest/api/content/" + page.ID + "?expand=body.storage,body.view.value,version,history.previousVersion"
	originalUrl := P.url + "spaces/" + space.Key + "/pages/" + page.ID
	content := Content{Content: string(HttpRequest("GET", source, P.email, P.token)), Source: source, OriginalUrl: originalUrl}

	contentChan <- content
	return content
}

// Utils
func HttpRequest(httpmethod string, uri string, email string, token string) []byte {
	var resp *http.Response
	var err error

	if email == "" && token == "" {
		resp, err = NonAuthenticatedHttpRequest(httpmethod, uri)
	} else {
		resp, err = AuthenticatedHttpRequest(httpmethod, uri, email, token)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Msg("Unauthorized!")
	}

	return body
}

func NonAuthenticatedHttpRequest(httpmethod string, uri string) (*http.Response, error) {
	var resp *http.Response
	var err error
	if httpmethod == "GET" {
		resp, err = http.Get(uri)
		if err != nil {
			log.Fatal().Msg("Unauthorized!")
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
		log.Fatal().Msg("Unauthorized!")
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
	Size    int            `json:size`
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
