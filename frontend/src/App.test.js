import { render, screen } from '@testing-library/react';
import App from './App';
import List from './components/List/List';

test('renders the App component', () => {
  render(<App />);
});

test('renders the Link component', () => {
  render(<List />);
});

test('renders resolve project on the page', () => {
  render(<App />);
  const linkElement = screen.getByText(/Displaying search results/i);
  expect(linkElement).toBeInTheDocument();
});
