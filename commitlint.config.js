module.exports = {
  extends: ["@commitlint/config-conventional"],
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
    "scope-enum": [
      2,
      "always",
      ["logging", "services", "docs", "dependencies", "deps", "auth", "api", "pkg", "proto", "test", "master", "examples"],
    ],
  },
};
