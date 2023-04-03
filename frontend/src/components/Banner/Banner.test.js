import { render, screen, waitFor } from '@testing-library/react';

import { bannerInstitutionInfo } from '../../aux/institutionInfo';
import Banner from './Banner';

const institutionNamesUpperCase = Object.keys(bannerInstitutionInfo).map(institutionName => institutionName.toUpperCase());

describe.each(institutionNamesUpperCase)(
  'Institution name: %s', (institutionNameUpperCase) => {
    beforeEach(() => {
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
      const { altLibraryLogoImageText } = bannerInstitutionInfo[institutionNameUpperCase.toLowerCase()]
      const linkElement = await waitFor(() => screen.getByAltText(altLibraryLogoImageText).closest('a'));
      const { logo, link } = bannerInstitutionInfo[institutionNameUpperCase.toLowerCase()]
      expect(linkElement).toHaveAttribute('href', link);
      const imgElement = linkElement.querySelector('img');
      expect(imgElement).toHaveAttribute('src', logo);
    });

    test(`sets the background of the logo correctly for ${institutionNameUpperCase}`, async () => {
      render(<Banner />, {
        route: `?institution=${institutionNameUpperCase}`,
      });
      const { altLibraryLogoImageText } = bannerInstitutionInfo[institutionNameUpperCase.toLowerCase()]
      const linkElement = await waitFor(() => screen.getByAltText(altLibraryLogoImageText));
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
    const linkElement = await waitFor(() => screen.getByAltText('NYU Libraries homepage.').closest('a'));
    expect(linkElement).toBeInTheDocument();
  });
});
