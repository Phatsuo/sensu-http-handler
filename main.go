package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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
	buf         bytes.Buffer
	requestBody io.Reader
	plugin      = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-http-handler",
			Short:    "Proof of concept generic http handler",
			Keyspace: "sensu.io/plugins/sensu-http-handler/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Env:       "HTTP_HANDLER_URL",
			Argument:  "url",
			Shorthand: "u",
			Default:   "",
			Usage:     "The http(s) url",
			Value:     &plugin.Url,
		},
		&sensu.PluginConfigOption[string]{
			Env:       "HTTP_HANDLER_METHOD",
			Argument:  "method",
			Shorthand: "m",
			Default:   "POST",
			Allow:     []string{"POST", "PATCH"},
			Usage:     "The http(s) method: POST and PATCH supported",
			Value:     &plugin.Method,
		},
		&sensu.PluginConfigOption[bool]{
			Env:       "HTTP_HANDLER_INSECURE_SKIP_VERIFY",
			Argument:  "insecure-skip-verify",
			Shorthand: "",
			Default:   false,
			Usage:     "Skip TLS verifications for https urls",
			Value:     &plugin.InsecureSkipVerify,
		},
		&sensu.PluginConfigOption[bool]{
			Env:       "HTTP_HANDLER_VERBOSE",
			Argument:  "verbose",
			Shorthand: "v",
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

	handler := sensu.NewHandler(&plugin.PluginConfig, options, checkArgs, sendRequest)
	//This handler is expected to be used with mutated events, and thus the json passed via stdin will not be a valid event
	//Disable event reading and handle reading stdin elsewhere.
	handler.DisableReadEvent()

	//Lets check to see if stdin has content, if it does copy stdin into a new reader we can use
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error checking stdin: %v\n", err)
		panic(err)
	}
	//Check the Mode bitmask for Named Pipe to indicate stdin is connected
	if fi.Mode()&os.ModeNamedPipe != 0 {
		requestBody = io.TeeReader(os.Stdin, &buf)
	}
	// execute the handler business logic: sendRequest
	handler.Execute()
}

func checkArgs(event *corev2.Event) error {
	if len(plugin.Url) == 0 {
		return fmt.Errorf("--url most be provided")
	}
	return nil
}

func sendRequest(event *corev2.Event) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: plugin.InsecureSkipVerify},
	}
	client := &http.Client{Transport: tr}

	//prep the request
	request, err := http.NewRequest(plugin.Method, plugin.Url, requestBody)
	if err != nil {
		return err
	}
	//Make request
	request.Header.Set("Content-Type", "application/json")
	for k, v := range plugin.Headers {
		request.Header.Set(k, v)
	}
	if plugin.Verbose {
		log.Println("sensu-http-handler request url:", plugin.Url)
		for k, v := range request.Header {
			log.Printf("sensu-http-handler request header:  %v :: %v", k, v)
		}
		var buf bytes.Buffer
		var requestBodyBytes []byte
		if requestBody != nil {
			requestBodyCopy := io.TeeReader(requestBody, &buf)
			requestBodyBytes, _ = ioutil.ReadAll(requestBodyCopy)
		}
		log.Println("sensu-http-handler request body:", strings.TrimSpace(string(requestBodyBytes)))
	}
	_, err = client.Do(request)
	return err
}
