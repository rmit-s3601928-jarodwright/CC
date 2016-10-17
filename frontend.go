package main 
import (
	"fmt"
	"net/http"
	"strconv"
)

func getCoords(coords []Coordinates) string {
	if coords == nil {
		return ""
	}
	s := ""
	for _, v := range coords {
		if (v.Latitude == 0 && v.Longitude == 0) {
			// skip bad values
		} else {
			s = s + ("new google.maps.LatLng(" + strconv.FormatFloat(v.Latitude, 'f', -1, 64) + ", " + strconv.FormatFloat(v.Longitude, 'f', -1, 64) + "),\n") 
		}
	}
	return s

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
src="https://maps.googleapis.com/maps/api/js?key=` + googlemapapikey + `&libraries=visualization&callback=initMap">
</script>
</body>
`)
}
