package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

func (P *Plugin) RunPlugin() []Content {
	contents := []Content{}

	var wg sync.WaitGroup

	start := time.Now()
	spaces := P.getTotalSpaces()
	log.Info().Msgf("Time to get all spaces with concurrency - %d ms", int(time.Since(start).Milliseconds()))

	start = time.Now()
	for _, space := range spaces {
		limit := make(chan struct{}, 5)
		pages := P.getTotalPages(space).Pages
		for _, p := range pages {
			limit <- struct{}{}
			wg.Add(1)
			go func() {
				content := P.getContent(p, space)
				contents = append(contents, content)
				<-limit
				wg.Done()
			}()
		}
		wg.Wait()
	}
	log.Info().Msgf("Time to get all pages with concurrency - %d:%d:%d", int(time.Since(start).Hours()), int(time.Since(start).Minutes()), int(time.Since(start).Seconds()))

	log.Info().Msg("Confluence plugin completed successfully")

	return contents
}

func (P *Plugin) getTotalSpaces() []Space {
	totalSpaces := P.getSpaces(0)
	var count int32 = 1
	var mutex sync.Mutex
	var wg sync.WaitGroup

	if totalSpaces.Size == 25 {
		for threadCount := 0; threadCount < 4; threadCount++ {
			wg.Add(1)
			go P.ThreadGetSpaces(&count, &totalSpaces, &mutex, &wg)
		}
	}

	wg.Wait()
	log.Info().Msgf(" Total of %d Spaces detected", len(totalSpaces.Results))

	return totalSpaces.Results
}

func (P *Plugin) ThreadGetSpaces(count *int32, totalSpaces *Space_Response, mutex *sync.Mutex, wg *sync.WaitGroup) {
	var moreSpaces Space_Response
	for {
		atomic.AddInt32(count, 1)
		lastSpaces := P.getSpaces(int(*count-1) * 25)
		moreSpaces.Results = append(moreSpaces.Results, lastSpaces.Results...)

		if lastSpaces.Size == 0 {
			mutex.Lock()
			totalSpaces.Results = append(totalSpaces.Results, moreSpaces.Results...)
			mutex.Unlock()
			wg.Done()
			break
		}
	}
}

func (P *Plugin) getSpaces(start int) Space_Response {
	resp := HttpRequest("GET", fmt.Sprintf("%srest/api/space?start=%d", P.url, start), P.email, P.token)

	data_obj := Space_Response{}
	jsonErr := json.Unmarshal(resp, &data_obj)
	if jsonErr != nil {
		log.Info().Msgf("Unauthorized to get spaces. Already got %d spaces")
	}

	return data_obj
}

func (P *Plugin) getTotalPages(space Space) Page_Result {
	totalPages := P.getPages(space, 0)
	var count int32 = 1
	var mutex sync.Mutex
	var wg sync.WaitGroup

	if len(totalPages.Pages) == 25 {
		for threadCount := 0; threadCount < 20; threadCount++ {
			wg.Add(1)
			go P.ThreadGetPages(space, &count, &totalPages, &mutex, &wg)
		}
	}

	wg.Wait()
	log.Info().Msgf(" Space - %s have %d pages", space.Name, len(totalPages.Pages))

	return totalPages
}

func (P *Plugin) ThreadGetPages(space Space, count *int32, totalPages *Page_Result, mutex *sync.Mutex, wg *sync.WaitGroup) {
	var morePages Page_Result
	for {
		atomic.AddInt32(count, 1)
		lastPages := P.getPages(space, int(*count-1)*25)
		morePages.Pages = append(morePages.Pages, lastPages.Pages...)

		if len(lastPages.Pages) == 0 {
			mutex.Lock()
			totalPages.Pages = append(totalPages.Pages, morePages.Pages...)
			mutex.Unlock()
			wg.Done()
			break
		}
	}
}

func (P *Plugin) getPages(space Space, start int) Page_Result {
	resp := HttpRequest("GET", fmt.Sprintf("%srest/api/space/%s/content?start=%d", P.url, space.Key, start), P.email, P.token)

	pages_obj := Page_Response{}
	jsonErr := json.Unmarshal(resp, &pages_obj)
	if jsonErr != nil {
		log.Fatal().Msg("Unauthorized!")
	}

	return pages_obj.Results
}

func (P *Plugin) getContent(page Page, space Space) Content {
	source := P.url + "rest/api/content/" + page.ID + "?expand=body.storage,body.view.value,version,history.previousVersion"
	originalUrl := P.url + "spaces/" + space.Key + "/pages/" + page.ID
	return Content{Content: string(HttpRequest("GET", source, P.email, P.token)), Source: source, OriginalUrl: originalUrl}
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
type Space struct {
	ID    int               `json:"id"`
	Key   string            `json:"key"`
	Name  string            `json:"name"`
	Links map[string]string `json:"_links"`
}

type Space_Response struct {
	Results []Space `json:`
	Size    int     `json:size`
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
