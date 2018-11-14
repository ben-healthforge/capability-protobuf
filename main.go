package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"net/http"
	"os"
	"strings"
	"text/template"
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
			SearchParam []struct {
				Name string
				Type string
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

	tmpl, err := template.New("server.proto.tmpl").Funcs(template.FuncMap{
		"normalise": normalise,
	}).ParseFiles("/home/ben/opt/go/src/github.com/ben-healthforge/capability-protobuf/server.proto.tmpl")
	if err != nil {
		glog.Exitf("%v", err)
	}

	for _, rest := range statement.Rest {
		if rest.Mode == "server" {
			err = tmpl.Execute(os.Stdout, rest)
			if err != nil {
				glog.Exitf("%v", err)
			}
		}
	}
}

func normalise(name string) string {
	return strings.Replace(name, "-",  "_", -1)
}
