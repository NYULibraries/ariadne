import 'bootstrap/dist/css/bootstrap.min.css';
import { useEffect, useRef } from 'react';
import { getCoverageStatement } from '../../aux/helpers';
import useApi from '../../hooks/useApi';
import linksApi from '../../api/fetchData';

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
      <div className="i-am-centered">
        <div className="list-group">
          {backendClient.resource?.map((link, idx) => (
            <div key={idx} className="list-group-item list-group-item-action flex-column" id="border-style">
              <div className="row">
                <h6>
                  <a href={link.target_url} target="_blank" rel="noopener noreferrer">
                    {link.target_public_name}
                  </a>
                </h6>
                <small>{getCoverageStatement(link)}</small>
              </div>
            </div>
          ))}
          {backendClient.resource?.length === 0 && <div className="i-am-centered">No results found</div>}
        </div>
        {backendClient.resourceLastElement && (
          <div className="ask-librarian">
            <h6>
              <a href={backendClient.resourceLastElement.target_url} target="_blank" rel="noopener noreferrer">
                {backendClient.resourceLastElement.target_public_name}
              </a>
            </h6>
          </div>
        )}
      </div>
    </>
  );
};

export const LOADING_TEXT = 'Loading...';
export const RESULTS_HEADER_TEXT = 'Displaying search results...';
export default List;
