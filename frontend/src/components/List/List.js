import 'bootstrap/dist/css/bootstrap.min.css';
import { useEffect } from 'react';
import { getCoverageStatement } from '../../aux/helpers';
import useApi from '../../hooks/useApi';
import linksApi from '../../api/fetchData';

const List = () => {
  const backendClient = useApi(linksApi.fetchData);

  useEffect(() => {
    backendClient.fetchResource();
    // We want to fetch the resource only once when component mounts, so we pass an empty array as a second argument
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <>
      <div className="jumbotron">
        <div className="container text-center">
          <p>Displaying search results...</p>
          <p>Note: Alternate titles might have matched your search terms</p>
        </div>
      </div>
      <div className="i-am-centered">
        <div className="list-group">
          {/* TODO: we could put a spinner here: */}
          {backendClient.loading && <div>Loading...</div>}
          {backendClient.error && <div className="i-am-centered">{backendClient.error}</div>}
          {backendClient.resource?.map((link, idx) => (
            <div key={idx} className="list-group-item list-group-item-action flex-column" style={{ border: 'none' }}>
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

export default List;
