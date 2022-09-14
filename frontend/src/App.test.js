import { render, screen } from '@testing-library/react';
import App from './App';

test('renders resolve project on the page', () => {
  render(<App />);
  const linkElement = screen.getByText(/resolve project/i);
  expect(linkElement).toBeInTheDocument();
});

