package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"os"
)

type CapabilityStatement struct {
	ResourceType string
	Rest []struct {
		Mode string
		Resource []struct {
			Type string
			Documentation string
			Interaction []struct {
				Code string
				Documentation string
			}
		}
	}
}

func main() {
	server := os.Args[1]

	client := &http.Client{}

	req, err := http.NewRequest("GET", server + "/metadata", nil)
	if err != nil {
		glog.Exitf("%v", err)
	}
	req.Header.Add("Accept", "application/fhir+json")

	resp, err := client.Do(req)
	if err != nil {
		glog.Exitf("%v", err)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	var statement CapabilityStatement
	err = decoder.Decode(&statement)
	if err != nil {
		glog.Exitf("%v", err)
	}

	fmt.Printf("%#v\n", statement)
}
