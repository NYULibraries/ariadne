# Ariadne backend

The Ariadne backend is an API server written in Go that takes
an [OpenURL](https://biblio.ugent.be/publication/760060/file/760063.pdf) submitted
via query string and returns JSON data containing electronic links from an
SFX Knowledgebase that represent NYU's e-holdings of the resource identified by
the OpenURL. It is essentially an API wrapper for the [SFX LinkResolver](https://exlibrisgroup.com/products/primo-discovery-service/sfx-link-resolver/),
and so is itself an OpenURL Link Resolver.

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

* Get the SFX HTTP POST request for The New Yorker **(make sure to keep the single-quotes
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

* Get the API server JSON response:

```shell
./ariadne debug api-json $( < the-new-yorker.txt )
```

* `debug api-json` output can be piped through
   [jq](https://stedolan.github.io/jq/) or
   [fx](https://github.com/antonmedv/fx).  Get list of all targets returned by the
previous command example which contain the string "gale" by piping through `fx`
with the appropriate reducer:

```shell
./ariadne debug api-json $( < the-new-yorker.txt ) | fx '.records.ctx_obj[0].ctx_obj_targets[0].target.filter( target => target.target_name.match( new RegExp( "gale", "i" ) ) ).map( target => ( { name: target.target_name, url: target.target_url } ) )'
```

...which produces a JSON array like this:

```json
[
  {
    "name": "GALEGROUP_DB_SHAKESPEARE_COLLECTION_PERI",
    "url": "http://proxy.library.nyu.edu/login?url=http://find.galegroup.com/openurl/openurl?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&rft.issn=0028-792X&ctx_enc=info%3Aofi%3Aenc%3AUTF-8&res_id=info%3Asid%2Fgale%3ASHAX&rft.date=2002&req_dat=info%3Asid%2Fgale%3Augnid%3Anew64731&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft.jtitle=New+Yorker"
  },
  {
    "name": "GALE_GENERAL_ONEFILE",
    "url": "http://proxy.library.nyu.edu/login?url=https://link.gale.com/apps/pub/1161/ITOF?u=nysl_me_newyorku"
  },
  {
    "name": "GALE_LITERATURE_RESOURCE_CENTER",
    "url": "http://proxy.library.nyu.edu/login?url=https://link.gale.com/apps/pub/1161/LitRC?u=new64731"
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
see _[backend/api/testdata/server/README.md](backend/api/testdata/server/README.md)_.

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

## Example

This is the existing SFX service response for The New Yorker:

> http://sfx.library.nyu.edu/sfxlcl41?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat

Ariadne backend should return the same set of links in its response:

JSON:
> http://localhost:8080/v0/?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat
