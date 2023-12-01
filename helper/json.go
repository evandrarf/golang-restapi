package helper

import (
	"encoding/json"
	"net/http"
	"restful_api/model/web"
)

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	PanicIfError(err)
}

func WriteToResponseBody(writer http.ResponseWriter, data web.WebResponse) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(data.Code)
	err := json.NewEncoder(writer).Encode(data)
	PanicIfError(err)
}