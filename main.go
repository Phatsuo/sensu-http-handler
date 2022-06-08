package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Url                string
	Method             string
	InsecureSkipVerify bool
	Verbose            bool
	Headers            map[string]string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-http-handler",
			Short:    "Proof of concept generic http handler",
			Keyspace: "sensu.io/plugins/sensu-http-handler/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "url",
			Env:       "HTTP_HANDLER_URL",
			Argument:  "url",
			Shorthand: "u",
			Default:   "",
			Usage:     "The http(s) url",
			Value:     &plugin.Url,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "method",
			Env:       "HTTP_HANDLER_METHOD",
			Argument:  "method",
			Shorthand: "m",
			Default:   "POST",
			Allow:     []string{"POST", "PATCH"},
			Usage:     "The http(s) method: POST and PATCH supported",
			Value:     &plugin.Method,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "insecure-skip-verify",
			Env:       "HTTP_HANDLER_INSECURE_SKIP_VERIFY",
			Argument:  "insecure-skip-verify",
			Shorthand: "",
			Default:   false,
			Usage:     "Skip TLS verifications for https urls",
			Value:     &plugin.InsecureSkipVerify,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "verbose",
			Env:       "HTTP_HANDLER_VERBOSE",
			Argument:  "verbose",
			Shorthand: "",
			Default:   false,
			Usage:     "Verbose logging",
			Value:     &plugin.Verbose,
		},
		&sensu.MapPluginConfigOption[string]{
			Argument: "header",
			Default:  map[string]string{},
			Usage:    "Add additional HTTP header in format key=value (ex: 'X-Sensu-Header=test value') can be used multiple times",
			Value:    &plugin.Headers,
		},
	}
)

func main() {

	check := sensu.NewHandler(&plugin.PluginConfig, options, checkArgs, executeCheck)
	check.Execute()
}

func checkArgs(event *corev2.Event) error {
	if len(plugin.Url) == 0 {
		return fmt.Errorf("--url most be provided")
	}
	return nil
}

func executeCheck(event *corev2.Event) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: plugin.InsecureSkipVerify},
	}
	client := &http.Client{Transport: tr}

	//Encode the data
	postJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	postBody := bytes.NewReader(postJSON)
	request, err := http.NewRequest(plugin.Method, plugin.Url, postBody)
	//Make request
	request.Header.Set("Content-Type", "application/json")
	for k, v := range plugin.Headers {
		request.Header.Set(k, v)
	}
	if plugin.Verbose {
		log.Println("sensu-http-handler --url", plugin.Url)
		for k, v := range request.Header {
			log.Printf("sensu-http-handler request headers  %v : %v", k, v)
		}
	}
	_, err = client.Do(request)
	return err
}
