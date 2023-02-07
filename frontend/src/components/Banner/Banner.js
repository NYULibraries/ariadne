import { DEFAULT_IMG_CLASS, DEFAULT_LINK, DEFAULT_LOGO } from '../../aux/institutionConstants';
import { getInstitution, getInstitutionQueryParameter } from '../../aux/helpers';
import { useEffect, useState } from 'react';

import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';

const Banner = () => {
  const [logo, setLogo] = useState(DEFAULT_LOGO);
  const [link, setLink] = useState(DEFAULT_LINK);
  const [imgClass, setImgClass] = useState(DEFAULT_IMG_CLASS);

  useEffect(() => {
    const institution = getInstitutionQueryParameter('institution');
    const { logo, link, imgClass } = getInstitution(institution);
    setLogo(logo);
    setLink(link);
    setImgClass(imgClass);
  }, []);
  return (
    <Navbar className="color-nav" expand="lg">
      <Container>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link href={link}>
              <img className={imgClass} src={logo} alt="NYU Libraries logo" width="220" height="60" />
            </Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

export default Banner;
