import apiClient from './apiClient';

var query = window.location.search;

const fetchData = () => apiClient.get('/' + query);

export default { fetchData };
