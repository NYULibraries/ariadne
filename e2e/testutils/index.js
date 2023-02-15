function removeSourceMappingUrlComments(html) {
  const regex = new RegExp( '/\\*#\\ssourceMappingURL=\\s*\\S+\\s\\*\\/', 'g' );

  return html.replace( regex, '/* E2E TEST EDIT: sourceMappingURL comments elided */' );
}

function updateGoldenFiles() {
  return process.env.UPDATE_GOLDEN_FILES &&
         process.env.UPDATE_GOLDEN_FILES.toLowerCase() !== 'false';
}

export {
  removeSourceMappingUrlComments,
  updateGoldenFiles,
};
