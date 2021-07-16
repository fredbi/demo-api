# demo-api

This repo demonstrates a simple REST API to upload and download images
from a persistent store.

Showcase: this simplistic webapp exposes an API that empowers a user to
upload images to a store, list all uploaded images, get, update or delete
a single image.

## Build and run

> **NOTE**: you need gcc installed

```
go get github.com/fredbi/demo-api/cmd/images
```

This should install the server binary in your path (e.g. `$GOPATH/bin`).

```
images serve
badger 2021/07/16 13:52:36 INFO: Set nextTxnTs to 0
2021/07/16 13:52:36 listening on 0.0.0.0:3000
```

You may now access the demo from your browser, e.g. ```http://localhost:3000```.

> **IMPORTANT**: POST, PATCH and DELETE operations are subject to very crude implementation on the client side.
> After such a call, the user is required to navigate back manually and get back to the list of images.
>
> I did not want to load more javascript than strictly necessary to test the app. Obvisouly, this should be
> done with some async call to refresh the page with the list.

## What is demonstrated?

I wanted to demonstrate the following things here:

* how to scaffold a simple set of REST endpoints in golang, e.g. CRUD
* how to use multipart forms as a flexible vehicle for data uploads (as proposed by the OpenAPI and GraphQL API standards)
* how to organize the app in modular parts: CLI, http handling, database repository
* how to use a KV store I wanted to experiment with: BadgerDB. This is my first hands-on experiment with this.
* how to carry out simple image processing in golang

Frontend is part of the demo, but I did not want to demonstrate anything there: the
frontend on demo is trivial, and essentially rendered by the server. Obviously,
this is not in the spirit of demonstrating an API, which would more likely render JSON messages.

The idea with this demo was _NOT_ to figure out a full-fledged UI (e.g. react or angular), hence
the rendering extra code on the server side. On a real implementation, we most likely would pare that part away.

## Repo structure

```
./assets            static HTML to run the demo (this is then baked in a generated variable)
./cmd/images        CLI to run the server
./app               server main runtime and app REST handlers
./pkg/repo          database access layer
./pkg/repo/images   database access implementation using BadgerDB as a KV store
./pkg/image-utils   utility to resize an image
```

## Blind spots

* [ ] Unit tests
    No unit tests are provided in this barebones demo.     
    The repo package defines an interface to allow for mocking and easier unit tests of the "app" package.

* [ ] Configurability
    No configuration is provided (flags, env vars, etc.). This is on purpose: I wanted a zero-config demo.
    The server runtime may easily be augmented to inject a config registry into the app (e.g. using github.com/pf13/viper).
  
* [ ] Interoperability
    On a real-life application, switch to JSON messages.

* [ ] Observability
    No observability. On a real-life micro-service deployment, we want logging, tracing, metrics and healthchecks.
    All of these would make a much beefier runtime, so this is skipped for brevity.

* [ ] Security
    A webapp or API operating under normal conditions should publish only TLS endpoints.
    Again the zero-config objective for the demo made me skip that part. 

    It is likely that this service would require some kind of authentication (e.g. cookie, API key, oauth2...).
    It is rather easy to inject middleware to resolve this for all routes.
    Authenticated users might then be restricted to their own private set of data items.

* [ ] Deployment
    To deploy the app, I would need to prepare a `Dockerfile` to capture the (very simple) build logic.
    In addition, I should prepare some kubernetes manifest for a deployment (or some AWS EC equivalent configuration).

* [ ] Scalability
    The http server for this simple API is naturally scalable and highly concurrent.
    However, things go differently with the underlying database (an embedded KV store): this db
    is essentially single-process and we cannot deploy multiple containers opening the same database.
    If required, this would advocate for a db server (e.g. postgres or similar) _or_  provide a RPC-like implementation
    for the repo package.

* [ ] Performances
    Performance is not what we were after for this demo. However, BadgerDB is very fast.
    A few comments on potential performance issues:
    * in the current implementation, we do not leverage streaming i/o's end-to-end: there is a buffering stage when
      interacting with the db. Interestingly, this is possible with postgres driver, but seems difficult to achieve
      with badgerDB
    * I favored simplicity over performance for the thumbsnail image resizing. Faster resizers exist out there however.
