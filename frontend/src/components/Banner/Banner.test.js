import { render, screen, waitFor } from '@testing-library/react';

import Banner from './Banner';

import { bannerInstitutionInfo } from '../../aux/institutionConstants';

const institutionNamesUpperCase = Object.keys(bannerInstitutionInfo).map(institutionName => institutionName.toUpperCase());

describe.each(institutionNamesUpperCase)(
  'Institution name: %s', (institutionNameUpperCase) => {
    beforeEach( () => {
      delete window.location;
      window.location = new URL(`${process.env.REACT_APP_API_URL}?institution=${institutionNameUpperCase}`);
    });

    afterEach(() => {
      delete window.location;
      window.location = new URL(process.env.REACT_APP_API_URL);
      jest.clearAllMocks();
    });

    test(`renders ${institutionNameUpperCase} page correctly`, () => {
      const { asFragment } = render(<Banner />);
      expect(asFragment()).toMatchSnapshot();
    });

    test(`renders the correct ${institutionNameUpperCase} logo and link based on institution query parameter`, async () => {
      render(<Banner />);
      const linkElement = await waitFor(() => screen.getByAltText(/NYU Libraries logo/i).closest('a'));
      const { logo, link } = bannerInstitutionInfo[institutionNameUpperCase.toLowerCase()]
      expect(linkElement).toHaveAttribute('href', link);
      const imgElement = linkElement.querySelector('img');
      expect(imgElement).toHaveAttribute('src', logo);
    });

    test(`sets the background of the logo correctly for ${institutionNameUpperCase}`, async () => {
      render(<Banner />, {
                           route: `?institution=${institutionNameUpperCase}`,
                         });
      const linkElement = await waitFor(() => screen.getByAltText(/NYU Libraries logo/i));
      const { imgClass } = bannerInstitutionInfo[institutionNameUpperCase.toLowerCase()]
      expect(linkElement).toHaveClass(imgClass);
    });
  }
)

describe('No `institution` parameter in query string', () => {
  test('renders the Banner component', () => {
    render(<Banner />);
  });

  test('renders the NYU Libraries logo', async () => {
    render(<Banner />);
    const linkElement = await waitFor(() => screen.getByAltText(/NYU Libraries logo/i));
    expect(linkElement).toBeInTheDocument();
  });
});
