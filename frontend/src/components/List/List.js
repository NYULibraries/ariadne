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

  useEffect(() => {
    backendClientRef.current.fetchResource();
  }, []);

  return (
    <>
      <div className="jumbotron">
        <div className="container text-center">
          <p>Displaying search results...</p>
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
export default List;
