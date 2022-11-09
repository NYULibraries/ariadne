const path = require('path');

const BACKEND_API_TESTDATA = path.join(
    __dirname,
    '..',
    '..',
    '..',
    'backend',
    'api',
    'testdata',
    'server',
);
const BACKEND_API_TEST_CASES_GOLDEN_FILES_DIR = path.join(BACKEND_API_TESTDATA, 'golden');
const BACKEND_API_TEST_CASES_INDEX = path.join(BACKEND_API_TESTDATA, 'test-cases.json');

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
