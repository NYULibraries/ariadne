import apiClient from './apiClient';

let escapedSearch = window.location.search.replace(';', '%3B');
const fetchData = () => apiClient.get(escapedSearch);

export default { fetchData };
