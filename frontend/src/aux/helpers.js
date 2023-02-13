import { institutions } from './institutionConstants';

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
    urlParams.delete(`umlaut.${parameterName.toLowerCase()}`);
    urlParams.set(parameterName.toLowerCase(), parameter);
    history.pushState(null, '', `?${urlParams.toString()}`)
  }
  return parameter;
};

const getInstitutionQueryParameter = (parameterName) => {
  const queryString = window.location.search;
  return getParameterFromQueryString(queryString, parameterName);
};

const getInstitution = (institution) => {
  const { logo, link, imgClass } = institutions[institution?.toLowerCase()] || institutions.nyu;
  return { logo, link, imgClass };
};


export {
  getCoverageStatement,
  getInstitution,
  getLinks,
  getParameterFromQueryString,
  getInstitutionQueryParameter
};

