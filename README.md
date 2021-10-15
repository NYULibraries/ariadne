# GetIt

## Backend

The GetIt backend is a simple API written in Golang, that takes an [OpenURL](https://biblio.ugent.be/publication/760060/file/760063.pdf) and returns electronic links from an SFX Knowledgebase that represent NYU's e-holdings of the resource identified by the OpenURL. It is essentially an API wrapper for the [SFX LinkResolver](https://exlibrisgroup.com/products/primo-discovery-service/sfx-link-resolver/) and so is itself an OpenURL Link Resolver.

Under the hood it translates the OpenURL from a GET request as querystring parameters to an XML `ContextObject` for posting to the [SFX Web Service](https://developers.exlibrisgroup.com/sfx/apis/web_services/openurl/), and then parses the resulting XML into a JSON string to pass on to the presentation layer.

### Usage

Since the backend is a Golang application, you can run it while developing with:

```
go run server.go
```

But we've also containerized it:

```
docker-compose up backend
# Visit localhost:8080 with an OpenURL, such as 
# http://localhost:8080/?sid=FirstSearch%3AWorldCat&genre=book&title=Fairy+tales&date=1898&aulast=Andersen&aufirst=H&auinitm=C&rfr_id=info%3Asid%2Ffirstsearch.oclc.org%3AWorldCat&rft.genre=book&rft_id=info%3Aoclcnum%2F7675437&rft.aulast=Andersen&rft.aufirst=H&rft.auinitm=C&rft.btitle=Fairy+tales&rft.date=1898&rft.place=Philadelphia&rft.pub=H.+Altemus+Co.&rft.genre=book
```

### Testing

```
docker-compose run backend-test
```

## Frontend



## TODO:

- Better cascading error handling/printing
- Testing