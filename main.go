package main

import (
	"fmt"
	//"io"
	"io/ioutil"
	//"log"
	//"html/template"
	"net/http"
	"net/url"
	//"time"
	"strings"
	//"oauth"
	"appengine"
	"appengine/urlfetch"
	b64 "encoding/base64"
	"encoding/json"
)

const (
	consumerkey = "L23T9SJUKk4zGrZf0lGjhXQZV" // replace with your consumer key from twitter.com
	consumersecretkey = "i7mCRyxSMUc1uS8c4EGGcZWM47gDTDOxNwE6PvURTCQIQlhi5f" // replace with your consumer secret key from twitter.com
	googlemapapikey = "AIzaSyCA2IXesNAu2eVxW2epTko-QTDxi5HqJkY" // replace with your api key from google maps api
	mapquestapikey  = "AbuyAhixGfSilbEtGF10ot8ZVQeC24KQ" // replace witth your api key
)

type Token struct {
	Tokentype    string
	Access_token string
}

func getConsumerKey() string {
	return consumerkey
}

func getConsumerSecret() string {
	return consumersecretkey
}

func getGoogleApiKey() string {
	return googlemapapikey
}
func getMapQuestApiKey() string {
	return mapquestapikey
}

func init() {
	
	http.HandleFunc("/", root)
	http.HandleFunc("/submit", submit)
}


func root(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<html><title>TweetMap</title>")
	fmt.Fprintf(w, "<!DOCTYPE html><html>")
	heatMapPage(w, r, nil)
	fmt.Fprintf(w, "</html>")
}

func submit(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<!DOCTYPE html><html><title>TweetMap</title>")
	keyword := url.QueryEscape(r.FormValue("keyword"))
	tweetArray := new(TwitterResponse)
	access_token, success := checkForAccessToken()
	if success == false {
		access_token = authorise(consumerkey, consumersecretkey, w, r)
	}
	heatMapPage(w, r, requestKeyword(keyword, access_token, w, r, tweetArray))
	fmt.Fprintf(w, "</html>")
}


func authorise(consumerkey string, consumersecretkey string, w http.ResponseWriter, r *http.Request) string {
	accesstoken := "nil"

	keys := consumerkey + ":" + consumersecretkey
	encoded := b64.StdEncoding.EncodeToString([]byte(keys))

	ctx := appengine.NewContext(r)
	hc := urlfetch.Client(ctx)
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader(form.Encode()))
	req.Header.Add("Authorization", "Basic "+encoded)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	resp, err := hc.Do(req)
	if err != nil {
		fmt.Fprintf(w, "<p> Error: %s </p>", err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(w, "<p> Error: %s </p>", err)
		} else {
			//fmt.Fprintf(w, "<p> %s </p><p> %s </p>", resp.Status, body)
			var t Token
			err = json.Unmarshal(body, &t)
			if err != nil {
				fmt.Fprintf(w, "<p> Error: %s </p>", err)
			} else {
				accesstoken = t.Access_token
			}
		}

	}

	return accesstoken

}

func checkForAccessToken() (string, bool) { /* TODO */
	return "nil", false
}