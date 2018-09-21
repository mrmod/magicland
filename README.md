## MagicLand

I love serverless compute concepts like Google's cloud functions and AWS's Lamdas. After doing a small demo of the same years ago using Nginx, a Go handler, and Mesos, I want to make something I can vomit a simple behavioral interface to and have it serve my most recent code.

## The Idea

* Create a MagicLand Service with a GitURL and a public DNS
* Create yields a CNAME to a MagicLand service URL
* The Magicland service URL is a reference to a list of IP:Port combinations
* On HTTP GET to the Public DNS, the CNAME is hit
* The service URL is called, and if empty, the GitURL for the service is cloned to a staging path
* The repo must have a `./index.js` exporting a `handle` function of the signature `(http.ClientRequest, http.ClientResponse)`

### Cloning the repo
* A list of available machine IPs is consulted and one with a knapsack large enough selected
* Clone begins
* A list of ports unused on the machine IP is consulted and the first free is selected
* A docker Aspen container is stood up in 0.1 CPU slice and 128 MB memory
* The container is started with an entrypoint of `magicland.init`
* Magicland init has an overlay of `/App` with the cloned repo
* Magicland init does `yarn install` if `package.json` is present
* Magicland init then executes `node magicland.js`
* An ExpressJS instance is drawn with a single route with `index.js.handle` as the handler
* A service entry is created in NGinx for the Magicland DNS (the customer CNAME) listening on the service port
* The HTTP request is sent to the service entry

There's probably a detail I didn't write yet. It'll be there. This is for fun.

## Scheduling

Mesos, ECS, Kubernoodle-salad; all fine at scheduling services. This is a "for fun" "not real life" project, so the scheduling is going to be crude and filled with bad assumptions. 

### Constraints

* Must be free/nearly free (eg: EC2 T1-micro)
* Must work in the development environment with only Docker and an internet connection
* Must consider Memory saturated when 96 MB remain on a system (62 MB + WTF space)

## Testing

`go test -v`

It sounds like "goatest Vee" if you wing it at the right velocity.

## Why Go?

"Where are the channels and goroutines and concurrency?". Your butt. Next question.
