# Test cases

## Chosen or created by the dev team

* **can-community-task-groups-learn-from-the-principles-of-group-therapy**: requires
  query string `sid` -> `rfr_id` to prevent SFX error "XSS violation occured [sic]."
* **contrived-frbr-group-test-case**: contrived test case to thoroughly exercise
  the `primo` package code.  Obviously also a Primo service test case (see next section).
* **corriere-fiorentino**: a simple, basic test case with a short response 
* **editorial-cartoon**: `genre` is "unknown"
* **history-today**: ISSN-based search for which Ariadne was originally incorrectly
constructing the SFX query due to testing only for the absence of `date` query param
and not testing for the existence of the `date` param with an empty value.
* **moral-psychology-is-relationship-regulation**: query string does not have `rft.genre`
  or `genre` in the query string
* **the-new-yorker**: *The New Yorker*, which has a fairly long list of links

## Test case group for fallback to Primo

* **contrived-frbr-group-test-case**: contrived test case created by the dev team
  to thoroughly exercise the `primo` package code.
* **hamlet**: ISBN-based search handled by Primo which has a FRBR type of "5",
  requiring extra processing including a second call the Primo API.
* **our-lady-of-everyday-life**: also a targeted test case (see next section).
ISBN-based search handled by Primo.  SFX returns no links.  Primo returns two eBook links.
* **the-sino-tibetan-languages**: neither SFX nor Primo return links.

## Test case groups used by the [sampler](https://github.com/NYULibraries/openurl-link-resolver-sampler)

### [Targeted](https://github.com/NYULibraries/openurl-link-resolver-sampler/blob/e056810c53bcf9fdd5b0232518b9cc5bd9f1b7f9/test-case-files/targeted/targeted-getit-test-OpenURLs.txt)
* **design-of-led-light-therapy-device-based-on-free-form-lens**
* **life-magazine-and-the-power-of-photography**
* **modelling-modular-living**
* **our-lady-of-everyday-life**: also a Primo service test case (see previous section)
* **publish-the-picture-at-your-peril**
* **publish-the-picture-at-your-peril_unescaped-ampersand** (same as **publish-the-picture-at-your-peril**,
  but with the ampersand character unescaped in the query string)
* **the-years-work-in-modern-language-studies**
