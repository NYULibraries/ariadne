# Ariadne

## Frontend

### Usage

Clone the repo and then:

```
cd frontend && yarn install
```

Start the client locally:

```
cd frontend && yarn start
```

Run in a container:

```
docker-compose up frontend
```

## Backend

The Ariadne backend is a REST API written in Go that takes
an [OpenURL](https://biblio.ugent.be/publication/760060/file/760063.pdf) submitted
via query string and returns JSON data containing electronic links from an
SFX Knowledgebase that represent NYU's e-holdings of the resource identified by
the OpenURL. It is essentially an API wrapper for the [SFX LinkResolver](https://exlibrisgroup.com/products/primo-discovery-service/sfx-link-resolver/),
and so is itself an OpenURL Link Resolver.

Under the hood it translates the OpenURL query string parameters to an XML `ContextObject`
for POST'ing to the [SFX Web Service](https://developers.exlibrisgroup.com/sfx/apis/web_services/openurl/),
parses the XML response into a JSON string which it delivers to the frontend.

### Usage

```
cd backend/
go run .
```

Run in a container:

```
docker-compose up backend
# Visit localhost:8080 with an OpenURL: 
# http://localhost:8080/v0/?sid=FirstSearch%3AWorldCat&genre=book&title=Fairy+tales&date=1898&aulast=Andersen&aufirst=H&auinitm=C&rfr_id=info%3Asid%2Ffirstsearch.oclc.org%3AWorldCat&rft.genre=book&rft_id=info%3Aoclcnum%2F7675437&rft.aulast=Andersen&rft.aufirst=H&rft.auinitm=C&rft.btitle=Fairy+tales&rft.date=1898&rft.place=Philadelphia&rft.pub=H.+Altemus+Co.&rft.genre=book
```

To run a [delve](https://github.com/go-delve/delve) debuggable containerized
instance:

```
docker-compose up -d backend-debug
```

To remotely debug this container instance in a
command-line [delve](https://github.com/go-delve/delve)
session, run `dlv`, making sure to map the remote container source path `/app/`
to
the local source path on the machine where the remote debugging session is being
run.
This will allow the `dlv` client to find the source file being referenced by the
`dlv` running in the container.

Example:

```
$ dlv connect localhost:2345
Type 'help' for list of commands.
(dlv) list main.main
Showing /app/main.go:11 (PC: 0x87c6ef)
Command failed: open /app/main.go: no such file or directory
(dlv) config substitute-path /app/ [LOCAL SOURCE CODE ABSOLUTE PATH]
(dlv) list main.main
Showing /app/main.go:11 (PC: 0x87c6ef)
     6:	)
     7:	
     8:	// Run on port 8080
     9:	const appPort = "8080"
    10:	
    11:	func main() {
    12:		router := NewRouter()
    13:	
    14:		log.Println("Listening on port", appPort)
    15:		log.Fatal(http.ListenAndServe(":"+appPort, router))
    16:	}
(dlv) 
```

To remote debug in Goland:
[Attach to a process in the Docker container](https://www.jetbrains.com/help/go/attach-to-running-go-processes-with-debugger.html#attach-to-a-process-in-the-docker-container)

To prevent the container from exiting after the headless `dlv` instance running
inside the container is killed, add this to the `backend-debug` service
definition
in _backend/Dockerfile.debug_:

```yaml
command: tail -f /dev/null
```

See [Keep a container running in compose \#1926](https://github.com/docker/compose/issues/1926)
for other methods.

### Testing

### frontend

Run all tests:

```
cd frontend/
yarn test
```

Run tests in a container:

```
docker-compose run --rm frontend-test
```

### backend

Run all tests:

```
cd backend/
go test ./...
```

Update test golden files:

```
cd backend/
go test --update-golden-files ./...
# Files ./testdata/server/golden/*.json have been updated. 
```

Run tests in a container:

```
docker-compose run backend-test
```

### E2E tests

Run tests:

```
cd e2e/
yarn test:e2e
```

## Example

This is the existing SFX service response for The New Yorker:

> http://sfx.library.nyu.edu/sfxlcl41?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat

Ariadne should return the same set of links for use in the frontend:

JSON:
> http://localhost:8080/v0/?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat

## TODO:

- Sanitize urls to prevent XML injection via query string
- Minimize load on the SFX server through backend and/or frontend debouncing, memo-izing, caching, etc.
