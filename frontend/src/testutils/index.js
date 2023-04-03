import * as fs from 'node:fs';
const path = require('path');

// Prefer actual backend/api/testdata/ files if they exist (i.e. running tests on
// local machine while working in the repo).
const BACKEND_API_TESTDATA_DIR_IN_REPO = path.join(
  __dirname,
  '..',
  '..',
  '..',
  'backend',
  'testutils',
  'testdata',
);
// If in a container, the backend/api/testdata/ files will have been copied into this directory.
const BACKEND_API_TESTDATA_DIR_IN_CONTAINER = path.join(
  __dirname,
  'backend-testdata',
);
const BACKEND_API_TESTDATA_DIR = fs.existsSync(BACKEND_API_TESTDATA_DIR_IN_REPO) ?
  BACKEND_API_TESTDATA_DIR_IN_REPO :
  BACKEND_API_TESTDATA_DIR_IN_CONTAINER;

const BACKEND_API_TEST_CASES_GOLDEN_FILES_DIR = path.join(BACKEND_API_TESTDATA_DIR, 'golden');
const BACKEND_API_TEST_CASES_INDEX = path.join(BACKEND_API_TESTDATA_DIR, 'test-cases.json');

function getTestCasesBackendFetchExceptions() {
  return [
    // Backend is down, or not accessible.
    {
      name: 'TypeError: Failed to fetch',
      // The New Yorker
      queryString: 'ctx_ver=Z39.88-2004&ctx_enc=info:ofi/enc:UTF-8&ctx_tim=2018-07-15T02:13:26IST&url_ver=Z39.88-2004&url_ctx_fmt=infofi/fmt:kev:mtx:ctx&rfr_id=info:sid/primo.exlibrisgroup.com:primo-dedupmrg524707060&rft_val_fmt=info:ofi/fmt:kev:mtx:journal&rft.genre=journal&rft.jtitle=Corriere%20Fiorentino&rft.btitle=Corriere%20Fiorentino&rft.aulast=&rft.aufirst=&rft.auinit=&rft.auinit1=&rft.auinitm=&rft.ausuffix=&rft.au=&rft.aucorp=&rft.volume=&rft.issue=&rft.part=&rft.quarter=&rft.ssn=&rft.spage=&rft.epage=&rft.pages=&rft.artnum=&rft.pub=&rft.place=Italy&rft.issn=&rft.eissn=&rft.isbn=&rft.sici=&rft.coden=&rft_id=info:doi/&rft.object_id=3400000000000901&rft.primo=dedupmrg524707060&rft.eisbn=&rft_dat=<NYUMARCIT>3400000000000901</NYUMARCIT><grp_id>582323038</grp_id><oa></oa><url></url>&rft_id=info:oai/&req.language=eng',
      error: 'TypeError: Failed to fetch',
    },
  ];
}

