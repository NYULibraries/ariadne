# SFX API OpenURL response example XML files with errors corrected

Examples of SFX API OpenURL responses of various types is given here:
https://developers.exlibrisgroup.com/sfx/apis/web_services/openurl/

Unfortunately the examples in the documentation have errors, so we are keeping
local corrected examples for reference purposes.

Corrections made:

* Illegal "&" characters have been changed to "&amp;".
* The XML contents of the <ctx_obj_attributes> have been escaped.  In the actual
response returned by the SFX API, the `<perldata>` element inside the `<ctx_obj_attributes>`
element is escaped, so presumably the documentation should have it escaped as well.

Other alterations:

* The examples were pretty-formatted to make them more readable.

Errors not corrected:

* The `<perldata>` XML in _multi_obj_detailed_xml-with-perldata-unescaped_ is not valid.
There appears to have been some mangling of data.  It would be better for ExLibris to fix it.
Some details:
    * Open and close tags do not match.
    * `<itemkey="rft.date">2004` should be `<item key="rft.date">2004`. 
