import React from 'react';
import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import { getInstitutionQueryParameter } from '../../aux/helpers';
import { bannerInstitutionInfo } from '../../aux/institutionInfo';

const Banner = () => {
  const institutionName = getInstitutionQueryParameter();
  const { logo, link, imgClass, altLibraryLogoImageText } = bannerInstitutionInfo[institutionName] || bannerInstitutionInfo.nyu;
  return (
    <Navbar className="color-nav" expand="lg" role="banner">
      <Container>
        <Nav className="me-auto" role="navigation">
          <Nav.Link href={link} aria-label={`${institutionName} home`}>
            <img
              className={imgClass}
              src={logo}
              alt={altLibraryLogoImageText}
              width="220"
              height="60"
            />
          </Nav.Link>
        </Nav>
      </Container>
    </Navbar>
  );
};

export default Banner;
