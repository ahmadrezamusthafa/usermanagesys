package response

import (
	"encoding/json"
	"net/http"
	"time"
	"tokopedia.se.training/Project1/usermanagesys/api/response/dto"
)

type ResponseModul struct {
	start     time.Time
	resWriter http.ResponseWriter
}

func New(res http.ResponseWriter) *ResponseModul {
	return &ResponseModul{
		start:     time.Now(),
		resWriter: res,
	}
}

func (m *ResponseModul) WriteSuccess(objPost interface{}) {
	m.resWriter.Header().Set("Content-Type", "application/json")
	m.resWriter.WriteHeader(http.StatusOK)
	var apiResult = dto.APIResultDto{Error: nil, Result: objPost, Success: true}
	var result, err = json.Marshal(apiResult)
	if err != nil {
		http.Error(m.resWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	m.resWriter.Write(result)
}

func (m *ResponseModul) WriteError(errMsg string) {
	m.resWriter.Header().Set("Content-Type", "application/json")
	m.resWriter.WriteHeader(http.StatusInternalServerError)
	var apiResult = dto.APIResultDto{Error: &errMsg, Result: nil, Success: false}
	var result, err = json.Marshal(apiResult)
	if err != nil {
		http.Error(m.resWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	m.resWriter.Write(result)
}

func (m *ResponseModul) WriteBadRequest(errMsg string) {
	m.resWriter.Header().Set("Content-Type", "application/json")
	m.resWriter.WriteHeader(http.StatusBadRequest)
	var apiResult = dto.APIResultDto{Error: &errMsg, Result: nil, Success: false}
	var result, err = json.Marshal(apiResult)
	if err != nil {
		http.Error(m.resWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	m.resWriter.Write(result)
}

func (m *ResponseModul) WriteUnauthorized(errMsg string) {
	m.resWriter.Header().Set("Content-Type", "application/json")
	m.resWriter.WriteHeader(http.StatusUnauthorized)
	var apiResult = dto.APIResultDto{Error: &errMsg, Result: nil, Success: false}
	var result, err = json.Marshal(apiResult)
	if err != nil {
		http.Error(m.resWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	m.resWriter.Write(result)
}
