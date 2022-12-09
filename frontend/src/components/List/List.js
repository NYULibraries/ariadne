import 'bootstrap/dist/css/bootstrap.min.css';
import { useEffect, useRef } from 'react';
import { getCoverageStatement } from '../../aux/helpers';
import useApi from '../../hooks/useApi';
import linksApi from '../../api/fetchData';
import { Col, Container, Row } from 'react-bootstrap';

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
  // eslint-disable-next-line no-console
  console.log('backendClient', backendClient.resource);

  return (
    <>
      <div className="jumbotron">
        <div className="container text-center">
          <p>{RESULTS_HEADER_TEXT}</p>
          <p>Note: Alternate titles might have matched your search terms</p>
        </div>
      </div>
      {/* TODO: we could put a spinner here: */}
      {backendClient.loading && <div className="loader">{LOADING_TEXT}</div>}
      {backendClient.error && <div className="i-am-centered">{backendClient.error}</div>}
      <Container>
        <Row md={12}>
          <Col md={8}>
            <div className="list-group">
              {backendClient.resource?.map((link, idx) => (
                <div key={idx} className="list-group-item list-group-item-action flex-column border-0">
                  <div className="row">
                    <span>
                      {/* Add metadata placeholders here */}
                      <p className="resource-type">{link.genre}</p>
                      <h2 className="title">{link.article_title}</h2>
                      <p>
                        {link.author} <span>•</span> {link.date}
                      </p>
                      <p style={{ margin: '0 0 10px' }}>
                        <span style={{ boxSizing: 'border-box' }}>Published in Journal </span>
                        <span style={{ fontStyle: 'italic' }}>{link.journal_title}. </span>
                        Volume {link.volume}. Issue {link.issue}. Page {link.start_page}-{link.end_page}.
                      </p>
                      <dl className="citation-info">
                        <dt>ISSN:</dt>
                        <dd>{link.issn}</dd>
                      </dl>
                    </span>
                    <h6>
                      <a href={link.target_url} target="_blank" rel="noopener noreferrer">
                        {link.target_public_name}
                      </a>
                    </h6>
                    <small>{getCoverageStatement(link)}</small>
                  </div>
                </div>
              ))}
              {backendClient.resource?.length === 0 && <p>No results found</p>}
            </div>
          </Col>
          <Col md={4}>
            {!backendClient.loading && (
              <div className="ask-librarian">
                <h5>Need help?</h5>
                <h6>
                  <a href={ASK_LIBRARIAN_URL} target="_blank" rel="noopener noreferrer">
                    {ASK_LIBRARIAN_TEXT}
                  </a>
                </h6>
              </div>
            )}
          </Col>
        </Row>
      </Container>
    </>
  );
};

export const ASK_LIBRARIAN_TEXT = 'Ask a Librarian';
export const ASK_LIBRARIAN_URL = 'https://library.nyu.edu/ask/';
export const LOADING_TEXT = 'Loading...';
export const RESULTS_HEADER_TEXT = 'Displaying search results...';
export default List;
