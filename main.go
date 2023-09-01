package main

import (
	"fmt"
	"github.com/gorilla/schema"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"net/url"
	"time"
)

/*
curl 'http://localhost:9091/model/allocation/compute?window=6d&aggregate=controller&step=1d&accumulate=false'

# https://github.com/opencost/opencost/blob/d4b47c5a75bb50ce1b3a2225c1fbc75968478ae3/pkg/costmodel/aggregation.go#L2227-L2297

window
resolution  duration
step duration
aggregate []string

includeIdle bool
accumulate bool

accumulateBy

idleByNode bool
includeProportionalAssetResourceCosts bool
includeAggregatedMetadata bool




*/

type AccumulateOption string

const (
	AccumulateOptionNone    AccumulateOption = ""
	AccumulateOptionAll     AccumulateOption = "all"
	AccumulateOptionHour    AccumulateOption = "hour"
	AccumulateOptionDay     AccumulateOption = "day"
	AccumulateOptionWeek    AccumulateOption = "week"
	AccumulateOptionMonth   AccumulateOption = "month"
	AccumulateOptionQuarter AccumulateOption = "quarter"
)

type Request struct {
	Window                                string           `json:"window,omitempty"`
	Resolution                            *metav1.Duration `json:"resolution,omitempty"`
	Step                                  *metav1.Duration `json:"step,omitempty"`
	Aggregate                             []string         `json:"aggregate,omitempty"`
	IncludeIdle                           bool             `json:"includeIdle,omitempty"`
	Accumulate                            bool             `json:"accumulate,omitempty"`
	AccumulateBy                          AccumulateOption `json:"accumulateBy,omitempty"`
	IdleByNode                            bool             `json:"idleByNode,omitempty"`
	IncludeProportionalAssetResourceCosts bool             `json:"includeProportionalAssetResourceCosts,omitempty"`
	IncludeAggregatedMetadata             bool             `json:"includeAggregatedMetadata,omitempty"`
}

/*
window=6d
aggregate=controller
step=1d
accumulate=false'
*/
func main() {
	r := Request{
		Window:                                "6d",
		Resolution:                            nil,
		Step:                                  &metav1.Duration{Duration: 6 * 24 * time.Hour},
		Aggregate:                             []string{"controller"},
		IncludeIdle:                           false,
		Accumulate:                            false,
		AccumulateBy:                          "",
		IdleByNode:                            false,
		IncludeProportionalAssetResourceCosts: false,
		IncludeAggregatedMetadata:             false,
	}
	var encoder = schema.NewEncoder()
	form := url.Values{}
	err := encoder.Encode(r, form)
	if err != nil {
		panic(err)
	}

	u, err := url.Parse("http://localhost:9091/model/allocation/compute")
	if err != nil {
		panic(err)
	}
	u.RawQuery = form.Encode()
	fmt.Println(u.String())

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
}
