REST
====

This example demonstrates a HTTP REST web service with some fixture data.
Follow along the example and patterns.

Also check routes.json for the generated docs from passing the -routes flag

Boot the server:
----------------
$ go run main.go

Client requests:
----------------
$ curl http://localhost:3333/

root.

$ curl http://localhost:3333/articles

[{"id":"1","title":"Hi"},{"id":"2","title":"sup"}]

$ curl http://localhost:3333/articles/1

{"id":"1","title":"Hi"}

$ curl -X DELETE http://localhost:3333/articles/1

{"id":"1","title":"Hi"}

$ curl http://localhost:3333/articles/1

"Not Found"

$ curl -X POST -d '{"id":"will-be-omitted","title":"awesomeness"}' http://localhost:3333/articles

{"id":"97","title":"awesomeness"}

$ curl http://localhost:3333/articles/97

{"id":"97","title":"awesomeness"}

$ curl http://localhost:3333/articles

[{"id":"2","title":"sup"},{"id":"97","title":"awesomeness"}]
