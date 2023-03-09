// Source: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/freeze
function deepFreeze(object) {
  // Retrieve the property names defined on object
  const propNames = Reflect.ownKeys(object);

  // Freeze properties before freezing self
  for (const name of propNames) {
    const value = object[name];

    if ((value && typeof value === "object") || typeof value === "function") {
      deepFreeze(value);
    }
  }

  return Object.freeze(object);
}

// Helper function for List.js component
function getCoverageStatement(link) {
  const coverage = link?.coverage?.[0];
  return coverage ? coverage.coverage_text[0].threshold_text[0].coverage_statement?.join('. ') : '';
}


function getInstitutionQueryParameter() {
  const queryString = window.location.search;
  return getParameterFromQueryString(queryString, 'institution');
}

//  Helper functions for useApi hook
function getLinks(jsonData) {
  return jsonData.ctx_obj[0].ctx_obj_targets[0].target;
}

// Helper functions for Banner.js component
function getParameterFromQueryString(queryString, parameterName) {
  const urlParams = new URLSearchParams(queryString.toLowerCase());
  return urlParams.get(parameterName.toLowerCase());
}

export {
  deepFreeze,
  getCoverageStatement,
  getInstitutionQueryParameter,
  getLinks,
  getParameterFromQueryString,
};

