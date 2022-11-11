import { useState } from 'react';
import { getLinks } from '../aux/helpers';

export default (apiFunc) => {
  const [resource, setResource] = useState(null);
  const [error, setError] = useState(null);
  const [resourceLastElement, setResourceLastElement] = useState(null);
  const [loading, setLoading] = useState(false);

  const fetchResource = async (...args) => {
    setLoading(true);
    try {
      const response = await apiFunc(...args);
      const arrOfLinks = getLinks(response.records);
      setResource(arrOfLinks.slice(0, -1));
      setResourceLastElement(arrOfLinks.at(-1));
    } catch (error) {
      setError('Something went wrong');
    } finally {
      setLoading(false);
    }
  };

  return { resource, fetchResource, error, resourceLastElement, loading };
};
