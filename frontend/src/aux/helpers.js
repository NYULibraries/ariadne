//  Helper functions for useApi hook
function getLinks(jsonData) {
  return jsonData.ctx_obj[0].ctx_obj_targets[0].target;
}

// Helper function for List.js component
function getCoverageStatement(link) {
  return link.coverage[0].coverage_text[0].threshold_text[0].coverage_statement?.join('. ');
}

// Helper functions for Banner.js component
function getParameterFromQueryString(queryString, parameterName) {
  const urlParams = new URLSearchParams(queryString.toLowerCase());
  return urlParams.get(parameterName);
}

function getInstitutionQueryParameter(parameterName) {
  const queryString = window.location.search;
  return getParameterFromQueryString(queryString, parameterName);
}

export {
  getCoverageStatement,
  getLinks,
  getParameterFromQueryString,
  getInstitutionQueryParameter
};

