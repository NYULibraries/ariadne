import { render, screen, waitFor, waitForElementToBeRemoved } from '@testing-library/react';
import '@testing-library/jest-dom';
import List from './List';
import { getTestCases } from '../../testutils';

const testCases = getTestCases();
testCases.forEach( testCase => {
  describe(testCase.name, () => {

    beforeEach(() => {
      delete window.location;
      window.location = new URL(`${process.env.REACT_APP_API_URL}?${testCase.queryString}`);
    });

    afterEach(() => {
      delete window.location;
      window.location = new URL('http://localhost:3000');
    });

    test('renders with a className of list-group', () => {
      const { container } = render(<List />);
      expect(container.getElementsByClassName('list-group').length).toBe(1);
    });

    test.skip('renders correctly', async () => {
      const actual = render(<List />);
      await waitForElementToBeRemoved(() => screen.getByText(/Loading.../i));
      expect(actual.asFragment()).toMatchSnapshot();
    });

    test('renders the search results', async () => {
      render(<List />);
      const linkElement = await waitFor(() => screen.getByText(/Displaying search results.../i));
      expect(linkElement).toBeInTheDocument();
    });

    test('renders Loading...', () => {
      render(<List />);
      const linkElement = screen.getByText(/Loading/i);
      expect(linkElement).toBeInTheDocument();
    });

    test('renders Ask a Librarian', async () => {
      render(<List />);
      const linkElement = await waitFor(() => screen.getByText(/Ask a Librarian/i));
      expect(linkElement).toBeInTheDocument();
    });

    test('Loading... no longer present in the DOM after loading data', async () => {
      const { getByText } = render(<List />);
      await waitForElementToBeRemoved(() => getByText(/Loading/i));
    });

  });
});
