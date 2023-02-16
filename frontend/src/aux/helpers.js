
//  Helper functions for useApi hook
const getLinks = (jsonData) => {
  return jsonData.ctx_obj[0].ctx_obj_targets[0].target;
};

// Helper function for List.js component
const getCoverageStatement = (link) => {
  return link.coverage[0].coverage_text[0].threshold_text[0].coverage_statement?.join('. ');
};

// Helper functions for Banner.js component
const getParameterFromQueryString = (queryString, parameterName) => {
  const urlParams = new URLSearchParams(queryString);
  let parameter = urlParams.get(parameterName.toLowerCase());
  if (urlParams.has(`umlaut.${parameterName.toLowerCase()}`)) {
    parameter = urlParams.get(`umlaut.${parameterName.toLowerCase()}`);
  }
  return parameter;
};


const getInstitutionQueryParameter = (parameterName) => {
  const queryString = window.location.search;
  return getParameterFromQueryString(queryString, parameterName);
};

export {
  getCoverageStatement,
  getLinks,
  getParameterFromQueryString,
  getInstitutionQueryParameter
};

