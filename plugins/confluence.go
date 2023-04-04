package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"strings"
)

const argConfluence = "confluence"
const argConfluenceSpaces = "confluence-spaces"
const argConfluenceUsername = "confluence-username"
const argConfluenceToken = "confluence-token"
const defaultConfluenceWindow = 25

type ConfluencePlugin struct {
	Plugin
	URL      string
	Token    string
	Username string
	Spaces   []string
}

func (p *ConfluencePlugin) IsEnabled() bool {
	return p.Enabled
}

func (p *ConfluencePlugin) DefineCommandLineArgs(cmd *cobra.Command) error {
	flags := cmd.Flags()
	flags.StringP(argConfluence, "", "", "scan confluence url")
	flags.StringP(argConfluenceSpaces, "", "", "confluence spaces")
	flags.StringP(argConfluenceUsername, "", "", "confluence username or email")
	flags.StringP(argConfluenceToken, "", "", "confluence token")
	return nil
}

func (p *ConfluencePlugin) Initialize(cmd *cobra.Command) error {
	flags := cmd.Flags()
	confluenceUrl, _ := flags.GetString(argConfluence)
	if confluenceUrl == "" {
		return errors.New("confluence URL arg is missing. Plugin initialization failed")
	}

	confluenceUrl = strings.TrimRight(confluenceUrl, "/")

	confluenceSpaces, _ := flags.GetString(argConfluenceSpaces)
	confluenceUsername, _ := flags.GetString(argConfluenceUsername)
	confluenceToken, _ := flags.GetString(argConfluenceToken)

	if confluenceUsername == "" || confluenceToken == "" {
		log.Warn().Msg("confluence credentials were not provided. The scan will be made anonymously only for the public pages")
	}

	p.Token = confluenceToken
	p.Username = confluenceUsername
	p.URL = confluenceUrl
	p.Spaces = strings.Split(confluenceSpaces, ",")
	p.Enabled = true
	return nil
}

func (p *ConfluencePlugin) GetItems(itemsChan chan Item) {
	spacesChan := make(chan []ConfluenceSpaceResult, 0)
	pagesChan := make(chan ConfluencePageResponse, 0)

	go p.getTotalSpaces(spacesChan)

	for spaces := range spacesChan {
		for _, space := range spaces {
			go p.getTotalPages(space, pagesChan)
		}
	}

	/*timer := time.NewTimer(time.Second * 30)

	go func() {
		<-timer.C
		if len(pagesChan) == cap(pagesChan) {
			log.Debug().Msg("Closing pages chan")
			close(pagesChan)
		} else {
			timer.Reset(time.Second * 30)
		}

	}()*/

	for pages := range pagesChan {
		for _, page := range pages.Results.Pages {
			pageContent, err := p.getContent(page, pages.Space)
			if err != nil {
				//return nil, err
			}

			itemsChan <- *pageContent

			/*if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(time.Second * 30)*/
		}
	}

	log.Debug().Msg("Closing items chan")
	//close(itemsChan)
	log.Debug().Msg("Confluence plugin completed successfully")

}

func (p *ConfluencePlugin) getTotalSpaces(spacesChan chan []ConfluenceSpaceResult) {
	totalSpaces, _ := p.getSpaces(0)
	spacesChan <- totalSpaces.Results

	actualSize := totalSpaces.Size
	lastSize := totalSpaces.Size

	if totalSpaces.Size == defaultConfluenceWindow {
		for lastSize != 0 {
			moreSpaces, _ := p.getSpaces(actualSize)
			spacesChan <- moreSpaces.Results
			actualSize += moreSpaces.Size
			lastSize = moreSpaces.Size
		}
	}
	log.Debug().Msg("Closing spaces chan")
	close(spacesChan)
}

func (p *ConfluencePlugin) getSpaces(start int) (*ConfluenceSpaceResponse, error) {
	url := fmt.Sprintf("%s/rest/api/space?start=%d", p.URL, start)
	body, err := p.httpRequest(http.MethodGet, url)
	if err != nil {
		return nil, fmt.Errorf("unexpected error creating an http request %w", err)
	}

	response := &ConfluenceSpaceResponse{}
	jsonErr := json.Unmarshal(body, response)
	if jsonErr != nil {
		return nil, fmt.Errorf("could not unmarshal response %w", err)
	}

	return response, nil
}

func (p *ConfluencePlugin) getTotalPages(space ConfluenceSpaceResult, pagesChan chan ConfluencePageResponse) {
	totalPages, err := p.getPages(space, 0)

	pagesChan <- *totalPages

	if err != nil {
		//return nil, fmt.Errorf("unexpected error creating an http request %w", err)
	}

	actualSize := len(totalPages.Results.Pages)
	lastSize := len(totalPages.Results.Pages)

	if actualSize == defaultConfluenceWindow {
		for lastSize == defaultConfluenceWindow {
			morePages, err := p.getPages(space, actualSize)

			if err != nil {
				//return nil, fmt.Errorf("unexpected error creating an http request %w", err)
			}

			pagesChan <- *morePages
			actualSize += len(morePages.Results.Pages)
			lastSize = len(morePages.Results.Pages)
		}
	}

	log.Info().Msgf(" Space - %s have %d pages", space.Name, len(totalPages.Results.Pages))
}

func (p *ConfluencePlugin) getPages(space ConfluenceSpaceResult, start int) (*ConfluencePageResponse, error) {
	url := fmt.Sprintf("%s/rest/api/space/%s/content?start=%d", p.URL, space.Key, start)
	body, err := p.httpRequest(http.MethodGet, url)

	if err != nil {
		return nil, fmt.Errorf("unexpected error creating an http request %w", err)
	}

	response := ConfluencePageResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		return nil, fmt.Errorf("could not unmarshal response %w", err)
	}

	response.Space = space.Key
	return &response, nil
	//return &response.Results, nil
}

func (p *ConfluencePlugin) getContent(page ConfluencePage, spaceKey string) (*Item, error) {
	url := p.URL + "/rest/api/content/" + page.ID + "?expand=body.storage,body.view.value,version,history.previousVersion"
	originalUrl := p.URL + "/spaces/" + spaceKey + "/pages/" + page.ID
	request, err := p.httpRequest(http.MethodGet, url)

	if err != nil {
		return nil, fmt.Errorf("unexpected error creating an http request %w", err)
	}

	content := &Item{
		Content: string(request),
		Source:  url,
		ID:      originalUrl,
	}
	return content, nil
}

func (p *ConfluencePlugin) httpRequest(method string, url string) ([]byte, error) {
	var err error

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unexpected error creating an http request %w", err)
	}

	if p.Username != "" && p.Token != "" {
		request.SetBasicAuth(p.Username, p.Token)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to send http request %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("unexpected error creating an http request %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("error calling http url \"%v\". status code: %v", url, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unexpected error reading http response body %w", err)
	}

	return body, nil
}

type ConfluenceSpaceResult struct {
	ID    int               `json:"id"`
	Key   string            `json:"key"`
	Name  string            `json:"Name"`
	Links map[string]string `json:"_links"`
}

type ConfluenceSpaceResponse struct {
	Results []ConfluenceSpaceResult `json:"results"`
	Size    int                     `json:"size"`
}

type ConfluencePage struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Space string ""
}

type ConfluencePageResult struct {
	Pages []ConfluencePage `json:"results"`
}

type ConfluencePageResponse struct {
	Results ConfluencePageResult `json:"page"`
	Space   string
}
