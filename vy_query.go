package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
    "log"

)

// Train json object from vy.no
type Train struct {
	ResultSetID string `json:"resultSetId"`
	Itineraries []struct {
		DebugInfo          interface{} `json:"debugInfo"`
		ID                 string      `json:"id"`
		From               string      `json:"from"`
		To                 string      `json:"to"`
		DepartureScheduled string      `json:"departureScheduled"`
		DepartureRealTime  string      `json:"departureRealTime"`
        cancelled          bool        `json:"cancelled"`
		ArrivalScheduled   string      `json:"arrivalScheduled"`
		ArrivalRealTime    interface{} `json:"arrivalRealTime"`
	}
}

// VyRequest request POST body
type VyRequest struct {
	From                  string `json:"from"`
	To                    string `json:"to"`
	Time                  string `json:"time"`
	LimitResultsToSameDay bool   `json:"limitResultsToSameDay"`
	Language              string `json:"language"`
	Passengers            []struct {
		Type           string        `json:"type"`
		CustomerNumber string        `json:"customerNumber"`
		Discount       string        `json:"discount"`
		Extras         []interface{} `json:"extras"`
	} `json:"passengers"`
	PriceNecessity string `json:"priceNecessity"`
}

// IsVyLate return returns how many seconds train is late
func CallVy(TrainTravel TrainTravel) Train {
	// Return delay in seconds for train home
	requestbyte := []byte(`
		{"from": "Lilleby",
		"to": "Hell",
		"time":"2019-09-09T15:00",
		"limitResultsToSameDay":true,
		"language":"no",
		"passengers": [{"type":"ADULT","customerNumber":"null","discount":"NONE","extras":[]}],
		"priceNecessity":"REQUIRED"}`)
	var vyReq VyRequest
	err := json.Unmarshal(requestbyte, &vyReq)
	vyReq.Time = time.Now().Format("2006-01-02") + "T" + TrainTravel.depatureTime
    vyReq.From = TrainTravel.from
    vyReq.To = TrainTravel.to
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(vyReq)
	req, err := http.NewRequest("POST", "https://booking.cloud.nsb.no/api/itineraries/search", b)

	req.Header.Add("Content-Type", "Application/Json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
        log.Fatal("CallVy err: ",err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	vy := Train{}
	err = json.Unmarshal(body, &vy)
    if err != nil {
        log.Fatal("CallVy err: ", err)
    }
    return vy

}
func IsVyLate(vy Train) float64 {

	if vy.Itineraries[0].DepartureRealTime == "" {
		return 0.0
	}
	DepartureRealTime, _ := time.Parse("2006-01-02T15:04:05", vy.Itineraries[0].DepartureRealTime)
	DepartureScheduled, _ := time.Parse("2006-01-02T15:04:05", vy.Itineraries[0].DepartureScheduled)

	return DepartureRealTime.Sub(DepartureScheduled).Seconds()
}

func IsCancelled(vy Train) float64 {
	if vy.Itineraries[0].cancelled {
        return 1.0
    }
    return 0.0
}
