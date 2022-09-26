import 'bootstrap/dist/css/bootstrap.min.css';
import { useEffect, useState } from 'react';
import axios from 'axios';

const List = () => {
  const [elements, setElements] = useState([]);

  const baseURl = process.env.REACT_APP_API_URL;
  const query =
    '?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat';

  const getElements = async () => {
    try {
      const response = await axios.get(`${baseURl}` + query);
      const jsonData = response.data;
      let arrOfLinks = jsonData.ctx_obj[0].ctx_obj_targets[0].target;
      setElements(arrOfLinks.slice(0, -1));
    } catch (err) {
      console.error(err.message);
    }
  };

  useEffect(() => {
    getElements();
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
                  <small>{element.coverage[0].coverage_text[0].threshold_text[0].coverage_statement}</small>
                </div>
              </div>
            ))}
        </div>
      </div>
    </>
  );
};

export default List;
