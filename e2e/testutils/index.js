function removeSourceMappingUrlComments(html) {
  const regex = new RegExp( '/\\*#\\ssourceMappingURL=\\s*\\S+\\s\\*\\/', 'g' );

  return html.replace( regex, '/* E2E TEST EDIT: sourceMappingURL comments elided */' );
}

// NOTE: it's current not possible to use a custom flag like `--update-golden-files`
// with `playwright`:
// "[Feature] Add support for test.each / describe.each #7036"
// https://github.com/microsoft/playwright/issues/7036
function updateGoldenFiles() {
  return process.env.UPDATE_GOLDEN_FILES &&
         process.env.UPDATE_GOLDEN_FILES.toLowerCase() !== 'false';
}

export {
  removeSourceMappingUrlComments,
  updateGoldenFiles,
};
