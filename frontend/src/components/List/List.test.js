import '@testing-library/jest-dom';

import { render, screen } from '@testing-library/react';

import List from './List';

const mockLinks = [
  {
    target_url: 'https://www.example.com',
    target_public_name: 'Example Website'
  },
  {
    target_url: 'https://www.example.org',
    target_public_name: 'Example Organization'
  }
];

describe('List', () => {
  it('renders without error when links prop is passed', () => {
    render(<List links={mockLinks} />);
    const linkElements = screen.getAllByRole('link');
    expect(linkElements.length).toBe(2);
    expect(linkElements[0]).toHaveAttribute('href', 'https://www.example.com');
    expect(linkElements[1]).toHaveAttribute('href', 'https://www.example.org');
  });

  it('renders "No results found" message when links prop is empty', () => {
    render(<List links={[]} />);
    const noResultsMessage = screen.getByText('No results found');
    expect(noResultsMessage).toBeInTheDocument();
  });

  it('renders error message when error prop is passed', () => {
    const errorMessage = 'Failed to load links';
    render(<List error={errorMessage} />);
    const errorElement = screen.getByText(errorMessage);
    expect(errorElement).toBeInTheDocument();
  });

  it('renders loading indicator when loading prop is true', () => {
    render(<List loading={true} />);
    const loadingElement = screen.getByLabelText('Loading...');
    expect(loadingElement).toBeInTheDocument();
  });

  it('matches snapshot when links prop is passed', () => {
    const { container } = render(<List links={mockLinks} />);
    expect(container.firstChild).toMatchSnapshot();
  });

  it('matches snapshot when links prop is empty', () => {
    const { container } = render(<List links={[]} />);
    expect(container.firstChild).toMatchSnapshot();
  });

  it('matches snapshot when error prop is passed', () => {
    const errorMessage = 'Failed to load links';
    const { container } = render(<List error={errorMessage} />);
    expect(container.firstChild).toMatchSnapshot();
  });

  it('matches snapshot when loading prop is true', () => {
    const { container } = render(<List loading={true} />);
    expect(container.firstChild).toMatchSnapshot();
  });

  it('renders the correct link text', () => {
    render(<List links={mockLinks} />);
    const linkElements = screen.getAllByRole('link');
    expect(linkElements.length).toBe(2);
    expect(linkElements[0]).toHaveTextContent('Example Website');
    expect(linkElements[1]).toHaveTextContent('Example Organization');
  });

  it('renders an Error component with the correct message when error prop is passed', () => {
    const errorMessage = 'Failed to load links';
    render(<List error={errorMessage} />);
    const errorComponent = screen.getByRole('alert');
    expect(errorComponent).toBeInTheDocument();
    expect(errorComponent).toHaveTextContent(errorMessage);
  });

  it('renders a Loader component when loading prop is true', () => {
    render(<List loading={true} />);
    const loaderComponent = screen.getByLabelText('Loading...');
    expect(loaderComponent).toBeInTheDocument();
  });
});

// import {
//   getTestCasesBackendFetchExceptions,
//   getTestCasesBackendHttpErrorResponses,
//   getTestCasesBackendResponsesIncludeErrors,
//   getTestCasesBackendSuccess
// } from '../../testutils';

// import apiClient from '../../api/apiClient';
// import { LOADING_TEXT } from '../Loader/Loader';
// import List from './List';

// const LOADING_TEXT_REGEXP = new RegExp(LOADING_TEXT, 'i');

// // Clearing "wrap in act()" warnings
// // =====================================================
// //
// // Even though @testing-library/react supposedly already wraps in `act()`, warnings
// // appeared in test output after `jest.spyOn(apiClient, 'get')...` was introduced
// // for https://nyu-lib.monday.com/boards/765008773/views/61942705.
// //
// // Sample warning:
// //
// // -----BEGIN SAMPLE-----
// // Warning: An update to List inside a test was not wrapped in act(...).
// //
// //     When testing, code that causes React state updates should be wrapped into act(...):
// //
// // act(() => {
// //   /* fire events that update state */
// // });
// // /* assert on the output */
// // -----END SAMPLE-----
// //
// // According to this article:
// //
// //     "React Testing Library and the “not wrapped in act” Errors"
// //     https://davidwcai.medium.com/react-testing-library-and-the-not-wrapped-in-act-errors-491a5629193b
// //
// // ...there are several cases which can generate warnings.  What seem to apply
// // to tests in this file are these two cases:
// //
// // - Case 1: Asynchronous Updates
// // - Case 3: Premature Exit
// //
// // Adding `await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));`
// // to tests that were generating these warnings resolved the issue.  In some cases
// // adding this line was something that would normally be done anyway to avoid race
// // conditions, but in other cases they had to be inserted even when it wouldn't
// // appear to be necessary (e.g. `waitFor` on the next line would prevent a race
// // condition).  These unintuitive cases have been labeled with comments that
// // refer to this long explanatory comment.

