import 'bootstrap/dist/css/bootstrap.min.css';
import { useEffect, useState } from 'react';

const List = () => {
  const [elements, setElements] = useState([]);

  const fetchElements = async () => {
    try {
      const response = await fetch('the-new-yorker.json', {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
        },
      });
      const jsonData = await response.json();
      let arrOfLinks = jsonData.ctx_obj[0].ctx_obj_targets[0].target;
      setElements(arrOfLinks.slice(0, -1));
    } catch (err) {
      console.error(err.message);
    }
  };

  useEffect(() => {
    fetchElements();
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
          {elements &&
            elements.map((element, idx) => (
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
                  <a href={element.target_url} target="_blank" rel="noopener noreferrer">
                    {element.target_url.substring(0, 35)}...
                  </a>
                </div>
              </div>
            ))}
        </div>
      </div>
    </>
  );
};

export default List;