// At the moment, all HTTP errors are treated basically the same by the frontend:
// codes and error text are simply printed out for the user.  No error-specific
// handling is done.  We test just a few possible error scenarios, but others are
// possible.  Note that these HTTP error statuses are being returned by the Ariadne
// API, which may or may not be simply "passing through" SFX API HTTP errors.
function getTestCasesBackendHttpErrorResponses() {
  return [
    {
      httpErrorCode: 400,
      httpErrorMessage: 'Bad Request',
      // The New Yorker
      queryString: 'ctx_ver=Z39.88-2004&ctx_enc=info:ofi/enc:UTF-8&ctx_tim=2018-07-15T02:13:26IST&url_ver=Z39.88-2004&url_ctx_fmt=infofi/fmt:kev:mtx:ctx&rfr_id=info:sid/primo.exlibrisgroup.com:primo-dedupmrg524707060&rft_val_fmt=info:ofi/fmt:kev:mtx:journal&rft.genre=journal&rft.jtitle=Corriere%20Fiorentino&rft.btitle=Corriere%20Fiorentino&rft.aulast=&rft.aufirst=&rft.auinit=&rft.auinit1=&rft.auinitm=&rft.ausuffix=&rft.au=&rft.aucorp=&rft.volume=&rft.issue=&rft.part=&rft.quarter=&rft.ssn=&rft.spage=&rft.epage=&rft.pages=&rft.artnum=&rft.pub=&rft.place=Italy&rft.issn=&rft.eissn=&rft.isbn=&rft.sici=&rft.coden=&rft_id=info:doi/&rft.object_id=3400000000000901&rft.primo=dedupmrg524707060&rft.eisbn=&rft_dat=<NYUMARCIT>3400000000000901</NYUMARCIT><grp_id>582323038</grp_id><oa></oa><url></url>&rft_id=info:oai/&req.language=eng',
    },
    {
      httpErrorCode: 403,
      httpErrorMessage: 'Forbidden',
      // The New Yorker
      queryString: 'ctx_ver=Z39.88-2004&ctx_enc=info:ofi/enc:UTF-8&ctx_tim=2018-07-15T02:13:26IST&url_ver=Z39.88-2004&url_ctx_fmt=infofi/fmt:kev:mtx:ctx&rfr_id=info:sid/primo.exlibrisgroup.com:primo-dedupmrg524707060&rft_val_fmt=info:ofi/fmt:kev:mtx:journal&rft.genre=journal&rft.jtitle=Corriere%20Fiorentino&rft.btitle=Corriere%20Fiorentino&rft.aulast=&rft.aufirst=&rft.auinit=&rft.auinit1=&rft.auinitm=&rft.ausuffix=&rft.au=&rft.aucorp=&rft.volume=&rft.issue=&rft.part=&rft.quarter=&rft.ssn=&rft.spage=&rft.epage=&rft.pages=&rft.artnum=&rft.pub=&rft.place=Italy&rft.issn=&rft.eissn=&rft.isbn=&rft.sici=&rft.coden=&rft_id=info:doi/&rft.object_id=3400000000000901&rft.primo=dedupmrg524707060&rft.eisbn=&rft_dat=<NYUMARCIT>3400000000000901</NYUMARCIT><grp_id>582323038</grp_id><oa></oa><url></url>&rft_id=info:oai/&req.language=eng',
    },
    {
      httpErrorCode: 408,
      httpErrorMessage: 'Request Timeout',
      // The New Yorker
      queryString: 'ctx_ver=Z39.88-2004&ctx_enc=info:ofi/enc:UTF-8&ctx_tim=2018-07-15T02:13:26IST&url_ver=Z39.88-2004&url_ctx_fmt=infofi/fmt:kev:mtx:ctx&rfr_id=info:sid/primo.exlibrisgroup.com:primo-dedupmrg524707060&rft_val_fmt=info:ofi/fmt:kev:mtx:journal&rft.genre=journal&rft.jtitle=Corriere%20Fiorentino&rft.btitle=Corriere%20Fiorentino&rft.aulast=&rft.aufirst=&rft.auinit=&rft.auinit1=&rft.auinitm=&rft.ausuffix=&rft.au=&rft.aucorp=&rft.volume=&rft.issue=&rft.part=&rft.quarter=&rft.ssn=&rft.spage=&rft.epage=&rft.pages=&rft.artnum=&rft.pub=&rft.place=Italy&rft.issn=&rft.eissn=&rft.isbn=&rft.sici=&rft.coden=&rft_id=info:doi/&rft.object_id=3400000000000901&rft.primo=dedupmrg524707060&rft.eisbn=&rft_dat=<NYUMARCIT>3400000000000901</NYUMARCIT><grp_id>582323038</grp_id><oa></oa><url></url>&rft_id=info:oai/&req.language=eng',
    },
    {
      httpErrorCode: 500,
      httpErrorMessage: 'Internal Server Error',
      // The New Yorker
      queryString: 'ctx_ver=Z39.88-2004&ctx_enc=info:ofi/enc:UTF-8&ctx_tim=2018-07-15T02:13:26IST&url_ver=Z39.88-2004&url_ctx_fmt=infofi/fmt:kev:mtx:ctx&rfr_id=info:sid/primo.exlibrisgroup.com:primo-dedupmrg524707060&rft_val_fmt=info:ofi/fmt:kev:mtx:journal&rft.genre=journal&rft.jtitle=Corriere%20Fiorentino&rft.btitle=Corriere%20Fiorentino&rft.aulast=&rft.aufirst=&rft.auinit=&rft.auinit1=&rft.auinitm=&rft.ausuffix=&rft.au=&rft.aucorp=&rft.volume=&rft.issue=&rft.part=&rft.quarter=&rft.ssn=&rft.spage=&rft.epage=&rft.pages=&rft.artnum=&rft.pub=&rft.place=Italy&rft.issn=&rft.eissn=&rft.isbn=&rft.sici=&rft.coden=&rft_id=info:doi/&rft.object_id=3400000000000901&rft.primo=dedupmrg524707060&rft.eisbn=&rft_dat=<NYUMARCIT>3400000000000901</NYUMARCIT><grp_id>582323038</grp_id><oa></oa><url></url>&rft_id=info:oai/&req.language=eng',
    },
  ];
}

