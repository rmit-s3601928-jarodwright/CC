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
	"strconv"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/submit", submit)
}

type Tweet struct {
	Content  string `json:"text"`
	Location string `json:"location"`
}

type Token struct {
	Tokentype    string
	Access_token string
}

type TwitterResponse struct {
	Statuses []struct {
		Text string `json:"text"`
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

type Coordinates struct {
	Longitude float64
	Latitude  float64
}

func root(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<html><title>TweetMap</title>")
	testkeyword := "Fitzroy"
	tweetArray := new(TwitterResponse)
	access_token, success := checkForAccessToken()
	if success == false {
		access_token = authorise("L23T9SJUKk4zGrZf0lGjhXQZV", "i7mCRyxSMUc1uS8c4EGGcZWM47gDTDOxNwE6PvURTCQIQlhi5f", w, r)
	}
	requestKeyword(testkeyword, access_token, w, r, tweetArray)
	//fmt.Fprintf(w, "<p>%s</p>", access_token)
	fmt.Fprintf(w, "<!DOCTYPE html><html>")
	heatMapPage(w, r, requestKeyword(testkeyword, access_token, w, r, tweetArray))
	fmt.Fprintf(w, "</html>")
}

func submit(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<!DOCTYPE html><html><title>TweetMap</title>")
	keyword := url.QueryEscape(r.FormValue("keyword"))
	tweetArray := new(TwitterResponse)
	access_token, success := checkForAccessToken()
	if success == false {
		access_token = authorise("L23T9SJUKk4zGrZf0lGjhXQZV", "i7mCRyxSMUc1uS8c4EGGcZWM47gDTDOxNwE6PvURTCQIQlhi5f", w, r)
	}
	//requestKeyword(testkeyword, access_token, w, r, tweetArray)
	//fmt.Fprintf(w, "<p>%s</p>", access_token)
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

func getCoords(coords []Coordinates) string {
	if coords == nil {
		return ""
	}
	s := ""
	for _, v := range coords {
		s = s + ("new google.maps.LatLng(" + strconv.FormatFloat(v.Latitude, 'f', -1, 64) + ", " + strconv.FormatFloat(v.Longitude, 'f', -1, 64) + "),\n")
	}
	return s

}

func requestKeyword(keyword string, accesstoken string, w http.ResponseWriter, r *http.Request, tweetArray *TwitterResponse) []Coordinates {
	ctx := appengine.NewContext(r)
	hc := urlfetch.Client(ctx)
	req, err := http.NewRequest("GET", "https://api.twitter.com/1.1/search/tweets.json?q="+keyword+"&result_type=mixed&count=100", nil)
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
			return compileLocationResults(&twitterResp, w, keyword)
		}
		//fmt.Fprintf(w, "<p>%s</p>", body)
	}
	return nil
}

func compileLocationResults(twitterResp *TwitterResponse, w http.ResponseWriter, keyword string) []Coordinates {
	i := 0
	coordSlice := make([]Coordinates, 1)
	for _, v := range twitterResp.Statuses {
		if (len(v.Geo.Coordinates) == 0) || len(v.Place.Bounds.Coordinates) == 0 {

		} else {
			if len(v.Geo.Coordinates) >= 1 {
				bufferCoords := Coordinates{
					Latitude:  v.Geo.Coordinates[0],
					Longitude: v.Geo.Coordinates[1],
				}
				coordSlice = append(coordSlice, bufferCoords)
				i = i + 1
			} else if len(v.Place.Bounds.Coordinates) >= 1 {
				bufferCoords := Coordinates{
					Latitude:  v.Place.Bounds.Coordinates[0][0][0],
					Longitude: v.Place.Bounds.Coordinates[0][0][1],
				}
				bufferCoords2 := Coordinates{
					Latitude:  v.Place.Bounds.Coordinates[0][1][0],
					Longitude: v.Place.Bounds.Coordinates[0][1][1],
				}
				bufferCoords3 := Coordinates{
					Latitude:  v.Place.Bounds.Coordinates[0][2][0],
					Longitude: v.Place.Bounds.Coordinates[0][2][1],
				}
				bufferCoords4 := Coordinates{
					Latitude:  v.Place.Bounds.Coordinates[0][3][0],
					Longitude: v.Place.Bounds.Coordinates[0][3][1],
				}
				coordSlice = append(coordSlice, bufferCoords)

				coordSlice = append(coordSlice, bufferCoords2)

				coordSlice = append(coordSlice, bufferCoords3)

				coordSlice = append(coordSlice, bufferCoords4)
			}
		}
	}
	fmt.Fprintf(w, "%+v", coordSlice)
	return coordSlice
}

func checkForAccessToken() (string, bool) { /* TODO */
	return "nil", false
}

func heatMapPage(w http.ResponseWriter, r *http.Request, coords []Coordinates) {
	fmt.Fprintf(w, `<head>
<meta charset="utf-8">
<title>Heatmaps</title>
<style>
html, body {
height: 100%%;
margin: 0;
padding: 0;
}
#map {
center: new google.maps.LatLng(0,0);
height: 100vh;
width: 100vh;
zoom:1;
float: left;
}
#floating-panel {
width: 25%%;
top: 10px;
float: left;
background-color: #fff;
padding: 5px;
border: 1px solid #999;
text-align: center;
font-family: 'Roboto','sans-serif';
line-height: 30px;
}
</style>
</head>

<body>
<div id="floating-panel">
<span>Search tweets</span>
<form action="/submit" "method="POST">
<input name="keyword" type="text"/>
<input type="submit" value="submit"/>
</form>
<button onclick="toggleHeatmap()">Toggle Heatmap</button>
<button onclick="changeGradient()">Change gradient</button>
<button onclick="changeRadius()">Change radius</button>
<button onclick="changeOpacity()">Change opacity</button>
</div>
<div id="map"></div>
<div id="keyWord"></div>
<script>

var map, heatmap;

function initMap() {
map = new google.maps.Map(document.getElementById("map"), {
zoom: 2,
center: {lat: 0, lng: 0},
mapTypeId: 'roadmap'
});

heatmap = new google.maps.visualization.HeatmapLayer({
data: getPoints(),
map: map
});
heatmap.set('radius', heatmap.get('radius') ? null : 75)
}

function toggleHeatmap() {
heatmap.setMap(heatmap.getMap() ? null : map);
}

function changeGradient() {
var gradient = [
'rgba(0, 255, 255, 0)',
'rgba(0, 255, 255, 1)',
'rgba(0, 191, 255, 1)',
'rgba(0, 127, 255, 1)',
'rgba(0, 63, 255, 1)',
'rgba(0, 0, 255, 1)',
'rgba(0, 0, 223, 1)',
'rgba(0, 0, 191, 1)',
'rgba(0, 0, 159, 1)',
'rgba(0, 0, 127, 1)',
'rgba(63, 0, 91, 1)',
'rgba(127, 0, 63, 1)',
'rgba(191, 0, 31, 1)',
'rgba(255, 0, 0, 1)'
]
heatmap.set('gradient', heatmap.get('gradient') ? null : gradient);
}

function changeRadius() {
heatmap.set('radius', heatmap.get('radius') ? null : 20);
}

function changeOpacity() {
heatmap.set('opacity', heatmap.get('opacity') ? null : 0.2);
}

function getPoints() {
return [
`+getCoords(coords)+`];
}
</script>
<script async defer
src="https://maps.googleapis.com/maps/api/js?key=AIzaSyCA2IXesNAu2eVxW2epTko-QTDxi5HqJkY&libraries=visualization&callback=initMap">
</script>
</body>
`)
}
