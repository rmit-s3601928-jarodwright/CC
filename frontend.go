package main 
import (
	"fmt"
	"net/http"
	"strconv"
)
// Turn a table of coordinates into a lot of js script strings compatible with google maps
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

// our html for this application
func heatMapPage(w http.ResponseWriter, r *http.Request, coords []Coordinates) {
	fmt.Fprintf(w, `<head>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
<script>
var keywordNumber = 0;

if (sessionStorage.getItem("keywordNumber"))
{
	keywordNumber = sessionStorage.getItem("keywordNumber");
}

function addKeywordSearch()
{
	for (i = 0; i < 7; i++)
	{
		if (sessionStorage.getItem(i) == document.getElementById("keyword").value)
		{
			for (j = i; j < 7; j++)
			{
				if (sessionStorage.getItem(j) != null && sessionStorage.getItem(j+1) != null)
				{
					sessionStorage.setItem(j, sessionStorage.getItem(j+1));
				}
			}
			keywordNumber--;
			break;
		}
	}

	if (keywordNumber == 7)
	{
		for (i = 0; i < 6; i++)
		{
			sessionStorage.setItem(i, sessionStorage.getItem(i+1));
		}
		sessionStorage.setItem(6, document.getElementById("keyword").value);
		}
	else
	{
		var keywordValue = keywordNumber;
		sessionStorage.setItem(keywordValue, document.getElementById("keyword").value);
		keywordNumber++;
		sessionStorage.setItem("keywordNumber", keywordNumber);
	}
}

function addKeyword(keyword)
{
	for (i = 0; i < 7; i++)
	{
		if (sessionStorage.getItem(i) == keyword)
		{
			for (j = i; j < 7; j++)
			{
				if (sessionStorage.getItem(j) != null && sessionStorage.getItem(j+1) != null)
				{
					sessionStorage.setItem(j, sessionStorage.getItem(j+1));
				}
			}
			keywordNumber--;
			break;
		}
	}

	if (keywordNumber == 7)
	{
		for (i = 0; i < 6; i++)
		{
			sessionStorage.setItem(i, sessionStorage.getItem(i+1));
		}
		sessionStorage.setItem(6, keyword);
		}
	else
	{
		var keywordValue = keywordNumber;
		sessionStorage.setItem(keywordValue, keyword);
		keywordNumber++;
		sessionStorage.setItem("keywordNumber", keywordNumber);
	}
}

</script>
<meta charset="utf-8">

<title>Heatmaps</title>

<style>
html, body {
	overflow: hidden;
	height: 100vh;
	margin: 0;
	padding: 0;
}
#map {
	float:left;
	center: new google.maps.LatLng(0,0);
	height: 100vh;
	width: 100vh;
	zoom:1;
}
#inputs {
	float:left;
	width: calc(100vw - 100vh);
	height: 100vh;
	font-size:72px;
}
#search {
	height:calc(12.5vh);
	width:calc(100vw - 100vh);
	float:left;
}
#search input {
	height:calc(12.5vh);
	width:calc(100vw - 100vh);
	font-size:12vh;
	text-align: center;
}
.floating-panel {
	width:calc(100vw - 100vh);
	height:calc(12.5vh);
	float: left;
	background-color: #9BC8FB;
	text-align: center;
	font-family: 'Roboto','sans-serif';
}
button {
	font-size:12vh;
	font-family: 'Roboto','sans-serif';
	border:none;
}
button:hover {
	background-color:#C7E0F9;
}
</style>
</head>

<body>
	<div id="inputs">
		<div id="search">
			<form action="/submit" onsubmit="addKeywordSearch();" "method="POST">
				<input name="keyword" id="keyword" type="text" autocomplete="off"/>
			</form>
		</div>
		<script>
		var j = 0;
		for (i = 6; i >= 0; i--)
		{
			if (sessionStorage.getItem(i))
			{
				jQuery('<button/>', {
					id: i,
					class: 'floating-panel',
					text: sessionStorage.getItem(i),
				}).click(
					function()
					{
						window.location = "\submit?keyword=" + $(this).text();
						addKeyword($(this).text());
					}).appendTo('#inputs');
			}
			else
			{
				j++;
			}
		}
		while (j>0)
		{
			jQuery('<div/>', {
				class: 'floating-panel',
			}).appendTo('#inputs');
			j--;
		}
		</script>
	</div>
	<div id="outputs">
		<div id="map"></div>
	</div>
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
