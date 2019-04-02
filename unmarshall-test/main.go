package main

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"strconv"
	"time"
)

type Metric struct {
	Ts         string `json:"ts"`
	DeviceID   string `json:"device_id"`
	Value      string `json:"value"`
	VariableID string `json:"variable_id"`
}

func main() {
	var metricsBuild []Metric

	l := 1000000
	for i := 0; i < l; i++ {
		metricsBuild = append(metricsBuild, Metric{
			Ts: strconv.Itoa(i),
			DeviceID: "685d04cc-ce30-4e79-ab42-10130e225f38",
			Value: "9600.00",
			VariableID: "111bfb55-4c68-4c81-bcfb-1539bb8dea39",
		})
	}

	sr := buildSearchResult(metricsBuild)

	var hits []json.RawMessage

	for _, hit := range sr.Hits.Hits {
		hits = append(hits, *hit.Source)
	}




	start := time.Now()







	metrics, _ := searchResultToBytes(hits)






	t := time.Now()
	fmt.Println(t.Sub(start))



	fmt.Println(len(metrics))

/*	var metrics []Metric
	err := json.Unmarshal(bytes, &metrics)
	if err != nil {
		fmt.Print(err)
	} else {

	}*/
}

func buildSearchResult(values []Metric) *elastic.SearchResult {
	searchResultMock := &elastic.SearchResult{}
	var hits []*elastic.SearchHit
	for _, v := range values {
		byteArray, _ := json.Marshal(&v)
		rawMsg := (*json.RawMessage)(&byteArray)

		hits = append(hits, &elastic.SearchHit{
			Source: rawMsg,
		})
	}
	searchHits := &elastic.SearchHits{
		Hits:      hits,
		TotalHits: int64(0),
		MaxScore:  nil,
	}
	searchResultMock.Hits = searchHits
	return searchResultMock
}





func searchResultToBytes(hits []json.RawMessage) ([]Metric, error) {
	return convert(hits), nil
}

func convert(hits []json.RawMessage) []Metric {
	var metrics []Metric
	if hits != nil {
		resultBytes, err := json.Marshal(hits)
		if err != nil {
			return nil
		}
		json.Unmarshal(resultBytes, &metrics)
	}
	return metrics
}




/*func searchResultToBytes(hits []json.RawMessage) ([]Metric, error) {
	var metrics []Metric


	ch := make(chan []Metric)
	for i := 0; i < len(hits); i+= 100 {
		// Correct way
		go func(i int) {
			convert(hits[i:i+100], ch)
		}(i)

		// Error
		go func() {
			fmt.Println(i)
			fmt.Println(i + 100)
			convert(hits[0:100], ch)
		}()
	}

	var i int
	t := len(hits)/100
	for {
		chMetrics := <-ch
		metrics = append(metrics, chMetrics...)
		i++
		if t == i {
			close(ch)
			break
		}
	}

	return metrics, nil
}

func convert(hits []json.RawMessage, ch chan<- []Metric) {
	if hits != nil {
		resultBytes, err := json.Marshal(hits)
		if err != nil {
			return
		}
		var metrics []Metric
		json.Unmarshal(resultBytes, &metrics)
		ch <- metrics
	}
}*/
