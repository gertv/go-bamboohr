// Copyright 2014 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package bamboo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Bamboo struct {
	subdomain, key string
	debug          bool
	client         *http.Client
	base           string
}

type Calendar struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Type     string   `xml:"type,attr"`
	Start    string   `xml:"start"`
	End      string   `xml:"end"`
	Employee Employee `xml:"employee"`
}

type Employee struct {
	Id   int    `xml:"id,attr"`
	Name string `xml:",chardata"`
}

// Start using the API here -
func BambooHR(subdomain, key string) Bamboo {
	return Bamboo{subdomain, key, false, &http.Client{}, "https://api.bamboohr.com/api/gateway.php"}
}

// Enable/disable debug mode. When debug mode is enabled,
// you will get additional logging showing the HTTP requests
// and responses
func (b *Bamboo) Debug(d bool) {
	b.debug = d
}

// Configure a custom HTTP client (e.g. to configure a proxy server)
func (b *Bamboo) Client(client *http.Client) {
	b.client = client
}

// Get a calendar that shows who's out
func (b Bamboo) WhosOut(from, to string) (cal Calendar, err error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/v1/time_off/whos_out?start=%s&end=%s", b.base, b.subdomain, from, to), nil)
	log.Printf("%v", req)
	if err != nil {
		return
	}

	req.SetBasicAuth(b.key, "x")

	resp, err := b.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := body(resp)
	if err != nil {
		return
	}
	if b.debug {
		log.Printf("Got response %s: %s", resp.Status, data)
	}

	err = xml.Unmarshal(data, &cal)

	return
}

// Extract body from the HTTP response
func body(resp *http.Response) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (i Item) StartTime() time.Time {
	time, err := time.Parse("2006-01-02", i.Start)
	if err != nil {
		panic(err)
	}
	return time
}

func (i Item) EndTime() time.Time {
	time, err := time.Parse("2006-01-02", i.End)
	if err != nil {
		panic(err)
	}
	return time
}
