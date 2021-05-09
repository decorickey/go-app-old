package bmonster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-app/app/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type APIClient struct {
	apiKey     string
	httpClient *http.Client
}

func New(apiKey string) *APIClient {
	return &APIClient{apiKey, &http.Client{}}
}

func (apiClient APIClient) createHeader() map[string]string {
	headers := make(map[string]string)
	if apiClient.apiKey != "" {
		headers["x-api-key"] = apiClient.apiKey
	}
	return headers
}

func (apiClient *APIClient) doRequest(rawurl, path, method string, params map[string]string, reqBody []byte) ([]byte, error) {
	// エンドポイントURL生成
	baseURL, err := url.Parse(rawurl)
	reference, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	endpoint := baseURL.ResolveReference(reference).String()

	// リクエスト生成
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// ヘッダー生成
	for key, value := range apiClient.createHeader() {
		req.Header.Add(key, value)
	}

	// GET用クエリ生成
	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	// リクエスト実行
	res, err := apiClient.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// レスポンス読み込み
	resBody, err := ioutil.ReadAll(res.Body)
	return resBody, err
}

func (apiClient APIClient) ScrapingProgramList(ch chan []map[string]string, studio models.Studio) {
	// HTMLデータ取得
	var baseURL = "https://www.b-monster.jp"
	params := map[string]string{"studio_name": studio.Name, "studio_code": studio.Code}
	body, err := apiClient.doRequest(baseURL, "/reserve/", "GET", params, nil)
	if err != nil {
		log.Println(err)
		ch <- make([]map[string]string, 0)
	}

	// HTML解析
	programList, err := analyzeHTML(body, studio)
	if err != nil {
		log.Println(err)
		ch <- make([]map[string]string, 0)
	 }
	 ch <- programList
}

func (apiClient APIClient) GetLatestProgramList() ([]models.Program, error) {
	var baseURL = "https://om14r3diye.execute-api.ap-northeast-1.amazonaws.com"
	body, err := apiClient.doRequest(baseURL, "/v1/bmonster", "GET", nil, nil)
	if err != nil {
		return nil, err
	}

	var programList []models.Program
	err = json.Unmarshal(body, &programList)
	if err != nil {
		return nil, err
	}
	return programList, nil
}

func analyzeHTML(body []byte, studio models.Studio) (programList []map[string]string, err error) {
	// バイトストリームで読み込むため一時的にファイル書き込み（Lambdaでは/tmp配下でないと権限がない）
	outputFile := "/tmp/" + studio.Name + ".out"
	if err = ioutil.WriteFile(outputFile, body, 0666); err != nil {
		return nil, err
	}

	// HTMLデータ読み込み
	file, err := os.Open(outputFile)
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	// 解析
	doc.Find("#scroll-box .flex-no-wrap").Each(func(i int, day *goquery.Selection) {
		// 今日を基準に計算
		date := time.Now().AddDate(0, 0, i).Format("2006/1/2") // 1月2日午後3時4分5秒2006年
		// デバッグ確認用の日付
		debugDate := day.Find(".column-header .smooth-text").Text()

		day.Find(".panel").Each(func(j int, program *goquery.Selection) {
			fromTo := program.Find(".tt-time").Text()
			instructor := program.Find(".tt-instructor").Text()
			vol := program.Find(".tt-mode").Text()
			vol = strings.Replace(vol, "\n", "", -1) // -1ですべて置換
			vol = strings.Replace(vol, " ", "", -1)  // -1ですべて置換

			if len(fromTo) > 0 && len(instructor) > 0 && len(vol) > 0 {
				// 時間情報を抽出
				start := fromTo[:5]
				end := fromTo[8:]
				start = fmt.Sprintf("%s %s", date, start)
				end = fmt.Sprintf("%s %s", date, end)
				// time.ParseはデフォルトでUTCなのでParseInLocationでJSTを指定
				jst, _ := time.LoadLocation("Asia/Tokyo")
				startTime, _ := time.ParseInLocation("2006/1/2 15:04", start, jst) // 1月2日午後3時4分5秒2006年
				endTime, _ := time.ParseInLocation("2006/1/2 15:04", end, jst)     // 1月2日午後3時4分5秒2006年

				// レスポンス用mapを生成
				p := make(map[string]string)
				p["studio_name"] = studio.Name
				p["start_time"] = startTime.Format(time.RFC3339)
				p["end_time"] = endTime.Format(time.RFC3339)
				p["performer"] = instructor
				p["vol"] = vol
				p["debugDate"] = debugDate
				programList = append(programList, p)
			}
		})
	})
	return programList, nil
}
