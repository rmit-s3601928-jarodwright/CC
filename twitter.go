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
	b64 "encoding/base64"
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
)

func init() {
    http.HandleFunc("/", root)
    //http.HandleFunc("/authorise", authorise)
}

type Tweet struct {
	Geo string
	Content    string
}

type Token struct {
	Tokentype string
	Access_token string
}

type Coordinates struct {
	longitude float64
	latitude float64
}

func root(w http.ResponseWriter, r *http.Request){

	fmt.Fprintf(w, "<html><title>TweetMap</title>")
	testkeyword := "Earthquake"
	access_token, success := checkForAccessToken()
	if success == false {
			access_token = authorise("L23T9SJUKk4zGrZf0lGjhXQZV", "i7mCRyxSMUc1uS8c4EGGcZWM47gDTDOxNwE6PvURTCQIQlhi5f", w, r)
	}
	requestKeyword(testkeyword, access_token, w, r)
	fmt.Fprintf(w, "<p>%s</p>", access_token)
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
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token",strings.NewReader(form.Encode()) )
	req.Header.Add("Authorization", "Basic " + encoded)
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
	
	return  accesstoken
}

func drawMap(keyword string, resultslimit int, coords ...Coordinates) {

}

func requestKeyword(keyword string, accesstoken string, w http.ResponseWriter,r *http.Request)  {
	ctx := appengine.NewContext(r)
	hc := urlfetch.Client(ctx)
	req, err := http.NewRequest("GET", "https://api.twitter.com/1.1/search/tweets.json?q=" + keyword, nil)
	req.Header.Add("Authorization", "Bearer " + accesstoken)
	resp, err := hc.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "<p> Error: %s </p>", err)
	} else {
		fmt.Fprintf(w, "<p>%s</p>", body)
	}
	
}	
func checkForAccessToken() (string, bool) {
	return "nil", false
}