import { render } from '@testing-library/react';
import App from './App';
import List from './components/List/List';
import Banner from './components/Banner/Banner';

test('renders the App component', () => {
  render(<App />);
});

test('renders the Link component', () => {
  render(<List />);
});

test('renders the Banner component', () => {
  render(<Banner />);
});

test('renders correctly', () => {
  const { asFragment } = render(<App />);
  expect(asFragment()).toMatchSnapshot();
});
