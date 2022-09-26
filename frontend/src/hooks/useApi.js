import { useState } from 'react';

const useApi = (url, query) => {
  const [elements, setElements] = useState([]);

  const getElements = async () => {
    try {
      const response = await fetch(url + query);
      const jsonData = await response.json();
      let arrOfLinks = jsonData.ctx_obj[0].ctx_obj_targets[0].target;
      setElements(arrOfLinks.slice(0, -1));
    } catch (err) {
      console.error(err.message);
    }
  };

  return { elements, getElements };
};

export default useApi;