// describe('Backend success', () => {
//   describe.each(getTestCasesBackendSuccess())('$name', (testCase) => {
//     beforeEach(() => {
//       delete window.location;
//       window.location = new URL(`${process.env.REACT_APP_API_URL}`);
//       window.history.pushState({}, null, `?${testCase.queryString}`);
//       jest
//         .spyOn(apiClient, 'get')
//         // Even though theoretically we should only need to intercept `apiClient.get`
//         // once using `.mockResolvedValueOnce`, we instead set no limit on
//         // number of interceptions because in React 18, the data fetch hook
//         // gets called twice in Strict Mode. See this Reddit thread:
//         // https://www.reddit.com/r/reactjs/comments/vi6q6f/what_is_the_recommended_way_to_load_data_for/
//         // We need to make sure that the tests always used the fake value, so we
//         // use `.mockResolvedValue`.
//         .mockResolvedValue(
//           new Response(JSON.stringify(testCase.response, null, '    '), { status: 200, statusText: 'OK' })
//         );
//     });

//     afterEach(() => {
//       delete window.location;
//       window.location = new URL('http://localhost:3000');
//       jest.clearAllMocks();
//     });

//     // test('renders correctly', async () => {
//     //   const actual = render(<List />);
//     //   await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
//     //   expect(actual.asFragment()).toMatchSnapshot();
//     // });

//     // test('renders the search results header text', async () => {
//     //   render(<List />);
//     //   // See comment at top of file: 'Clearing "wrap in act()" warnings'
//     //   await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
//     //   const resultsHeaderText = await waitFor(() => screen.getByText(new RegExp(RESULTS_HEADER_TEXT, 'i')));
//     //   expect(resultsHeaderText).toBeInTheDocument();
//     // });

//   });

//   describe('Backend fetch exceptions', () => {
//     describe.each(getTestCasesBackendFetchExceptions())('$name', (testCase) => {
//       beforeEach(() => {
//         delete window.location;
//         window.location = new URL(`${process.env.REACT_APP_API_URL}?${testCase.queryString}`);
//         jest.spyOn(apiClient, 'get').mockImplementation(() => Promise.reject(new TypeError('Failed to fetch')));
//       });

//       afterEach(() => {
//         delete window.location;
//         window.location = new URL('http://localhost:3000');
//         jest.clearAllMocks();
//       });

//       test('renders backend fetch errors correctly', async () => {
//         const actual = render(<List />);
//         await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
//         expect(actual.asFragment()).toMatchSnapshot();
//       });
//     });
//   });

//   describe('Backend HTTP error responses', () => {
//     describe.each(getTestCasesBackendHttpErrorResponses())(
//       `HTTP $httpErrorCode ($httpErrorMessage) error`,
//       (testCase) => {
//         beforeEach(() => {
//           delete window.location;
//           window.location = new URL(`${process.env.REACT_APP_API_URL}?${testCase.queryString}`);
//           jest
//             .spyOn(apiClient, 'get')
//             .mockResolvedValue(
//               new Response(null, { status: testCase.httpErrorCode, statusText: testCase.httpErrorMessage })
//             );
//         });

//         afterEach(() => {
//           delete window.location;
//           window.location = new URL('http://localhost:3000');
//           jest.clearAllMocks();
//         });

//         test(`renders ${LOADING_TEXT}`, async () => {
//           render(<List />);
//           const loadingIndicator = screen.getByText(LOADING_TEXT_REGEXP);
//           expect(loadingIndicator).toBeInTheDocument();
//           // See comment at top of file: 'Clearing "wrap in act()" warnings'
//           await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
//         });

//         test(`${LOADING_TEXT} no longer present in the DOM after loading data`, async () => {
//           const { getByText } = render(<List />);
//           await waitForElementToBeRemoved(() => getByText(LOADING_TEXT_REGEXP));
//         });

//         test(`is rendered correctly`, async () => {
//           const actual = render(<List />);
//           await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
//           expect(actual.asFragment()).toMatchSnapshot();
//         });
//       }
//     );
//   });

//   describe('Backend response includes errors', () => {
//     describe.each(getTestCasesBackendResponsesIncludeErrors())('$name', (testCase) => {
//       beforeEach(() => {
//         delete window.location;
//         window.location = new URL(`${process.env.REACT_APP_API_URL}?${testCase.queryString}`);
//         jest
//           .spyOn(apiClient, 'get')
//           .mockResolvedValue(
//             new Response(JSON.stringify(testCase.response, null, '    '), { status: 200, statusText: 'OK' })
//           );
//       });

//       afterEach(() => {
//         delete window.location;
//         window.location = new URL('http://localhost:3000');
//         jest.clearAllMocks();
//       });

//       test(`renders ${LOADING_TEXT}`, async () => {
//         render(<List />);
//         const loadingIndicator = screen.getByText(LOADING_TEXT_REGEXP);
//         expect(loadingIndicator).toBeInTheDocument();
//         // See comment at top of file: 'Clearing "wrap in act()" warnings'
//         await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
//       });

//       test(`${LOADING_TEXT} no longer present in the DOM after loading data`, async () => {
//         const { getByText } = render(<List />);
//         await waitForElementToBeRemoved(() => getByText(LOADING_TEXT_REGEXP));
//       });

//       test('renders errors included in backend response correctly', async () => {
//         const actual = render(<List />);
//         await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
//         expect(actual.asFragment()).toMatchSnapshot();
//       });
//     });
//   });
// });
