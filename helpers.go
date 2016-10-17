package main 
import (
	"net/http"
	"encoding/json"
)

func decoder(resp *http.Response) *json.Decoder {
	return json.NewDecoder(resp.Body)
}