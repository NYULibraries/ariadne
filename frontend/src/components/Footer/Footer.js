import './Footer.css';

const PageFooter = () => {

    return (
        <footer className="primary-footer" data-swiftype-index="false">
            <div className="wrap">
                <div className="block-container">
                    <div className="block block--25 footer__give">
                        <a className="button" href="https://library.nyu.edu/giving/">Give to the Libraries</a>
                    </div>
                    <div className="block block--25 footer__menu">
                        <ul className="list" role="menu">
                            <li className="menu__li menu__li--login-to-nyu-home">
                                <a href="http://home.nyu.edu" className="menu__link menu__link--login-to-nyu-home" target="_blank" rel="noreferrer">
                                    Login to NYU Home
                                </a>
                            </li>
                            <li className="menu__li menu__li--departments">
                                <a href="https://library.nyu.edu/departments/" className="menu__link menu__link--departments">
                                    Departments
                                </a>
                            </li>
                            <li className="menu__li menu__li--staff-wiki">
                                <a href="https://wiki.library.nyu.edu/" className="menu__link menu__link--staff-wiki" target="_blank" rel="noreferrer">
                                    Staff Wiki
                                </a>
                            </li>
                            <li className="menu__li menu__li--staff-directory">
                                <a href="https://library.nyu.edu/people/" className="menu__link menu__link--staff-directory">
                                    Staff Directory
                                </a>
                            </li>
                            <li className="menu__li menu__li--status-page">
                                <a href="https://nyulibraries.statuspage.io" className="menu__link menu__link--status-page" target="_blank" rel="noreferrer">
                                    Status Page
                                </a>
                            </li>
                        </ul>
                    </div>
                    <div className="block block--25 footer__menu">
                        <ul className="list" role="menu">
                            <li className="menu__li menu__li--research-guides">
                                <a href="https://guides.nyu.edu/" className="menu__link menu__link--research-guides" target="_blank" rel="noreferrer">
                                    Research Guides
                                </a>
                            </li>
                            <li className="menu__li menu__li--faqs">
                                <a href="http://library.answers.nyu.edu" className="menu__link menu__link--faqs" target="_blank" rel="noreferrer">
                                    FAQs
                                </a>
                            </li>
                            <li className="menu__li menu__li--career-opportunities">
                                <a href="https://library.nyu.edu/about/who-we-are/career-opportunities/" className="menu__link menu__link--career-opportunities">
                                    Career Opportunities
                                </a>
                            </li>
                            <li className="menu__li menu__li--contact-us">
                                <a href="https://library.nyu.edu/contact/" className="menu__link menu__link--contact-us">
                                    Contact Us
                                </a>
                            </li>
                            <li className="menu__li menu__li--accessibility">
                                <a href="https://www.nyu.edu/footer/accessibility.html" className="menu__link menu__link--accessibility" target="_blank" rel="noreferrer">
                                    Accessibility
                                </a>
                            </li>
                        </ul>
                    </div>
                    <div className="block block--25 footer__social">
                        <p>
                            Find out about upcoming programs, events, and resources.<br />
                            <a className="ss-navigateright right" href="https://bit.ly/nyu-liblink">Subscribe to our email list</a>
                        </p>
                        <a href="https://twitter.com/nyulibraries" className="ss-icon" target="_blank" aria-label="Twitter" rel="noreferrer">
                            <img src="/images/twitter.svg" alt="Twitter logo" height="20" />
                        </a>
                        <a href="https://www.facebook.com/nyulibraries" className="ss-icon" target="_blank" aria-label="Facebook" rel="noreferrer">
                            <img src="/images/facebook.svg" alt="Facebook logo" height="20" />
                        </a>
                        <a href="https://www.instagram.com/nyulibraries" className="ss-icon" target="_blank" aria-label="Instagram" rel="noreferrer">
                            <img src="/images/instagram.svg" alt="Instagram logo" height="20" />
                        </a>
                    </div>
                </div>
                <div className="footer__copyright">
                    Unless otherwise noted, all content copyright New York University. All rights reserved.
                    <a href="https://library.nyu.edu/privacy-policy/">Privacy policy</a>
                    <a className="footer__logo" href="https://www.nyu.edu">
                        <img src="/images/nyu-footer-logo.svg" alt="New York University logo" height="27" />
                    </a>
                </div>
            </div>
        </footer>
    );
};

export default PageFooter;
