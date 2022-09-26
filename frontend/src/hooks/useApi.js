import { useState } from 'react';

const useApi = (url, query) => {
  const [elements, setElements] = useState([]);
  const [error, setError] = useState(null);
  const [lastElement, setLastElement] = useState(null);

  const getElements = async () => {
    try {
      const response = await fetch(url + query);
      const jsonData = await response.json();
      let arrOfLinks = jsonData.ctx_obj[0].ctx_obj_targets[0].target;
      setElements(arrOfLinks.slice(0, -1));
      setLastElement(arrOfLinks.at(-1));
    } catch (error) {
      setError(error.message || 'Something went wrong');
    }
  };

  return { elements, getElements, lastElement, error };
};

export default useApi;
