import 'bootstrap/dist/css/bootstrap.min.css';

import { useEffect, useRef } from 'react';
import { Col, Container, Row } from 'react-bootstrap';

import linksApi from '../../api/fetchData';
import useApi from '../../hooks/useApi';
import AskLibrarian from '../AskLibrarian/AskLibrarian';
import Citation from '../Citation/Citation';
import Error from '../Error/Error';
import List from '../List/List';
import Loader from '../Loader/Loader';
import StableLink from '../StableLink/StableLink';

const Main = () => {
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
        <Container style={{ "marginBottom": "26px" }}>
          <Row md={12}>
            <Col md={8}>
              <div className="jumbotron">
                <div className="container text-left">
                  <h1>{RESULTS_HEADER_TEXT}</h1>
                  <div>
                    <Citation />
                  </div>
                </div>
              </div>
              {backendClient.loading && <Loader />}
              <div className="mt-3 mb-3"> {/* Add Bootstrap margin-top and margin-bottom classes */}
                <StableLink />
              </div>
              {(backendClient.resource?.length === 0 || backendClient.error) ?
                (
                  <>
                    <div role="alert">
                      {backendClient.error && <Error message={backendClient.error} />}
                    </div>
                    <div>
                      <p>No results found</p>
                    </div>
                  </>) :
                <List links={backendClient.resource} loading={backendClient.loading} />}
            </Col>
            <Col md={4}>
              <aside title="ask-librarian">
                <AskLibrarian loading={backendClient.loading} />
              </aside>
            </Col>
          </Row>
        </Container>
      </>
    </main>
  );
};

export const RESULTS_HEADER_TEXT = 'GetIt Search Results:';

export default Main;
