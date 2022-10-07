import { useState } from 'react';
import fetchData from '../api/fetchData';
import { getLinks } from '../aux/helpers';

const useApi = () => {
  const [resource, setResource] = useState([]);
  const [error, setError] = useState(null);
  const [resourceLastElement, setResourceLastElement] = useState(null);
  const [loading, setLoading] = useState(false);

  const fetchResource = async () => {
    setLoading(true);
    try {
      const response = await fetchData();
      const jsonData = await response.data;
      // TODO: add a getLinks helper method to retrieve the links
      let arrOfLinks = getLinks(jsonData);
      setResource(arrOfLinks.slice(0, -1));
      setResourceLastElement(arrOfLinks.at(-1));
    } catch (error) {
      setError('Something went wrong');
    } finally {
      setLoading(false);
    }
  };

  return { resource, fetchResource, resourceLastElement, error, loading };
};

export default useApi;
