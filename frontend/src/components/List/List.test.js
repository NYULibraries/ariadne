import { render, screen, waitFor, waitForElementToBeRemoved } from '@testing-library/react';
import '@testing-library/jest-dom';
import List, {LOADING_TEXT} from './List';
import { getTestCases } from '../../testutils';
import apiClient from '../../api/apiClient';

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

const testCases = getTestCases();
testCases.forEach( testCase => {
  describe(testCase.name, () => {

    beforeEach(() => {
      delete window.location;
      window.location = new URL(`${process.env.REACT_APP_API_URL}?${testCase.queryString}`);
      jest.spyOn(apiClient, 'get').mockResolvedValue( { data: testCase.response });
    });

    afterEach(() => {
      delete window.location;
      window.location = new URL('http://localhost:3000');
      jest.clearAllMocks();
    });

    test('renders with a className of list-group', async () => {
      const { container } = render(<List />);
      await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
      expect(container.getElementsByClassName('list-group').length).toBe(1);
    });

    test.skip('renders correctly', async () => {
      const actual = render(<List />);
      await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
      expect(actual.asFragment()).toMatchSnapshot();
    });

    test('renders the search results', async () => {
      render(<List />);
      // See comment at top of file: 'Clearing "wrap in act()" warnings'
      await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
      const linkElement = await waitFor(() => screen.getByText(/Displaying search results.../i));
      expect(linkElement).toBeInTheDocument();
    });

    test(`renders ${LOADING_TEXT}`, async () => {
      render(<List />);
      const linkElement = screen.getByText(LOADING_TEXT_REGEXP);
      expect(linkElement).toBeInTheDocument();
      // See comment at top of file: 'Clearing "wrap in act()" warnings'
      await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
    });

    test('renders Ask a Librarian', async () => {
      render(<List />);
      // See comment at top of file: 'Clearing "wrap in act()" warnings'
      await waitForElementToBeRemoved(() => screen.getByText(LOADING_TEXT_REGEXP));
      const linkElement = await waitFor(() => screen.getByText(/Ask a Librarian/i));
      expect(linkElement).toBeInTheDocument();
    });

    test(`${LOADING_TEXT} no longer present in the DOM after loading data`, async () => {
      const { getByText } = render(<List />);
      await waitForElementToBeRemoved(() => getByText(LOADING_TEXT_REGEXP));
    });

  });
});
