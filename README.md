# Resolve

## Backend

The Resolve backend is a simple API written in Golang, that takes an [OpenURL](https://biblio.ugent.be/publication/760060/file/760063.pdf) and returns electronic links from an SFX Knowledgebase that represent NYU's e-holdings of the resource identified by the OpenURL. It is essentially an API wrapper for the [SFX LinkResolver](https://exlibrisgroup.com/products/primo-discovery-service/sfx-link-resolver/) and so is itself an OpenURL Link Resolver.

Under the hood it translates the OpenURL from a GET request as querystring parameters to an XML `ContextObject` for posting to the [SFX Web Service](https://developers.exlibrisgroup.com/sfx/apis/web_services/openurl/), and then parses the resulting XML into a JSON string to pass on to the presentation layer.

### Usage

Since the backend is a Golang application, you can run it while developing with:

```
go run server.go
```

But it's also containerized:

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

The frontend takes OpenURLs as querystring params and makes async calls to the backend with the OpenURL. Should parse the returned JSON to just pull out the target URLs with display text and coverage text if available.

For instance this is the existing service for The New Yorker:

> http://sfx.library.nyu.edu/sfxlcl41?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat

And this should return the same set of links for use in the frontend:

> http://localhost:8080/?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat

## TODO:

- Use Gin or other more performant http framework
- Testing 
  - round out request coverage - ~80%
  - test server.go
- GraphQL?
- Frontend 
  - react? vue?
  - Debounce, memoize and otherwise cache requests so as not to bombard the SFX server