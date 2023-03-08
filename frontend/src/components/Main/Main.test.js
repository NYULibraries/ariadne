import '@testing-library/jest-dom';

import { render, screen, waitFor, waitForElementToBeRemoved } from '@testing-library/react';
import {
    getTestCasesBackendFetchExceptions,
    getTestCasesBackendHttpErrorResponses,
    getTestCasesBackendResponsesIncludeErrors,
    getTestCasesBackendSuccess
} from '../../testutils';
import { ASK_LIBRARIAN_TEXT, ASK_LIBRARIAN_URL } from '../AskLibrarian/AskLibrarian';
import Main, { RESULTS_HEADER_TEXT } from './Main';

import apiClient from '../../api/apiClient';
import { LOADING_TEXT } from '../Loader/Loader';

const LOADING_TEXT_REGEXP = new RegExp(LOADING_TEXT, 'i');

// Clearing "wrap in act()" warnings
// =====================================================
//
// Even though @testing-library/react supposedly already wraps in `act()`, warnings
// appeared in test output after `jest.spyOn(apiClient, 'get')...` was introduced
// for https://nyu-lib.monday.com/boards/765008773/views/61942705.
//
// Sample warning:
//
// -----BEGIN SAMPLE-----
// Warning: An update to List inside a test was not wrapped in act(...).
//
//     When testing, code that causes React state updates should be wrapped into act(...):
//
// act(() => {
//   /* fire events that update state */
// });
// /* assert on the output */
// -----END SAMPLE-----
//
// According to this article:
//
//     "React Testing Library and the “not wrapped in act” Errors"
//     https://davidwcai.medium.com/react-testing-library-and-the-not-wrapped-in-act-errors-491a5629193b
//
// ...there are several cases which can generate warnings.  What seem to apply
// to tests in this file are these two cases:
//
// - Case 1: Asynchronous Updates
// - Case 3: Premature Exit
//
// Adding `await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));`
// to tests that were generating these warnings resolved the issue.  In some cases
// adding this line was something that would normally be done anyway to avoid race
// conditions, but in other cases they had to be inserted even when it wouldn't
// appear to be necessary (e.g. `waitFor` on the next line would prevent a race
// condition).  These unintuitive cases have been labeled with comments that
// refer to this long explanatory comment.

describe('Backend success', () => {
    describe.each(getTestCasesBackendSuccess())('$name', (testCase) => {
        beforeEach(() => {
            delete window.location;
            window.location = new URL(`${process.env.REACT_APP_API_URL}`);
            window.history.pushState({}, null, `?${testCase.queryString}`);
            jest
                .spyOn(apiClient, 'get')
                // Even though theoretically we should only need to intercept `apiClient.get`
                // once using `.mockResolvedValueOnce`, we instead set no limit on
                // number of interceptions because in React 18, the data fetch hook
                // gets called twice in Strict Mode. See this Reddit thread:
                // https://www.reddit.com/r/reactjs/comments/vi6q6f/what_is_the_recommended_way_to_load_data_for/
                // We need to make sure that the tests always used the fake value, so we
                // use `.mockResolvedValue`.
                .mockResolvedValue(
                    new Response(JSON.stringify(testCase.response, null, '    '), { status: 200, statusText: 'OK' })
                );
        });

        afterEach(() => {
            delete window.location;
            window.location = new URL('http://localhost:3000');
            jest.clearAllMocks();
        });

        test(`renders ${LOADING_TEXT}`, async () => {
            render(<Main />);
            const loadingIndicator = screen.getByText(LOADING_TEXT_REGEXP);
            expect(loadingIndicator).toBeInTheDocument();
            // See comment at top of file: 'Clearing "wrap in act()" warnings'
            await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
        });

        test(`${LOADING_TEXT} no longer present in the DOM after loading data`, async () => {
            const { getByText } = render(<Main />);
            await waitForElementToBeRemoved(() => getByText(LOADING_TEXT_REGEXP));
        });

        test('renders correctly', async () => {
            const actual = render(<Main />);
            await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
            expect(actual.asFragment()).toMatchSnapshot();
        });

        test('renders the search results header text', async () => {
            render(<Main />);
            // See comment at top of file: 'Clearing "wrap in act()" warnings'
            await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
            const resultsHeaderText = await waitFor(() => screen.getByText(new RegExp(RESULTS_HEADER_TEXT, 'i')));
            expect(resultsHeaderText).toBeInTheDocument();
        });

        test('renders with a className of list-group', async () => {
            const { container } = render(<Main />);
            await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
            expect(container.getElementsByClassName('list-group').length).toBe(1);
        });

        test(`renders ${ASK_LIBRARIAN_TEXT}`, async () => {
            render(<Main />);
            // See comment at top of file: 'Clearing "wrap in act()" warnings'
            await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
            const askALibrarian = await waitFor(() => screen.getByText(ASK_LIBRARIAN_TEXT));
            expect(askALibrarian).toBeInTheDocument();
        });

        test(`renders ${ASK_LIBRARIAN_TEXT} with a link to ${ASK_LIBRARIAN_URL}`, async () => {
            render(<Main />);
            await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
            const askALibrarian = await waitFor(() => screen.getByText(ASK_LIBRARIAN_TEXT));
            expect(askALibrarian).toHaveAttribute('href', ASK_LIBRARIAN_URL);
        });
    });
});

