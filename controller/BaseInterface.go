package controller

import "net/http"

type BaseActionInterface interface {
	Execute(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error)
	Post(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error)
	Patch(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error)
}
