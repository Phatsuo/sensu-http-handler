# sensu-http-handler
Proof of concept generic http handler

This handler should be considered experimental and entirely unsupported.

If you are interested in extending or fixing this, you are encourage to fork this repo.


## Usage
```
sensu-http-handler --help
Proof of concept generic http handler

Usage:
  sensu-http-handler [flags]
  sensu-http-handler [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of this plugin

Flags:
      --header stringToString   Add additional HTTP header in format key=value (ex: 'X-Sensu-Header=test value') can be used multiple times (default [])
  -h, --help                    help for sensu-http-handler
      --insecure-skip-verify    Skip TLS verifications for https urls
  -m, --method string           The http(s) method: POST and PATCH supported (default "POST")
  -u, --url string              The http(s) url
  -d, --data                    The data to pass, will be treated as a template
  -v, --verbose                 Verbose logging

```

## Important Arguments

### `--url`
Required. This is the full path to the http(s) endpoint you need to POST/PATCH

### `--method`
The HTTP method, currently supported POST and PATCH. Defaults to POST

### `--header`
This allows you to add headers using a `key=value` pattern. You can use this multiple times, each time calling a different key to set multiple headers.

### `--data`
The POST or PATCH body. Treated as a Sensu template and processed with the event object.


## Environment Variables
|Argument   |Environment Variable |
|-----------|---------------------|
|--url      |HTTP_HANDLER_URL     |
|--method   |HTTP_HANDLER_METHOD  |


## What it's suppose to do.
Ideally this handler should allow you to post json representation of Sensu events to a random http endpoint (like a webhook).

The json passed in is treated as a Sensu template and processed with the event object.

### Why not just use curl in a shell script as a handler command?

Good question.  This is exactly what I would normally do. 
But there are situations where using Sensu assets are preferred and curl is difficult to package as an asset because of its library dependencies. And a call to curl is not able to be treated as a Sensu template.

So here we are, This golang executable should be relatively easy to package as a Sensu asset, and should expose just enough http configuration to allow you to send data to a simple webhook url expecting json data.

It will not be as featureful as a curl script, though. So if you need advanced http features like proxy support or private cert, this proof-of-concept handler problably isn't going to get there out of the box.
