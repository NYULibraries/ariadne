import { useEffect, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import { getQueryParameter, getInstitution } from '../../aux/helpers';

const Banner = () => {
  const [logo, setLogo] = useState('https://cdn.library.nyu.edu/images/nyulibraries-logo.svg');
  const [link, setLink] = useState('http://library.nyu.edu');
  const [imgClass, setImgClass] = useState('image');

  useEffect(() => {
    const institution = getQueryParameter('institution');
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
