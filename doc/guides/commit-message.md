# Commit message guidelines

[https://www.conventionalcommits.org/en/v1.0.0/#summary](https://www.conventionalcommits.org/en/v1.0.0/#summary)

The Conventional Commits specification is a lightweight convention on top of commit messages. It provides an easy set of rules for creating an explicit commit history; which makes it easier to write automated tools on top of.
This convention dovetails with SemVer, by describing the features, fixes, and breaking changes made in commit messages.

The commit message should be structured as follows:

```md
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

The commit contains the following structural elements, to communicate intent to the consumers of your library:

- **fix**: a commit of the _type_ `fix` patches a bug in your codebase (this correlates with **PATCH** in Semantic Versioning).
- **feat**: a commit of the _type_ `feat` introduces a new feature to the codebase (this correlates with **MINOR** in Semantic Versioning).
- **BREAKING CHANGE**: a commit that has a footer `BREAKING CHANGE:`, or appends a `!` after the type/scope, introduces a breaking API change (correlating with **MAJOR** in Semantic Versioning). A BREAKING CHANGE can be part of commits of any type.
- types other than fix: and feat: are allowed, for example @commitlint/config-conventional (based on the the Angular convention) recommends build:, chore:, ci:, docs:, style:, refactor:, perf:, test:, and others.
- footers other than BREAKING CHANGE: `description` may be provided and follow a convention similar to git trailer format.
  Additional types are not mandated by the Conventional Commits specification, and have no implicit effect in Semantic Versioning (unless they include a BREAKING CHANGE).
  A scope may be provided to a commit’s type, to provide additional contextual information and is contained within parenthesis, e.g., `feat(parser): add ability to parse arrays.`

Allowed types:

```json
["build", "chore", "ci", "docs", "feat", "fix", "perf", "refactor", "revert", "style", "test"]
```

Allowed scopes:

```json
["logging", "services", "docs", "dependencies", "auth", "api", "pkg", "proto", "cypher", "sdk", "schema", "test", "master", "examples"]
```

## Examples

## Commit message with description, scope and breaking change footer

```md
feat(component): add new indykite feature

BREAKING CHANGE: replaces configuration file syntax
```

## Commit message with no body

```md
feat: add new indykite feature
```

## Commit message with scope

```md
feat(component): add new indykite feature
```

If you are new to contributing to IndyKite, please try to do your best at
conforming to these guidelines, but do not worry if you get something wrong.
One of the existing contributors will help get things situated and the
contributor landing the Pull Request will ensure that everything follows
the project guidelines.
