import {
  DEFAULT_IMG_CLASS,
  DEFAULT_LINK,
  DEFAULT_LOGO,
  NYUAD,
  NYUAD_IMG_CLASS,
  NYUAD_LINK,
  NYUAD_LOGO,
  NYUSH,
  NYUSH_IMG_CLASS,
  NYUSH_LINK,
  NYUSH_LOGO
} from './institutionConstants';

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

// Helper function for Banner.js component

const getInstitution = (institution) => {
  let logo = DEFAULT_LOGO;
  let link = DEFAULT_LINK;
  let imgClass = DEFAULT_IMG_CLASS;

  if (institution.toUpperCase() === NYUAD) {
    logo = NYUAD_LOGO;
    link = NYUAD_LINK;
    imgClass = NYUAD_IMG_CLASS;
  } else if (institution.toUpperCase() === NYUSH) {
    logo = NYUSH_LOGO;
    link = NYUSH_LINK;
    imgClass = NYUSH_IMG_CLASS;
  }

  return { logo, link, imgClass };
};

export {
  getCoverageStatement,
  getInstitution,
  getLinks,
  getParameterFromQueryString,
  getInstitutionQueryParameter
};

