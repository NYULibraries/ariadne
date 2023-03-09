import { cleanup, render } from '@testing-library/react';

import App from './App';
import Banner from './components/Banner/Banner';
import PageFooter from './components/Footer/Footer';
import Main from './components/Main/Main';

afterEach(cleanup);

test('renders the App component', () => {
  render(<App />);
});

test('renders the List component', () => {
  render(<Main />);
});

test('renders the Banner component', () => {
  render(<Banner />);
});

test('renders the PageFooter component', () => {
  render(<PageFooter />);
});

