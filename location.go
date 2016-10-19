package main 
import (
	"appengine"
	"net/http"
	"net/url"	
	"appengine/urlfetch"
	"fmt"
)

type geocodingResults struct {
	Results []struct {
		Locations []Location `json:"locations"`
	} `json:"results"`
}


type Coordinates struct {
	Longitude float64
	Latitude  float64
}

type Location struct {
	Street      string `json:"street"`
	City        string `json:"adminArea5"`
	State       string `json:"adminArea3"`
	PostalCode  string `json:"postalCode"`
	County      string `json:"adminArea4"`
	CountryCode string `json:"adminArea1"`
	LatLng      GeoPoint `json:"latLng"`
	Type        string `json:"type"`
	DragPoint   bool   `json:"dragPoint"`
}



type GeoPoint struct {
	Lng float64	
	Lat  float64
}

// meaty function for attempting to get as many locations from a TwitterResponse as possible. It uses both coordinate data and string-based location data to create a list of coordinates.
func compileLocationResults(twitterResp *TwitterResponse, w http.ResponseWriter, r *http.Request, keyword string) []Coordinates {
	i := 0
	coordSlice := make([]Coordinates, 1)
	coordSlice = append(coordSlice, getStoredData(w, r, keyword)...)
	i = len(coordSlice)
	for _, v := range twitterResp.Statuses {
		if (((len(v.Geo.Coordinates) == 0) || len(v.Place.Bounds.Coordinates) == 0) && i < 10) {
				bufferCoords := geocodeSearch(w, r, v.User.Location)
				if (bufferCoords.Latitude != 0 && bufferCoords.Longitude != 0) {			
								storeTweet(v.Id, v.Text, bufferCoords, w, r)
								coordSlice = append(coordSlice, bufferCoords)
								i = i + 1
				}
		} else {
			if len(v.Geo.Coordinates) >= 1 {
				bufferCoords := Coordinates{
					Latitude:  v.Geo.Coordinates[0],
					Longitude: v.Geo.Coordinates[1],
				}
				if bufferCoords.Latitude != 0 && bufferCoords.Longitude != 0 {			
					storeTweet(v.Id, v.Text, bufferCoords, w, r)
					coordSlice = append(coordSlice, bufferCoords)
					i = i + 1
				}
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
				if bufferCoords.Latitude != 0 && bufferCoords.Longitude != 0 {			
					storeTweet(v.Id, v.Text, bufferCoords, w, r)
					coordSlice = append(coordSlice, bufferCoords)
					i = i + 1
				}
				if bufferCoords2.Latitude != 0 && bufferCoords2.Longitude != 0 {			
					storeTweet(v.Id, v.Text, bufferCoords2, w, r)
					coordSlice = append(coordSlice, bufferCoords2)
					i = i + 1
				}

				if bufferCoords3.Latitude != 0 && bufferCoords3.Longitude != 0 {			
					storeTweet(v.Id, v.Text, bufferCoords3, w, r)
					coordSlice = append(coordSlice, bufferCoords3)
					i = i + 1
				}

				if bufferCoords4.Latitude != 0 && bufferCoords4.Longitude != 0 {			
					storeTweet(v.Id, v.Text, bufferCoords4, w, r)
					coordSlice = append(coordSlice, bufferCoords4)
					i = i + 1
				}
			}
		}
	}
	
	fmt.Fprintf(w, `<script>document.title = "TweetMap (" + %d + " unique locations)"</script>`, i-1)
	//fmt.Fprintf(w, "%+v", coordSlice)
	return coordSlice
}
// Use mapquest's api to determine a set of coordinates from a string-based location. Has a rate limit of 15000 requests per month so we use it as little as possible.
func geocodeSearch(w http.ResponseWriter, r *http.Request, location string) Coordinates {
				ctx := appengine.NewContext(r)
				hc := urlfetch.Client(ctx)
				req, err := http.NewRequest("GET", "https://open.mapquestapi.com/geocoding/v1/address?inFormat=kvp&outFormat=json&location=" + url.QueryEscape(location) + "&key=" + mapquestapikey, nil)
				resp, err := hc.Do(req)

				if err != nil {
					return Coordinates{
						Latitude: 0,
						Longitude: 0,
					}
				} else {

					defer resp.Body.Close()

					var result geocodingResults
					err = decoder(resp).Decode(&result)

					if err != nil {
						return Coordinates{
						Latitude: 0,
						Longitude: 0,
					}
					} else {

					
						if len(result.Results[0].Locations) > 0 {
							bufferCoords := Coordinates{
								Latitude: result.Results[0].Locations[0].LatLng.Lat,
								Longitude: result.Results[0].Locations[0].LatLng.Lng,
							}
							return bufferCoords
						}
					}
				}
				return Coordinates{
						Latitude: 0,
						Longitude: 0,
					}
}