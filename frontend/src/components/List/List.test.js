import { render, screen, waitFor, waitForElementToBeRemoved } from '@testing-library/react';
import '@testing-library/jest-dom';
import List from './List';

test('renders with a className of list-group', () => {
  const { container } = render(<List />);
  expect(container.getElementsByClassName('list-group').length).toBe(1);
});

test('render correctly', () => {
  const { asFragment } = render(<List />);
  expect(asFragment()).toMatchSnapshot();
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

test('renders a E Journal Full Text link', async () => {
  render(<List />);
  const linkElement = await waitFor(() => screen.getByText(/E Journal Full Text/i));
  expect(linkElement).toBeInTheDocument();
});

test('renders a Gale General OneFile link', async () => {
  render(<List />);
  const linkElement = await waitFor(() => screen.getByText(/Gale General OneFile/i));
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
