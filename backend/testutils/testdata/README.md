# Test cases

## Chosen by the dev team

* **can-community-task-groups-learn-from-the-principles-of-group-therapy**: requires
  query string `sid` -> `rfr_id` to prevent SFX error "XSS violation occured [sic]."
* **corriere-fiorentino**: a simple, basic test case with a short response 
* **the-new-yorker**: *The New Yorker*, which has a fairly long list of links
* **moral-psychology-is-relationship-regulation**: query string does not have `rft.genre`
  or `genre` in the query string
* **editorial-cartoon**: `genre` is "unknown"

## Test case group for fallback to Primo

* **our-lady-of-everyday-life**: also a targeted test case (see next section).
ISBN-based search handled by Primo.  SFX returns no links.  Primo returns two eBook links.
* **history-today**: ISSN-based search handled by Primo.  SFX returns no links.  Primo returns 5 links.
* **the-sino-tibetan-languages**: neither SFX nor Primo return links.

## Test case groups used by the [sampler](https://github.com/NYULibraries/openurl-link-resolver-sampler)

### [Targeted](https://github.com/NYULibraries/openurl-link-resolver-sampler/blob/e056810c53bcf9fdd5b0232518b9cc5bd9f1b7f9/test-case-files/targeted/targeted-getit-test-OpenURLs.txt)
* **design-of-led-light-therapy-device-based-on-free-form-lens**
* **life-magazine-and-the-power-of-photography**
* **modelling-modular-living**
* **our-lady-of-everyday-life**: also a Primo service test case (see previous section)
* **the-years-work-in-modern-language-studies**
* **publish-the-picture-at-your-peril**
* **publish-the-picture-at-your-peril_unescaped-ampersand** (same as **publish-the-picture-at-your-peril**,
  but with the ampersand character unescaped in the query string)
