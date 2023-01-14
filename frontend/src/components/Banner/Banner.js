import { useEffect, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';

const Banner = () => {
  const [logo, setLogo] = useState('https://cdn.library.nyu.edu/images/nyulibraries-logo.svg');
  const [link, setLink] = useState('http://library.nyu.edu');
  const [imgClass, setImgClass] = useState('image');

  useEffect(() => {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    let institution = urlParams.get('institution');
    if (urlParams.has('umlaut.institution')) {
      institution = urlParams.get('umlaut.institution');
      urlParams.delete('umlaut.institution');
      urlParams.set('institution', institution);
      window.location.search = urlParams.toString();
    }
    if (institution === 'NYUAD') {
      setLogo(`${process.env.PUBLIC_URL}/images/abudhabi-logo-color.svg`);
      setLink('https://nyuad.nyu.edu/en/library.html');
      setImgClass('image white-bg');
    } else if (institution === 'NYUSH') {
      setLogo(`${process.env.PUBLIC_URL}/images/shanghai-logo-color.svg`);
      setLink('https://shanghai.nyu.edu/academics/library');
      setImgClass('image white-bg');
    }
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
