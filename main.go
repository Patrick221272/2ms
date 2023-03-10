package main

import (
	"encoding/json"
	"fmt"
	"github.com/checkmarx/2ms/Reporting"
	"github.com/zricethezav/gitleaks/v8/cmd/generate/config/rules"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	"io"
	"net/http"
	"regexp"
)

const API_KEY = "<your_api_key>"
const USER = "<your_user_email>"
const BASE_URL = "https://checkmarx.atlassian.net/wiki/"

func main() {

	err, rules := loadRules()
	if err != nil {
		return
	}

	cfg := config.Config{
		Rules: rules,
	}

	detector := detect.NewDetector(cfg)

	rawText := string(`
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCKLwIHewTIhcpH3WLnxZ61xBAk2lnkdahFxjHYi+khrENzbGr8
EeJDZ1FMUDDYGeLtjlROLHT41ovicFbsmgIU0QQVFewIAwvKIw5hBtq0TtO9CsXe
BaNmzw8ZduXJ/clOpdOF7/1ro485a+v956ZAhB2ohbk6qRqGyg3kaxclOQIDAQAB
AoGAV7z5QN6vbtLkWTUMc7VazHas+Xla0mCSc5sgUyqi4CqMuWEBnQON8tZLHHVe
ThhBqixRA0HfE5DGSQSjbJ9s6fD+Sjt0Qj2yer70FuEiR0uGM4tOAE7WbX+Ny7PT
gmDiWOITe7v0yzIgZzbLgPhg5SlCmiy8Nv2Zf/v54yLVPLECQQDbwpsuu6beMDip
kRB/msCAEEAstdfSPY8L9QySYxskkJvtWpWBu5trnRatiGoLYWvnsBzcL4xWGrs8
Tpr4hTirAkEAoPiRDHrVbkKAgrmLW/TrSDiOG8uXSTuvz4iFgzCG6Cd8bp7mDKhJ
l98Upelf0Is5sEnLDqnFl62LZAyckeThqwJAOjZChQ6QFSsQ11nl1OdZNpMXbMB+
euJzkedHfT9jYTwtEaJ9F/BqKwdhinYoIPudabHs8yZlNim+jysDQfGIIQJAGqlx
JPcHeO7M6FohKgcEHX84koQDN98J/L7pFlSoU7WOl6f8BKavIdeSTPS9qQYWdQuT
9YbLMpdNGjI4kLWvZwJAJt8Qnbc2ZfS0ianwphoOdB0EwOMKNygjnYx7VoqR9/h1
4Xgur9w/aLZrLM3DSatR+kL+cVTyDTtgCt9Dc8k48Q==
-----END RSA PRIVATE KEY-----



rdme_hdhehehehehehekjashdlkjedrtyfjaslfjdtyurhgowuhfjkasndlkythbsdfgetfde93

`)

	fragment := detect.Fragment{
		Raw:      rawText,
		FilePath: "directory\\rawStringAsFile.txt",
	}

	findings := detector.Detect(fragment)

	fmt.Printf("total secrets found: %v\n", len(findings))
	report := Reporting.Report{}
	for _, value := range findings {
		secret := Reporting.Secret{Description: value.Description, StartLine: value.StartLine, StartColumn: value.StartColumn, EndLine: value.EndLine, EndColumn: value.EndColumn, Value: value.Secret}
		if len(report.Results) > 0 {
			_, ok := report.Results[value.File]
			if ok {
				report = Reporting.AddSecretToFile(report, value, secret)
			} else {
				report = Reporting.CreateNewResult(report, value, secret)
			}
		} else {
			report = Reporting.CreateNewResult(report, value, secret)
		}

	}
	Reporting.ShowReport(report)

}

