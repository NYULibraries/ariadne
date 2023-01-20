import { render, screen, waitFor } from '@testing-library/react';
import Banner from './Banner';
import { getInstitution, getParameterFromQueryString } from '../../aux/helpers';

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

/* 
   It is generally not recommended to use the window object in a test environment.
   The window object is not available in all test environments, such as Node.js
   In certain cases mocking methods can be used https://jestjs.io/docs/manual-mocks#mocking-methods-which-are-not-implemented-in-jsdom
   However it's still a work in progress for the Jest team to fully support the window object in the test environment
   Jest runs on Node.js and uses a virtual DOM to simulate a browser environment, so it does not directly interact with a real browser.
   This allows it to run tests in a faster and more consistent environment than running tests in a real browser.
   Similarly, React Testing Library (RTL) is also designed to run in a Node environment and interact with a virtual DOM rather than a real browser.
   */
test('renders the correct NYUAD logo and link based on institution query parameter', async () => {
  window.history.pushState({}, null, '/?institution=NYUAD');
  render(<Banner />);
  const linkElement = await waitFor(() => screen.getByAltText(/NYU Libraries logo/i).closest('a'));
  expect(linkElement).toHaveAttribute('href', 'https://nyuad.nyu.edu/en/library.html');
  const imgElement = linkElement.querySelector('img');
  expect(imgElement).toHaveAttribute('src', `${process.env.PUBLIC_URL}/images/abudhabi-logo-color.svg`);
});

test('renders the correct NYUSH logo and link based on institution query parameter', async () => {
  window.history.pushState({}, null, '/?institution=NYUSH');
  render(<Banner />);
  const linkElement = await waitFor(() => screen.getByAltText(/NYU Libraries logo/i).closest('a'));
  expect(linkElement).toHaveAttribute('href', 'https://shanghai.nyu.edu/academics/library');
  const imgElement = linkElement.querySelector('img');
  expect(imgElement).toHaveAttribute('src', `${process.env.PUBLIC_URL}/images/shanghai-logo-color.svg`);
});

test('redirects correctly when institution query parameter is "umlaut.institution"', async () => {
  const institution = 'NYUSH';
  // Source: https://stackoverflow.com/questions/54090231/how-to-fix-error-not-implemented-navigation-except-hash-changes
  // This is a workaround for to clear `jsdom` error caused by:
  // https://github.com/jsdom/jsdom/blob/16.0.0/lib/jsdom/living/window/navigation.js#L74-L78
  const assignMock = jest.fn();
  delete window.location;
  window.location = { assign: assignMock };
  window.location.search = `?umlaut.institution=${institution}`;
  render(<Banner />);
  const linkElement = screen.getByAltText(/NYU Libraries logo/i).closest('a');
  expect(linkElement).toHaveAttribute('href', 'https://shanghai.nyu.edu/academics/library');
  const imgElement = linkElement.querySelector('img');
  expect(imgElement).toHaveAttribute('src', `${process.env.PUBLIC_URL}/images/shanghai-logo-color.svg`);
});

test('changes the background of the logo correctly when institution is NYUSH or NYUAD', async () => {
  const institution = 'NYUAD';
  render(<Banner />, {
    route: `?institution=${institution}`,
  });
  const linkElement = await waitFor(() => screen.getByAltText(/NYU Libraries logo/i));
  expect(linkElement).toHaveClass('image white-bg');
});

describe('getInstitution', () => {
  test('returns the correct object for institution "NYUAD"', () => {
    const institution = 'NYUAD';
    const result = getInstitution(institution);
    expect(result).toEqual({
      logo: `${process.env.PUBLIC_URL}/images/abudhabi-logo-color.svg`,
      link: 'https://nyuad.nyu.edu/en/library.html',
      imgClass: 'image white-bg',
    });
  });
  test('returns the correct object for institution "NYUSH"', () => {
    const institution = 'NYUSH';
    const result = getInstitution(institution);
    expect(result).toEqual({
      logo: `${process.env.PUBLIC_URL}/images/shanghai-logo-color.svg`,
      link: 'https://shanghai.nyu.edu/academics/library',
      imgClass: 'image white-bg',
    });
  });
  test('returns the default object for NYU institution', () => {
    const institution = 'NYU';
    const result = getInstitution(institution);
    expect(result).toEqual({
      logo: 'https://cdn.library.nyu.edu/images/nyulibraries-logo.svg',
      link: 'http://library.nyu.edu',
      imgClass: 'image',
    });
  });
});

describe('getParameterFromQueryString', () => {
  test('returns the correct value for "institution" query parameter', () => {
    const queryString = '?institution=NYU';
    const parameterName = 'institution';
    const expectedValue = 'NYU';
    const returnedValue = getParameterFromQueryString(queryString, parameterName);
    expect(returnedValue).toBe(expectedValue);
  });

  test('returns the correct value for "umlaut.institution" query parameter', () => {
    const queryString = '?umlaut.institution=NYUAD';
    const parameterName = 'institution';
    const expectedValue = 'NYUAD';
    const returnedValue = getParameterFromQueryString(queryString, parameterName);
    expect(returnedValue).toBe(expectedValue);
  });
});
