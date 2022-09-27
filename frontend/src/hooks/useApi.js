import { useState } from 'react';
import getLinks from '../api/getLinks';

const useApi = () => {
  const [elements, setElements] = useState([]);
  const [error, setError] = useState(null);
  const [lastElement, setLastElement] = useState(null);
  const [loading, setLoading] = useState(false);

  const request = async () => {
    setLoading(true);
    try {
      const response = await getLinks();
      const jsonData = await response.data;
      let arrOfLinks = jsonData.ctx_obj[0].ctx_obj_targets[0].target;
      setElements(arrOfLinks.slice(0, -1));
      setLastElement(arrOfLinks.at(-1));
    } catch (error) {
      setError(error.message || 'Something went wrong');
    } finally {
      setLoading(false);
    }
  };

  return { elements, request, lastElement, error, loading };
};

export default useApi;
