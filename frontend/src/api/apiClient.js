const apiClient = {
  get: async function get(path) {
    const url = `${process.env.REACT_APP_API_URL}${path}`;

    return await fetch(url);
  },
};

export default apiClient;
