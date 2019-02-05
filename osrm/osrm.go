package osrm

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Point struct {
	Longitude float32
	Latitude  float32
}

func httpGet(url string) (interface{}, error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))
	return body, err
}

func GetRouteJson(pointA Point, pointB Point) {

	url := fmt.Sprintf("http://35.193.174.32:5000/route/v1/driving/%f,%f;%f,%f?overview=full&alternatives=false&steps=false&geometries=geojson", pointA.Longitude, pointA.Latitude, pointB.Longitude, pointB.Latitude)
	fmt.Println(url)
	resp, _ := httpGet(url)

	fmt.Println(resp)

}
