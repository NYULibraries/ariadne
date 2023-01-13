import { useEffect, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import shanghaiLogo from './images/shanghai-logo-color.svg';
import abudhabiLogo from './images/abudhabi-logo-color.svg';

const Banner = () => {
  const [logo, setLogo] = useState('https://cdn.library.nyu.edu/images/nyulibraries-logo.svg');
  // const [institution, setInstitution] = useState(null);

  useEffect(() => {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const institution = urlParams.get('institution');
    if (institution === 'NYUAD') {
      setLogo(abudhabiLogo);
    } else if (institution === 'NYUSH') {
      setLogo(shanghaiLogo);
    } else if (institution === 'umlaut.institution') {
      // setInstitution(institution);
      urlParams.delete('umlaut.institution');
      window.location.search = urlParams.toString();
    }
    // setInstitution(institution);
  }, []);
  return (
    <Navbar className="color-nav" expand="lg">
      <Container>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link href="http://library.nyu.edu">
              <img className="image" src={logo} alt="NYU Libraries logo" width="220" height="30" />
            </Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

export default Banner;
