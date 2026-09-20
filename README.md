System One Model Diff Evaluator

The System One Model Diff Evaluator is a GitHub Action and CLI tool that analyzes git diffs to assess impact severity, risk, and breaking changes.

Features

Analyzes code diffs against a target git branch.

Calls the TypeSafe Jev API to compute an impact score and confidence level.

Identifies potential breaking changes automatically.

Runs as a pre-compiled container image on GitHub Container Registry for fast execution.

Includes a standalone CLI binary for local development and workflow integration.

Basic Usage

```
To use this action in a GitHub Actions workflow, add the following step:

- name: Checkout code
  uses: actions/checkout@v4
  with:
    fetch-depth: 0

- name: Evaluate Git Diff
  uses: varsamlewis/eval-diff@v1
  with:
    typesafe_api_key: ${{ secrets.TYPESAFE_API_KEY }}
```

Documentation

Development, building, testing, and releasing: see docs/CONTRIBUTING.md

Local CLI usage and flags: see docs/CLI.md
