package controllers

import (
	"encoding/json"
	"fmt"
	"go-app/app/models"
	"go-app/config"
	"html/template"
	"io/ioutil"
	"net/http"
)

func StartWebServer() error {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello/", helloHandler)
	http.HandleFunc("/bmonster", bmonsterHandler)
	http.HandleFunc("/api/bmonster", apiBmonsterHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.WebPort), nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// ResponseWriterにレスポンスの内容を書き込む
	fmt.Fprint(w, "<p>Index Page</p>")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// URLの情報を抽出
	path := r.URL.Path         // /hello/xxx
	v := path[len("/hello/"):] // xxx（スライスをうまく使って抽出）
	// ResponseWriterにレスポンスの内容を書き込む
	fmt.Fprintf(w, "<p>Hello! %v!</p>", v)
}

var templateFiles = []string{
	"app/views/bmonster.html",
}
var templates = template.Must(template.ParseFiles(templateFiles...))

func bmonsterHandler(w http.ResponseWriter, r *http.Request) {
	//t, _ := template.ParseFiles("app/views/bmonster.html")
	//t.Execute(w, nil)
	err := templates.ExecuteTemplate(w, "bmonster.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func APIError(w http.ResponseWriter, errorMessage string, code int) {
	// 200以外のステータスコードをセット
	w.WriteHeader(code)
	// レスポンスヘッダーをセット
	w.Header().Set("Content-Type", "application/json")
	// JSONレスポンス生成
	jsonError, _ := json.Marshal(JSONError{Error: errorMessage, Code: code})
	// JSON返却
	w.Write(jsonError)
}

func apiBmonsterHandler(w http.ResponseWriter, r *http.Request) {
	var performer string

	// クエリ解析
	if r.Method == "GET" {
		performer = r.URL.Query().Get("performer")
	} else if r.Method == "POST" && r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		performer = r.PostFormValue("performer")
	} else if r.Method == "POST" && r.Header.Get("Content-Type") == "application/json" {
		body, _ := ioutil.ReadAll(r.Body)
		var p models.Program
		_ = json.Unmarshal(body, &p)
		performer = p.Performer
	} else {
		APIError(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// データ取得
	var df *models.DataFrameProgram
	var err error
	if performer != "" {
		df, err = models.GetProgramByPerformer(performer)
	} else {
		df, err = models.GetAllProgram()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// データがなければエラー返却
	if df.Programs == nil {
		APIError(w, "No programs found.", http.StatusBadRequest)
		return
	}

	// JSON生成
	js, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// JSON返却
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
