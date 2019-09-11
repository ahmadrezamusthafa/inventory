package response

import (
	"encoding/json"
	"github.com/rezamusthafa/inventory/api/response/results"
	"net/http"
)

func WriteSuccess(objPost interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	var apiResult = results.APIResult{Error: nil, Result: objPost, Success: true}
	var result, err = json.Marshal(apiResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func WriteError(errMsg string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	var apiResult = results.APIResult{Error: &errMsg, Result: nil, Success: false}
	var result, err = json.Marshal(apiResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func WriteBadRequest(errMsg string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	var apiResult = results.APIResult{Error: &errMsg, Result: nil, Success: false}
	var result, err = json.Marshal(apiResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}
