import { useState } from 'react';

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
          // Currently the backend only returns single-record responses.
          // That may change in the future.
          const arrOfLinks = responseBody.records[0].links;
          arrOfLinks.sort((a, b) => a.display_name.localeCompare(b.display_name));
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
