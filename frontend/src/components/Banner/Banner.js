import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import { getInstitutionQueryParameter } from '../../aux/helpers';
import { bannerInstitutionInfo } from '../../aux/institutionInfo';

const Banner = () => {
  const institutionName = getInstitutionQueryParameter();
  const { logo, link, imgClass, altLibraryLogoImageText } = bannerInstitutionInfo[institutionName] || bannerInstitutionInfo.nyu;
  return (
    <Navbar className="color-nav" expand="lg" role="navigation">
      <Container>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link href={link}>
              <img
                className={imgClass}
                src={logo}
                alt={altLibraryLogoImageText}
                width="220"
                height="60"
              />
            </Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

export default Banner;