func loadRules() (error, map[string]config.Rule) {
	var configRules []*config.Rule

	configRules = append(configRules, rules.AdafruitAPIKey())
	configRules = append(configRules, rules.AdobeClientID())
	configRules = append(configRules, rules.AdobeClientSecret())
	configRules = append(configRules, rules.AgeSecretKey())
	configRules = append(configRules, rules.Airtable())
	configRules = append(configRules, rules.AlgoliaApiKey())
	configRules = append(configRules, rules.AlibabaAccessKey())
	configRules = append(configRules, rules.AlibabaSecretKey())
	configRules = append(configRules, rules.AsanaClientID())
	configRules = append(configRules, rules.AsanaClientSecret())
	configRules = append(configRules, rules.Atlassian())
	configRules = append(configRules, rules.AWS())
	configRules = append(configRules, rules.BitBucketClientID())
	configRules = append(configRules, rules.BitBucketClientSecret())
	configRules = append(configRules, rules.BittrexAccessKey())
	configRules = append(configRules, rules.BittrexSecretKey())
	configRules = append(configRules, rules.Beamer())
	configRules = append(configRules, rules.CodecovAccessToken())
	configRules = append(configRules, rules.CoinbaseAccessToken())
	configRules = append(configRules, rules.Clojars())
	configRules = append(configRules, rules.ConfluentAccessToken())
	configRules = append(configRules, rules.ConfluentSecretKey())
	configRules = append(configRules, rules.Contentful())
	configRules = append(configRules, rules.Databricks())
	configRules = append(configRules, rules.DatadogtokenAccessToken())
	configRules = append(configRules, rules.DigitalOceanPAT())
	configRules = append(configRules, rules.DigitalOceanOAuthToken())
	configRules = append(configRules, rules.DigitalOceanRefreshToken())
	configRules = append(configRules, rules.DiscordAPIToken())
	configRules = append(configRules, rules.DiscordClientID())
	configRules = append(configRules, rules.DiscordClientSecret())
	configRules = append(configRules, rules.Doppler())
	configRules = append(configRules, rules.DropBoxAPISecret())
	configRules = append(configRules, rules.DropBoxLongLivedAPIToken())
	configRules = append(configRules, rules.DropBoxShortLivedAPIToken())
	configRules = append(configRules, rules.DroneciAccessToken())
	configRules = append(configRules, rules.Duffel())
	configRules = append(configRules, rules.Dynatrace())
	configRules = append(configRules, rules.EasyPost())
	configRules = append(configRules, rules.EasyPostTestAPI())
	configRules = append(configRules, rules.EtsyAccessToken())
	configRules = append(configRules, rules.Facebook())
	configRules = append(configRules, rules.FastlyAPIToken())
	configRules = append(configRules, rules.FinicityClientSecret())
	configRules = append(configRules, rules.FinicityAPIToken())
	configRules = append(configRules, rules.FlickrAccessToken())
	configRules = append(configRules, rules.FinnhubAccessToken())
	configRules = append(configRules, rules.FlutterwavePublicKey())
	configRules = append(configRules, rules.FlutterwaveSecretKey())
	configRules = append(configRules, rules.FlutterwaveEncKey())
	configRules = append(configRules, rules.FrameIO())
	configRules = append(configRules, rules.FreshbooksAccessToken())
	configRules = append(configRules, rules.GoCardless())
	configRules = append(configRules, rules.GCPAPIKey())
	configRules = append(configRules, rules.GitHubPat())
	configRules = append(configRules, rules.GitHubFineGrainedPat())
	configRules = append(configRules, rules.GitHubOauth())
	configRules = append(configRules, rules.GitHubApp())
	configRules = append(configRules, rules.GitHubRefresh())
	configRules = append(configRules, rules.GitlabPat())
	configRules = append(configRules, rules.GitlabPipelineTriggerToken())
	configRules = append(configRules, rules.GitlabRunnerRegistrationToken())
	configRules = append(configRules, rules.GitterAccessToken())
	configRules = append(configRules, rules.GrafanaApiKey())
	configRules = append(configRules, rules.GrafanaCloudApiToken())
	configRules = append(configRules, rules.GrafanaServiceAccountToken())
	configRules = append(configRules, rules.Hashicorp())
	configRules = append(configRules, rules.Heroku())
	configRules = append(configRules, rules.HubSpot())
	configRules = append(configRules, rules.Intercom())
	configRules = append(configRules, rules.JWT())
	configRules = append(configRules, rules.KrakenAccessToken())
	configRules = append(configRules, rules.KucoinAccessToken())
	configRules = append(configRules, rules.KucoinSecretKey())
	configRules = append(configRules, rules.LaunchDarklyAccessToken())
	configRules = append(configRules, rules.LinearAPIToken())
	configRules = append(configRules, rules.LinearClientSecret())
	configRules = append(configRules, rules.LinkedinClientID())
	configRules = append(configRules, rules.LinkedinClientSecret())
	configRules = append(configRules, rules.LobAPIToken())
	configRules = append(configRules, rules.LobPubAPIToken())
	configRules = append(configRules, rules.MailChimp())
	configRules = append(configRules, rules.MailGunPubAPIToken())
	configRules = append(configRules, rules.MailGunPrivateAPIToken())
	configRules = append(configRules, rules.MailGunSigningKey())
	configRules = append(configRules, rules.MapBox())
	configRules = append(configRules, rules.MattermostAccessToken())
	configRules = append(configRules, rules.MessageBirdAPIToken())
	configRules = append(configRules, rules.MessageBirdClientID())
	configRules = append(configRules, rules.NetlifyAccessToken())
	configRules = append(configRules, rules.NewRelicUserID())
	configRules = append(configRules, rules.NewRelicUserKey())
	configRules = append(configRules, rules.NewRelicBrowserAPIKey())
	configRules = append(configRules, rules.NPM())
	configRules = append(configRules, rules.NytimesAccessToken())
	configRules = append(configRules, rules.OktaAccessToken())
	configRules = append(configRules, rules.PlaidAccessID())
	configRules = append(configRules, rules.PlaidSecretKey())
	configRules = append(configRules, rules.PlaidAccessToken())
	configRules = append(configRules, rules.PlanetScalePassword())
	configRules = append(configRules, rules.PlanetScaleAPIToken())
	configRules = append(configRules, rules.PlanetScaleOAuthToken())
	configRules = append(configRules, rules.PostManAPI())
	configRules = append(configRules, rules.Prefect())
	configRules = append(configRules, rules.PrivateKey())
	configRules = append(configRules, rules.PulumiAPIToken())
	configRules = append(configRules, rules.PyPiUploadToken())
	configRules = append(configRules, rules.RapidAPIAccessToken())
	configRules = append(configRules, rules.ReadMe())
	configRules = append(configRules, rules.RubyGemsAPIToken())
	configRules = append(configRules, rules.SendbirdAccessID())
	configRules = append(configRules, rules.SendbirdAccessToken())
	configRules = append(configRules, rules.SendGridAPIToken())
	configRules = append(configRules, rules.SendInBlueAPIToken())
	configRules = append(configRules, rules.SentryAccessToken())
	configRules = append(configRules, rules.ShippoAPIToken())
	configRules = append(configRules, rules.ShopifyAccessToken())
	configRules = append(configRules, rules.ShopifyCustomAccessToken())
	configRules = append(configRules, rules.ShopifyPrivateAppAccessToken())
	configRules = append(configRules, rules.ShopifySharedSecret())
	configRules = append(configRules, rules.SidekiqSecret())
	configRules = append(configRules, rules.SidekiqSensitiveUrl())
	configRules = append(configRules, rules.SlackAccessToken())
	configRules = append(configRules, rules.SlackWebHook())
	configRules = append(configRules, rules.StripeAccessToken())
	configRules = append(configRules, rules.SquareAccessToken())
	configRules = append(configRules, rules.SquareSpaceAccessToken())
	configRules = append(configRules, rules.SumoLogicAccessID())
	configRules = append(configRules, rules.SumoLogicAccessToken())
	configRules = append(configRules, rules.TeamsWebhook())
	configRules = append(configRules, rules.TelegramBotToken())
	configRules = append(configRules, rules.TravisCIAccessToken())
	configRules = append(configRules, rules.Twilio())
	configRules = append(configRules, rules.TwitchAPIToken())
	configRules = append(configRules, rules.TwitterAPIKey())
	configRules = append(configRules, rules.TwitterAPISecret())
	configRules = append(configRules, rules.TwitterAccessToken())
	configRules = append(configRules, rules.TwitterAccessSecret())
	configRules = append(configRules, rules.TwitterBearerToken())
	configRules = append(configRules, rules.Typeform())
	configRules = append(configRules, rules.VaultBatchToken())
	configRules = append(configRules, rules.VaultServiceToken())
	configRules = append(configRules, rules.YandexAPIKey())
	configRules = append(configRules, rules.YandexAWSAccessToken())
	configRules = append(configRules, rules.YandexAccessToken())
	configRules = append(configRules, rules.ZendeskSecretKey())
	configRules = append(configRules, rules.GenericCredential())

	// ensure rules have unique ids
	rules := make(map[string]config.Rule)
	for _, rule := range configRules {
		// required to be empty when not running via cli. otherwise rule will be ignored
		rule.Keywords = []string{}
		rules[rule.RuleID] = *rule
	}

	return nil, rules
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
	resp := HttpRequest("GET", BASE_URL+"rest/api/space/"+space.Key+"/content?limit=90")

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
