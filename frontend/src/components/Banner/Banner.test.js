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

test('renders the correct NYUAD logo and link based on institution query parameter', async () => {
  window.history.pushState({}, null, '/?institution=NYUAD');
  render(<Banner />);
  const linkElement = await waitFor(() => screen.getByAltText(/NYU Libraries logo/i).closest('a'));
  expect(linkElement).toHaveAttribute('href', 'https://nyuad.nyu.edu/en/library.html');
  const imgElement = linkElement.querySelector('img');
  expect(imgElement).toHaveAttribute('src', `/images/abudhabi-logo-color.svg`);
});

test('renders the correct NYUSH logo and link based on institution query parameter', async () => {
  window.history.pushState({}, null, '/?institution=NYUSH');
  render(<Banner />);
  const linkElement = await waitFor(() => screen.getByAltText(/NYU Libraries logo/i).closest('a'));
  expect(linkElement).toHaveAttribute('href', 'https://shanghai.nyu.edu/academics/library');
  const imgElement = linkElement.querySelector('img');
  expect(imgElement).toHaveAttribute('src', `/images/shanghai-logo-color.svg`);
});

test('changes the background of the logo correctly when institution is NYUSH or NYUAD', async () => {
  const institution = 'NYUAD';
  render(<Banner />, {
    route: `?institution=${institution}`,
  });
  const linkElement = await waitFor(() => screen.getByAltText(/NYU Libraries logo/i));
  expect(linkElement).toHaveClass('image white-bg');
});
