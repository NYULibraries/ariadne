import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';

const Banner = () => {
  return (
    <Navbar bg="#﻿57068c" expand="lg">
      <Container>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link href="http://library.nyu.edu">
              <img
                className="image"
                src="https://cdn.library.nyu.edu/images/nyulibraries-logo.svg"
                alt="NYU Libraries logo"
                width="220"
                height="30"
              />
            </Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

export default Banner;
