
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
  const urlParams = new URLSearchParams(queryString.toLowerCase());
  let parameter = urlParams.get(parameterName);

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

