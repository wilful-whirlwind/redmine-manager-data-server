package controller

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type BaseController struct {
	base   BaseActionInterface
	logger *slog.Logger
}

func generateLogger(r *http.Request) *slog.Logger {
	logGroup := slog.Group("request", "method", r.Method, "url", r.URL)
	t := time.Now()
	fileName := t.Format("20060102")
	file, err := os.OpenFile("./log/application_"+fileName+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{AddSource: true})).With(logGroup)
}
func PreExecute(w http.ResponseWriter, r *http.Request, targetAction BaseActionInterface) *BaseController {
	baseAction := BaseController{
		base:   targetAction,
		logger: generateLogger(r),
	}
	return &baseAction
}

func Execute(w http.ResponseWriter, r *http.Request, targetAction BaseActionInterface) {
	responseBody := make(map[string]interface{})
	switch r.Method {
	case "GET":
		responseBody, _ = targetAction.Get(w, r)
		break
	case "POST":
		responseBody, _ = targetAction.Post(w, r)
		break
	}
	setDefaultHeader(w)
	convertBodyToJSON(responseBody, w)
}

func setDefaultHeader(w http.ResponseWriter) {
	uuidWithHyphen := uuid.New()
	w.Header().Set("X-redmine-manager", uuidWithHyphen.String())
	w.Header().Set("Content-Type", "application/json")
}

func convertBodyToJSON(responseBody map[string]interface{}, w http.ResponseWriter) {
	bytes, err := json.Marshal(responseBody)
	if err != nil {
		fmt.Println("JSON marshal error: ", err)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		fmt.Println("JSON write error: ", err)
		return
	} else {
		fmt.Println(string(bytes))
	}
}

func convertResponse(body any) map[string]interface{} {
	jsonStr, _ := json.Marshal(body)
	result := make(map[string]interface{})
	err := json.Unmarshal(jsonStr, &result)
	if err != nil {
		return nil
	}
	return result
}
