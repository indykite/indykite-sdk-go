module.exports = {
  extends: ["@commitlint/config-conventional"],
  defaultIgnores: false,
  rules: {
    //   0 - Disabled, 1 - Warning, 2 - Error
    "body-max-line-length": [2, "always", 72],
    "header-max-length": [2, "always", 72],
    "subject-max-length": [2, "always", 50],
    "type-enum": [
      2,
      "always",
      ["build", "chore", "ci", "docs", "feat", "fix", "perf", "refactor", "revert", "style", "test"],
    ],
    // Is warning only. But with action (commitling.yaml), warnings are converted to errors on non-required check.
    // So we can fail if commit is missing JIRA ticket, but it will be still possible to merge.
    "contains-jira-ticket": [1, "always"],
    "scope-enum": [
      2,
      "always",
      [
        "logging",
        "services",
        "docs",
        "dependencies",
        "deps",
        "auth",
        "api",
        "pkg",
        "proto",
        "test",
        "master",
        "examples",
      ],
    ],
  },

  plugins: [
    {
      // Based on documentation, we can have only 1 local plugin.
      // But it can implement multiple rules if needed.
      // https://commitlint.js.org/#/reference-plugins?id=local-plugins
      rules: {
        "contains-jira-ticket": ({ header, body, footer }) => {
          let regex = /\bENG\-[0-9]{2,}\b/;
          if (regex.test(header)) {
            return [false, "Please, move your JIRA ticket into commit body, and do not include in header directly."];
          }

          // There must be ENG-XX where XX is at least 2 digit number
          if (regex.test(body) === false && regex.test(footer) === false) {
            return [
              false,
              "Your commit message is missing JIRA ticket. Consider adding it to commit body or ignore this error.",
            ];
          }
          return [true];
        },
      },
    },
  ],
};
