package main

import (
	"net/http"
	"net/url"
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strings"
)

type TwitterResponse struct {
	Statuses []struct {
		Text string `json:"text"`
		Id string `json:"id_str"`
		Geo  struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geo"`
		Place struct {
			Id     string `json:"id"`
			Bounds struct {
				Coordinates [][][]float64 `json:"coordinates"`
			} `json:"bounding_box"`
		} `json:"place"`
		User struct {
			Name     string `json:"name"`
			Location string `json:"location"`
		} `json:"user"`
	} `json:"statuses"`
}


type TweetData struct {
	Content string
	Gp GeoPoint
}



func tweetStorageKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "TweetStorage", "default_tweetstorage", 0, nil)
}

func requestKeyword(keyword string, accesstoken string, w http.ResponseWriter, r *http.Request, tweetArray *TwitterResponse) []Coordinates {
	ctx := appengine.NewContext(r)
	hc := urlfetch.Client(ctx)
	req, err := http.NewRequest("GET", "https://api.twitter.com/1.1/search/tweets.json?q="+ url.QueryEscape(keyword)+"&result_type=mixed&count=100", nil)
	req.Header.Add("Authorization", "Bearer "+accesstoken)
	resp, err := hc.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "<p> Error: %s </p>", err)
	} else {
		//fmt.Fprintf(w, "<p> %s </p></br>", body)
		//var teststring = `{ "statuses":[{"location": "6", "name": "test1", "meme": "none"},{"location": null, "name": "test"}],"meta": "2"}`
		var twitterResp TwitterResponse
		err := json.Unmarshal(body, &twitterResp)
		//fmt.Fprintf(w, "<p><strong>%s</strong></p>", body)

		if err != nil {
			fmt.Fprintf(w, "<p> Error: %s </p>", err)
		} else {
			fmt.Fprintf(w, "<p>%+v</p>", twitterResp)
			//fmt.Fprintf(w, "<p>%d</p>", len(twitterResp.Statuses))
			//return
			return compileLocationResults(&twitterResp, w, r, keyword)
		}
		//fmt.Fprintf(w, "<p>%s</p>", body)
	}
	return nil
}



func getStoredData(w http.ResponseWriter, r *http.Request, keyword string) []Coordinates {
	c := appengine.NewContext(r)
	coordSlice := make([]Coordinates, 1)
	query := datastore.NewQuery("TweetStorage").Ancestor(tweetStorageKey(c)).Limit(2000)
	tweets := make([]TweetData, 1)
	if _, err := query.GetAll(c, &tweets); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
            return nil
	}
	for _,v := range tweets {
		if strings.Contains(v.Content, keyword) {
			bufferCoords := Coordinates{
				Latitude:  v.Gp.Lat,
				Longitude: v.Gp.Lng,
			}
			coordSlice = append(coordSlice, bufferCoords)
		}
	}
	return coordSlice
}

func storeTweet(tweetid string, content string, coords Coordinates, w http.ResponseWriter,  r *http.Request) {
	c := appengine.NewContext(r)
	newGP := GeoPoint{
		Lat: coords.Latitude,
		Lng: coords.Longitude,
	}
	data := TweetData{
		Content: content,
		Gp: newGP,
	}
	key := datastore.NewKey(c, "TweetStorage", tweetid, 0, nil)

	_, err := datastore.Put(c, key, &data)
	if err != nil {
		fmt.Fprintf(w, "<p>%s</p>", err)
	}
}


