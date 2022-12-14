import { useState } from 'react';
import { getLinks } from '../aux/helpers';

export default (apiFunc) => {
  const [resource, setResource] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const fetchResource = async (...args) => {
    setLoading(true);
    try {
      const response = await apiFunc(...args);
      if (response.ok) {
        const responseBody = await response.json();
        if (responseBody.errors.length === 0) {
          const arrOfLinks = getLinks(responseBody.records);
          setResource(arrOfLinks);
        } else {
          setError(`The backend API returned errors: ${responseBody.errors.map((error) => `"${error}"`).join(', ')}`);
        }
      } else {
        setError(`The backend API returned an HTTP error response: ${response.status} (${response.statusText})`);
      }
    } catch (error) {
      setError(`Error fetching data from the Ariadne API: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  return { resource, fetchResource, error, loading };
};