function getTestCasesBackendResponsesIncludeErrors() {
  return [
    {
      name: 'Response includes 2 errors and an empty SFX response',
      // The New Yorker
      queryString: 'ctx_ver=Z39.88-2004&ctx_enc=info:ofi/enc:UTF-8&ctx_tim=2018-07-15T02:13:26IST&url_ver=Z39.88-2004&url_ctx_fmt=infofi/fmt:kev:mtx:ctx&rfr_id=info:sid/primo.exlibrisgroup.com:primo-dedupmrg524707060&rft_val_fmt=info:ofi/fmt:kev:mtx:journal&rft.genre=journal&rft.jtitle=Corriere%20Fiorentino&rft.btitle=Corriere%20Fiorentino&rft.aulast=&rft.aufirst=&rft.auinit=&rft.auinit1=&rft.auinitm=&rft.ausuffix=&rft.au=&rft.aucorp=&rft.volume=&rft.issue=&rft.part=&rft.quarter=&rft.ssn=&rft.spage=&rft.epage=&rft.pages=&rft.artnum=&rft.pub=&rft.place=Italy&rft.issn=&rft.eissn=&rft.isbn=&rft.sici=&rft.coden=&rft_id=info:doi/&rft.object_id=3400000000000901&rft.primo=dedupmrg524707060&rft.eisbn=&rft_dat=<NYUMARCIT>3400000000000901</NYUMARCIT><grp_id>582323038</grp_id><oa></oa><url></url>&rft_id=info:oai/&req.language=eng',
      response: {
        errors: [
          '[ERROR 1]',
          '[ERROR 2]',
        ],
        records: {},
      },
    },
  ];
}

function getTestCasesBackendSuccess() {
  const testCasesFromAriadneBackendApi = require(BACKEND_API_TEST_CASES_INDEX)
    .filter( testCase => testCase.frontendTest );

  testCasesFromAriadneBackendApi.forEach(testCase => {
      testCase.response = getResponse(testCase.key);
  });

  return testCasesFromAriadneBackendApi;
}

function getResponse(key) {
  const goldenFile = path.join(BACKEND_API_TEST_CASES_GOLDEN_FILES_DIR, `${key}.json`);

  return require(goldenFile);
}

export {
  getTestCasesBackendFetchExceptions,
  getTestCasesBackendHttpErrorResponses,
  getTestCasesBackendResponsesIncludeErrors,
  getTestCasesBackendSuccess,
};

export default [
  {
    name: 'Give to the Libraries',
    href: 'https://library.nyu.edu/giving/'
  },
  {
    name: 'Login to NYU Home',
    href: 'http://home.nyu.edu',
    target: '_blank',
    rel: 'noreferrer'
  },
  {
    name: 'Departments',
    href: 'https://library.nyu.edu/departments/'
  },
  {
    name: 'Staff Wiki',
    href: 'https://wiki.library.nyu.edu/',
    target: '_blank',
    rel: 'noreferrer'
  },
  {
    name: 'Staff Directory',
    href: 'https://library.nyu.edu/people/'
  },
  {
    name: 'Status Page',
    href: 'https://nyulibraries.statuspage.io/',
    target: '_blank',
    rel: 'noreferrer'
  },
  {
    name: 'Research Guides',
    href: 'https://guides.nyu.edu/',
    target: '_blank'
  },
  {
    name: 'FAQs',
    href: 'https://library.answers.nyu.edu/',
    target: '_blank',
    rel: 'noreferrer'
  },
  {
    name: 'Career Opportunities',
    href: 'https://library.nyu.edu/about/who-we-are/career-opportunities/'
  },
  {
    name: 'Contact Us',
    href: 'https://library.nyu.edu/contact/'
  },
  {
    name: 'Accessibility',
    href: 'https://www.nyu.edu/footer/accessibility.html',
    target: '_blank',
    rel: 'noreferrer'
  },
  {
    name: 'Subscribe to our email list',
    href: 'https://signup.e2ma.net/signup/1934378/1922970/'
  },
  {
    name: 'Privacy policy',
    href: 'https://library.nyu.edu/privacy-policy/'
  }
];
