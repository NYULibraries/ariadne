# Ariadne backend

The Ariadne backend is an API server written in Go that takes
an [OpenURL](https://biblio.ugent.be/publication/760060/file/760063.pdf) submitted
via query string and returns JSON data containing electronic links from an
SFX Knowledgebase that represent NYU's e-holdings of the resource identified by
the OpenURL. It is essentially an API wrapper for the [SFX LinkResolver](https://exlibrisgroup.com/products/primo-discovery-service/sfx-link-resolver/),
and so is itself an OpenURL Link Resolver.  In the event that SFX does not return
any links, the API server will fall back to the Primo API.

The backend application is a CLI which in addition to providing the REST API server
functionality also includes useful commands for debugging in the terminal.  

## Start the API server

Using `go run`:

```shell
cd backend/
go run . server
```

Building and running the binary:

```shell
cd backend/
go build
./ariadne server
```

Running with a logging disabled:

```shell
cd backend/
go build
# Run `./ariadne server -h` to get a list of all valid `--logging-level` options
./ariadne server --logging-level disabled
```

Running on a different port:

```shell
cd backend/
go build
./ariadne server --port 8081
```

Get help on the `server` command:

```shell
./ariadne server --help
```

Run in a container (port 8080):

```
docker-compose up backend
```

To run a [delve](https://github.com/go-delve/delve) debuggable containerized
instance (port 8080):

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
Showing /app/main.go:5 (PC: 0x922346)
Warning: listing may not match stale executable
     1:	package main
     2:	
     3:	import "ariadne/cmd"
     4:	
     5:	func main() {
     6:		cmd.Execute()
     7:	}
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

## CLI tools for debugging

Get help on `debug` command:

```shell
./ariadne debug --help
```

Get help on a `debug` sub-command (e.g. `params`):

```shell
./ariadne debug params --help
```

### Examples

* Get the SFX HTTP GET request for The New Yorker **(make sure to keep the single-quotes
around the query string argument!)**:

```shell
./ariadne debug sfx-request '?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404<fssessid>0<%2Ffssessid>&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat'
```

* Query string arguments can be stored in files (in this example, _the-new-yorker.txt_),
and then fed to the ariadne command.  This command is exactly the same as the previous one:

```shell
./ariadne debug sfx-request $( < the-new-yorker.txt )
```

* Get the HTTP XML response from SFX:

```shell
./ariadne debug sfx-response $( < the-new-yorker.txt )
```

* Get the list of targets returned by SFX as formatted JSON:

```shell
./ariadne debug sfx-targets $( < the-new-yorker.txt )
```

* Get the initial Primo ISBN member search HTTP request for Hamlet:

```shell
./ariadne debug primo-isbn-search-request $( < hamlet.txt )
```

* Get the Primo FRBR member search HTTP requests for Hamlet (these secondary requests
take time to generate because the response from the initial ISBN search query must
be fetched and analyzed):

```shell
./ariadne debug primo-frbr-member-requests $( < hamlet.txt )
```

* Get the HTTP responses from Primo:

```shell
./ariadne debug primo-responses $( < hamlet.txt )
```

* Get the response bodies from the HTTP response from Primo as formatted JSON:

```shell
./ariadne debug primo-api-responses $( < hamlet.txt )
```

* Get the Ariadne-filtered list of links from Primo as formatted JSON:

```shell
./ariadne debug primo-links $( < hamlet.txt )
```

* Get the API server JSON response:

```shell
./ariadne debug api-json $( < the-new-yorker.txt )
```

* `debug api-json` output can be piped through
   [jq](https://stedolan.github.io/jq/) or
   [fx](https://github.com/antonmedv/fx).  Get list of all links returned by the
previous command example which contain the string "gale" in the `display_name` by
piping through `fx` with the appropriate reducer (this example also prints
"[no coverage information available]" if `coverage_text` is empty, which in the
case of The New Yorker is never the case):

```shell
./ariadne debug api-json $( < the-new-yorker.txt ) | fx '.records[0].links.filter( link => link.display_name.match( new RegExp( "gale", "i" ) ) ).map( link => ( { name: link.display_name, url: link.url, coverage: link.coverage_text ? link.coverage_text : "[no coverage information available]" } ) )'
```

...which produces a JSON array like this:

```json
[
   {
      "name": "Gale General OneFile",
      "url": "http://proxy.library.nyu.edu/login?url=https://link.gale.com/apps/pub/1161/ITOF?u=nysl_me_newyorku",
      "coverage": "Available from 2002/01/14"
   },
   {
      "name": "Gale Literature Resource Center",
      "url": "http://proxy.library.nyu.edu/login?url=https://link.gale.com/apps/pub/1161/LitRC?u=new64731",
      "coverage": "Available from 1978/01/01  until 1978/12/31. Available from 1982/01/01  until 1982/12/31. Available from 1989/01/01  until 1989/12/31. Available from 1996/01/01  until 1996/12/31. Available from 2002/01/01"
   }
]
```

* Get the URL-decoded params from the given query string in a JSON object:

```shell
./ariadne debug params $( < the-new-yorker.txt )
```

## Generate a shell autocompletion script

* Get a list of all shells for which an autocompletion script can be automatically generated:

```shell
./ariadne completion --help
```

* Generate an autocompletion script for Bash

```shell
./ariadne completion bash
```
To immediately enable autocompletion in your current terminal session:

```shell
source <(./ariadne completion bash)
```

For `zsh`, use the commands above and replace `bash` with `zsh`.

## Testing

For documentation on the test cases used by both the backend and the frontend,
see _[backend/testutils/testdata/README.md](backend/testutils/testdata/README.md)_.

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

### Update SFX and Primo response fixture files

These examples use `go run main.go` instead of `./ariadne`, to emphasize that one
wants to make sure to use the most current version of the code to update the fixture
files, which in many cases will be the code in the source code of the working directory.
`./ariadne` can also be used, as long as it is understood that there exists the
possibility the at the binary is the wrong version.

**SFX**

```shell
cd backend/
go run main.go debug sfx-response '?sid=&aulast=Shakespeare&aufirst=William&genre=book&title=The%20Oxford%20Shakespeare:%20Hamlet&date=1987&isbn=9780198129103' > testutils/testdata/fixtures/sfx-fake-responses/hamlet.xml
```

**Primo**

A lot trickier, because there might be multiple Primo responses that need
to be made into fixtures.  In theory, in addition to the response to the initial
ISBN search, there might be multiple FRBR member searches each of which will
generate a response.  In practice, it might be the case that there will always or
almost always be only one subsequent FRBR member search for the active FRBR group.

The example below assumes only 1 FRBR member search response.  The `debug primo-api-responses`
fetches all Primo responses and returns them in an JSON array.  The first `debug primo-api-responses`
command fetches the full response and uses `jq` to extract the first element in
the array and redirects it into the ISBN search response fixture.  The second
command fetches the second element in the array and redirects it into the FRBR
member search fixture file.

The Primo fake and test currently assume only one single FRBR member request and
response for what is designated the active FRBR group (type 5).  If a new test
case is added that requires more than one FRBR member search, we will need to
rewrite the Primo fake to be able to handle that.

```shell
cd backend/
go run main.go debug primo-api-responses '?sid=&aulast=Shakespeare&aufirst=William&genre=book&title=The%20Oxford%20Shakespeare:%20Hamlet&date=1987&isbn=9780198129103' | jq '.[0]' > testutils/testdata/fixtures/primo-fake-responses/hamlet.json
go run main.go debug primo-api-responses '?sid=&aulast=Shakespeare&aufirst=William&genre=book&title=The%20Oxford%20Shakespeare:%20Hamlet&date=1987&isbn=9780198129103' | jq '.[1]' > testutils/testdata/fixtures/primo-fake-responses/frbr-member-search-data/hamlet.json 
```

## Example

This is the existing SFX service response for The New Yorker:

> http://sfx.library.nyu.edu/sfxlcl41?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat

Ariadne backend should return the same set of links in its response:

JSON:
> http://localhost:8080/v0/?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat
