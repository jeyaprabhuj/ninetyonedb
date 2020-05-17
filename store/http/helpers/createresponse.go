package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/jeyaprabhuj/ninetyonedb/store/http/response"
)

func CreateErrorResponse(code int, msg string) []byte {
	var bError []byte
	kvError := &response.ErrorResponse{Code: 404, Message: msg}
	bError, err := json.Marshal(kvError)
	if err != nil {
		fmt.Println(err)
	}
	return bError
}

func CreateJSONResponse(w http.ResponseWriter, reply []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(reply)
}
