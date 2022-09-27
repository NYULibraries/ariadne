import 'bootstrap/dist/css/bootstrap.min.css';
import { useEffect } from 'react';
import useApi from '../../hooks/useApi';

const List = () => {
  const getLinksApi = useApi();

  useEffect(() => {
    getLinksApi.request();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <>
      <div className="jumbotron" style={{ backgroundColor: '#F7EDA3' }}>
        <div className="container text-center">
          <p>Displaying search results...</p>
          <p>Note: Alternate titles might have matched your search terms</p>
        </div>
      </div>
      <div className="i-am-centered">
        <div className="list-group">
          {/* we could put a spinner here: */}
          {getLinksApi.loading && <div>Loading...</div>}
          {getLinksApi.error && <div className="i-am-centered">{getLinksApi.error}</div>}
          {getLinksApi.elements?.map((element, idx) => (
            <div key={idx} className="list-group-item list-group-item-action flex-column" style={{ border: 'none' }}>
              <div className="d-flex w-100 justify-content-between">
                <h6 className="mb-1">
                  <a
                    style={{ textDecoration: 'none', color: '#6c07ae' }}
                    href={element.target_url}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    {element.target_public_name}
                  </a>
                </h6>
                <small>{element.coverage[0].coverage_text[0].threshold_text[0].coverage_statement}</small>
              </div>
            </div>
          )) ?? <div className="i-am-centered">No results found</div>}
        </div>
        {getLinksApi.lastElement && (
          <div className="ask-librarian">
            <h6>
              <a
                style={{ textDecoration: 'none', color: '#6c07ae' }}
                href={getLinksApi.lastElement.target_url}
                target="_blank"
                rel="noopener noreferrer"
              >
                {getLinksApi.lastElement.target_public_name}
              </a>
            </h6>
          </div>
        )}
      </div>
    </>
  );
};

export default List;
