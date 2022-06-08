package main

import (
	"fmt"
	"log"

	"github.com/sensu/sensu-go/types"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Example string
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
			Path:      "example",
			Env:       "CHECK_EXAMPLE",
			Argument:  "example",
			Shorthand: "e",
			Default:   "",
			Usage:     "An example string configuration option",
			Value:     &plugin.Example,
		},
	}
)

func main() {

	check := sensu.NewHandler(&plugin.PluginConfig, options, checkArgs, executeCheck)
	check.Execute()
}

func checkArgs(event *types.Event) error {
	if len(plugin.Example) == 0 {
		return fmt.Errorf("--example or CHECK_EXAMPLE environment variable is required")
	}
	return nil
}

func executeCheck(event *types.Event) error {
	log.Println("executing check with --example", plugin.Example)
	return nil
}
