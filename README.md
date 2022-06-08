# sensu-http-handler
Proof of concept generic http handler

This handler should be considered experimental and entirely unsupported.

If you are interested in extending or fixing this, you are encourage to fork this repo.

### What it's suppose to do.
Ideally this handler should allow you to post json representation of Sensu events to a random http endpoint (like a webhook)

If you need to modify the json event representation in any way, you'll need to do that as part of a Sensu mutator called prior to using the handler.

### Why not just use curl in a shell script as a handler command?

Good question.  This is exactly what I would normally do. 
But there are situations where using Sensu assets are preferred and curl is difficult to package as an asset because of its library dependencies.

So here we are, This golang executable should be relatively easy to package as a Sensu asset, and should expose just enough http configuration to allow you to send data to a simple webhook url expecting json data.

It will not be as featureful as a curl script, though. 

