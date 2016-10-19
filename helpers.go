package main 
import (
	"net/http"
	"encoding/json"
)

// Quickly create a json decoder for use in geocodeSearch
func decoder(resp *http.Response) *json.Decoder {
	return json.NewDecoder(resp.Body)
}