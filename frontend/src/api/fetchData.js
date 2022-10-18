import apiClient from './apiClient';

const fetchData = () => apiClient.get('/' + window.location.search);

export default { fetchData };