describe('Backend fetch exceptions', () => {
    describe.each(getTestCasesBackendFetchExceptions())('$name', (testCase) => {
        beforeEach(() => {
            delete window.location;
            window.location = new URL(`${process.env.REACT_APP_API_URL}?${testCase.queryString}`);
            jest.spyOn(apiClient, 'get').mockImplementation(() => Promise.reject(new TypeError('Failed to fetch')));
        });

        afterEach(() => {
            delete window.location;
            window.location = new URL('http://localhost:3000');
            jest.clearAllMocks();
        });

        test('renders backend fetch errors correctly', async () => {
            const actual = render(<Main />);
            await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
            expect(actual.asFragment()).toMatchSnapshot();
        });
    });
});

describe('Backend HTTP error responses', () => {
    describe.each(getTestCasesBackendHttpErrorResponses())(
        `HTTP $httpErrorCode ($httpErrorMessage) error`,
        (testCase) => {
            beforeEach(() => {
                delete window.location;
                window.location = new URL(`${process.env.REACT_APP_API_URL}?${testCase.queryString}`);
                jest
                    .spyOn(apiClient, 'get')
                    .mockResolvedValue(
                        new Response(null, { status: testCase.httpErrorCode, statusText: testCase.httpErrorMessage })
                    );
            });

            afterEach(() => {
                delete window.location;
                window.location = new URL('http://localhost:3000');
                jest.clearAllMocks();
            });

            test(`renders ${LOADING_TEXT}`, async () => {
                render(<Main />);
                const loadingIndicator = screen.getByText(LOADING_TEXT_REGEXP);
                expect(loadingIndicator).toBeInTheDocument();
                // See comment at top of file: 'Clearing "wrap in act()" warnings'
                await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
            });

            test(`${LOADING_TEXT} no longer present in the DOM after loading data`, async () => {
                const { getByText } = render(<Main />);
                await waitForElementToBeRemoved(() => getByText(LOADING_TEXT_REGEXP));
            });

            test(`is rendered correctly`, async () => {
                const actual = render(<Main />);
                await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
                expect(actual.asFragment()).toMatchSnapshot();
            });
        }
    );
});

describe('Backend response includes errors', () => {
    describe.each(getTestCasesBackendResponsesIncludeErrors())('$name', (testCase) => {
        beforeEach(() => {
            delete window.location;
            window.location = new URL(`${process.env.REACT_APP_API_URL}?${testCase.queryString}`);
            jest
                .spyOn(apiClient, 'get')
                .mockResolvedValue(
                    new Response(JSON.stringify(testCase.response, null, '    '), { status: 200, statusText: 'OK' })
                );
        });

        afterEach(() => {
            delete window.location;
            window.location = new URL('http://localhost:3000');
            jest.clearAllMocks();
        });

        test(`renders ${LOADING_TEXT}`, async () => {
            render(<Main />);
            const loadingIndicator = screen.getByText(LOADING_TEXT_REGEXP);
            expect(loadingIndicator).toBeInTheDocument();
            // See comment at top of file: 'Clearing "wrap in act()" warnings'
            await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
        });

        test(`${LOADING_TEXT} no longer present in the DOM after loading data`, async () => {
            const { getByText } = render(<Main />);
            await waitForElementToBeRemoved(() => getByText(LOADING_TEXT_REGEXP));
        });

        test('renders errors included in backend response correctly', async () => {
            const actual = render(<Main />);
            await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
            expect(actual.asFragment()).toMatchSnapshot();
        });
    });
});
