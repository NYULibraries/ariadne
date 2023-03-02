import 'bootstrap/dist/css/bootstrap.min.css';

import { useEffect, useRef } from 'react';
import { Col, Container, Row } from 'react-bootstrap';

import Citation from '../Citation/Citation';
//import metaData from '../../metadata.json';
import { getCoverageStatement } from '../../aux/helpers';
//import metaData from '../../metadata.json';
import linksApi from '../../api/fetchData';
import useApi from '../../hooks/useApi';

const List = () => {
  const backendClient = useApi(linksApi.fetchData);

  // This is a ref to the backendClient object. It has a property called .current. The value is persisted between renders.
  // useRef doesn’t notify you when its content changes. Mutating the .current property doesn’t cause a re-render. Refs are not counted as dependencies for useEffect
  // https://reactjs.org/docs/hooks-reference.html#useref
  const backendClientRef = useRef(null);
  backendClientRef.current = backendClient;


  // Be aware that this hook will run twice when running in Strict Mode in React 18.
  //
  // Note also that while it's been customary to do data fetching via hooks
  // as low down in the component hierarchy as possible to avoid the unnecessary/undesired
  // child component re-renderings that could happen when fetching data in components
  // higher up in the hierarchy, the React core team and the community are rethinking
  // data fetching best practices, and it might be good to use newer methods in future projects.
  //
  // See this React thread, which has an important 6/22 post from Dan Abramov of
  // the React core team that includes a summary of the issues as well as links
  // to the Beta React Docs outlining the new guidance:
  //
  //   "What is the recommended way to load data for React 18?"
  //   https://www.reddit.com/r/reactjs/comments/vi6q6f/what_is_the_recommended_way_to_load_data_for/
  //
  useEffect(() => {
    backendClientRef.current.fetchResource();
  }, []);

  return (
    <main>
      <>
        {/* TODO: we could put a spinner here: */}
        <Container>
          <Row md={12}>
            <Col md={8}>
              <div className="jumbotron">
                <div className="container text-left">
                  <h1 style={{ "fontWeight": "bold" }}>{RESULTS_HEADER_TEXT}</h1>
                  <div>
                    <Citation />
                  </div>
                </div>
              </div>
              {/* </Col>
            <Col md={8}> */}
              {backendClient.loading && <div className="loading-container">
                <div aria-label="Loading...">{LOADING_TEXT}</div>
              </div>}
              {backendClient.error && <div>{backendClient.error}</div>}
              {/* </Col>
            <Col md={8}> */}
              <div className="list-group">
                <div className="list-group-item list-group-item-action flex-column border-0">
                  {/* <Citation /> */}
                </div>
                {backendClient.resource?.map((link, idx) => (
                  <div key={idx} className="list-group-item list-group-item-action flex-column border-0">
                    <div className="row">
                      <h3>
                        <a href={link.target_url} target="_blank" rel="noopener noreferrer">
                          {link.target_public_name}
                        </a>
                      </h3>
                      <p>{getCoverageStatement(link)}</p>
                    </div>
                  </div>
                ))}
                {(backendClient.resource?.length === 0 || backendClient.error) && (
                  <>
                    <div className="list-group-item list-group-item-action flex-column border-0">
                      <p>No results found</p>
                    </div>
                  </>
                )}
              </div>
            </Col>
            <Col md={4}>
              <aside>
                {!backendClient.loading && (
                  <div className="ask-librarian">
                    <h4>Need help?</h4>
                    <h5>
                      <a href={ASK_LIBRARIAN_URL} target="_blank" rel="noopener noreferrer">
                        {ASK_LIBRARIAN_TEXT}
                      </a>
                    </h5>
                    <md-card-content>
                      <p>Use <a href="https://library.nyu.edu/ask/" target="_blank" className="md-primoExplore-theme" rel="noreferrer">Ask A Librarian</a> or the &quot;Chat with Us&quot; icon at the bottom right corner for any question you have about the Libraries&apos; services.</p>
                      <p>Visit our <a href="https://guides.nyu.edu/online-tutorials/finding-sources" target="_blank" className="md-primoExplore-theme" rel="noreferrer">online tutorials</a> for tips on searching the catalog and getting library resources.</p>
                      <h3 className="md-subhead">Additional Resources</h3>
                      <ul>
                        <li>Use <a href="https://ezborrow.reshare.indexdata.com/" target="_blank" className="md-primoExplore-theme" rel="noreferrer">EZBorrow</a> or <a href="https://library.nyu.edu/services/borrowing/from-non-nyu-libraries/interlibrary-loan/" target="_blank" className="md-primoExplore-theme" rel="noreferrer">InterLibrary Loan (ILL)</a> for materials unavailable at NYU</li>
                        <li>Discover subject specific resources using <a href="http://guides.nyu.edu" target="_blank" className="md-primoExplore-theme" rel="noreferrer">expert curated research guides</a></li>
                        <li>Explore the <a href="https://library.nyu.edu/services/" target="_blank" className="md-primoExplore-theme" rel="noreferrer">complete list of library services</a></li>
                        <li>Reach out to the Libraries on <a href="https://www.instagram.com/nyulibraries/" target="_blank" className="md-primoExplore-theme" rel="noreferrer">our Instagram</a></li>
                        <li>Search <a href="https://www.worldcat.org/search?qt=worldcat_org_all" target="_blank" className="md-primoExplore-theme" rel="noreferrer">WorldCat</a> for items in nearby libraries</li>
                      </ul>
                    </md-card-content>
                  </div>
                )}
              </aside>
            </Col>
          </Row>
        </Container>
      </>
    </main>
  );
};

export const ASK_LIBRARIAN_TEXT = 'Ask a Librarian';
export const ASK_LIBRARIAN_URL = 'https://library.nyu.edu/ask/';
export const LOADING_TEXT = 'Loading...';
export const RESULTS_HEADER_TEXT = 'Getit Search Results:';

export default List;
