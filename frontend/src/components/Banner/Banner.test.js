import { render, screen, waitFor } from '@testing-library/react';
import Banner from './Banner';

test('renders the Banner component', () => {
  render(<Banner />);
});

test('renders the NYU Libraries logo', async () => {
  render(<Banner />);
  const linkElement = await waitFor(() => screen.getByAltText(/NYU Libraries logo/i));
  expect(linkElement).toBeInTheDocument();
});

test('renders correctly', () => {
  const { asFragment } = render(<Banner />);
  expect(asFragment()).toMatchSnapshot();
});
