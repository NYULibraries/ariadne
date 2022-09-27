import { useState } from 'react';
import getLinks from '../api/getLinks';

const useApi = () => {
  const [elements, setElements] = useState([]);
  const [error, setError] = useState(null);
  const [lastElement, setLastElement] = useState(null);

  const request = async () => {
    try {
      const response = await getLinks();
      const jsonData = await response.data;
      let arrOfLinks = jsonData.ctx_obj[0].ctx_obj_targets[0].target;
      setElements(arrOfLinks.slice(0, -1));
      setLastElement(arrOfLinks.at(-1));
    } catch (error) {
      setError(error.message || 'Something went wrong');
    }
  };

  return { elements, request, lastElement, error };
};

export default useApi;
