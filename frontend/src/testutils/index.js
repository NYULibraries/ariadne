const fs = require('fs');
const path = require('path');

// Prefer actual backend/api/testdata/ files if they exist (i.e. running tests on
// local machine while working in the repo).  If in a container, the backend/api/testdata/
// files will have been copied into this directory.
const BACKEND_API_TESTDATA_DIR_IN_REPO = path.join(
    __dirname,
    '..',
    '..',
    '..',
    'backend',
    'api',
    'testdata',
    'server',
);
const BACKEND_API_TESTDATA_DIR_IN_CONTAINER = path.join(
    __dirname,
    'backend-testdata',
);
const BACKEND_API_TESTDATA_DIR = fs.existsSync(BACKEND_API_TESTDATA_DIR_IN_REPO) ?
                                 BACKEND_API_TESTDATA_DIR_IN_REPO :
                                 BACKEND_API_TESTDATA_DIR_IN_CONTAINER;

const BACKEND_API_TEST_CASES_GOLDEN_FILES_DIR = path.join(BACKEND_API_TESTDATA_DIR, 'golden');
const BACKEND_API_TEST_CASES_INDEX = path.join(BACKEND_API_TESTDATA_DIR, 'test-cases.json');

function getTestCases() {
    const testCasesFromAriadneBackendApi = require(BACKEND_API_TEST_CASES_INDEX);

    testCasesFromAriadneBackendApi.forEach( testCase => {
        testCase.response = getResponse(testCase.key);
    });

    return testCasesFromAriadneBackendApi;
}

function getResponse(key) {
    const goldenFile = path.join(BACKEND_API_TEST_CASES_GOLDEN_FILES_DIR, `${key}.json`);

    return require(goldenFile);
}

export {
    getTestCases,
};
