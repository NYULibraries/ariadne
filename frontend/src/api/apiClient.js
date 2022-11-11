const apiClient = {
  get: async function get(path) {
    const url = `${process.env.REACT_APP_API_URL}${path}`;
    const response = await fetch(url);

    return response.json();
  },
};

export default apiClient;
