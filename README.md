# Eval Diff

Evaluate a Git diff for change impact and potential breaking changes using the TypeSafe Jev API.

This action compares `origin/main...HEAD`, sends the resulting diff to the evaluation service, and returns an impact score, confidence score, and breaking-change flag.

## Note

This codebase is still pre v1.0.0 and should not be used for anything beyond testing and POC. 

It has a few surface layers left to harden and some performance improvements before use on large code bases.  

## Usage

Add the action after checking out the repository with its history available:

```yaml
name: Evaluate diff

on:
  pull_request:

jobs:
  evaluate:
    runs-on: ubuntu-latest
    permissions:
      contents: read
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Evaluate Git Diff
        id: eval_diff
        uses: varsamlewis/eval-diff@v1
        with:
          typesafe_api_key: ${{ secrets.TYPESAFE_API_KEY }}

      - name: Show evaluation
        run: |
          echo "Impact: ${{ steps.eval_diff.outputs.impact_score }}"
          echo "Confidence: ${{ steps.eval_diff.outputs.confidence }}"
          echo "Breaking change: ${{ steps.eval_diff.outputs.is_breaking }}"
```

The checkout must fetch `origin/main`, because the action calculates `git diff origin/main...HEAD`. The action runs a Docker container, so the runner also needs Docker and outbound network access to GitHub Container Registry and TypeSafe's API.

## Supported action reference

Use the GitHub Action release tag in workflows—not a container tag and not the TypeSafe HTTP endpoint directly:

```yaml
uses: varsamlewis/eval-diff@v1
```

`v1` is the supported rolling major-version tag and receives compatible updates after release validation. To keep a workflow fixed to one release, use an exact action tag instead:

```yaml
uses: varsamlewis/eval-diff@v0.0.2
```

Do not use `@latest`; it is not a supported GitHub Action release channel. The GHCR container version is an internal implementation detail selected by the action's checked-in `CONTAINER_VERSION` file. Consumers should not invoke the TypeSafe API or the container image directly unless they have separately integrated and approved those interfaces.

## Inputs

| Input | Required | Default | Description |
| --- | --- | --- | --- |
| `typesafe_api_key` | Yes | — | API key for the TypeSafe evaluation service. Store it as a GitHub Actions secret. |
| `format` | No | `standard` | CLI output format: `standard`, `score`, or `full`. |
| `max_chars` | No | `4000` | Maximum number of diff characters evaluated. Content beyond the limit is omitted. |
| `dry_run` | No | `false` | When `true`, reports payload statistics without calling the external API. |

## Outputs

| Output | Description |
| --- | --- |
| `impact_score` | Model-generated impact score on a 1–5 scale. |
| `confidence` | Model-generated confidence score. |
| `is_breaking` | `true` when the model's breaking-change probability is greater than `0.8`; otherwise `false`. |

## Security and data handling

Unless `dry_run` is enabled, the action sends the evaluated Git diff to the third-party TypeSafe Jev API. Diffs can include proprietary code, security fixes, internal paths, or accidentally committed secrets. Run this action only when you are authorized to disclose that material to TypeSafe and its service providers.

`--max-chars` reduces the amount sent, but it does not redact sensitive data and may omit relevant changes. The result is an advisory signal—not a substitute for code review, testing, security review, or change-control processes.

See [security, privacy, and operational considerations](docs/security_and_data_handling.md) before enabling the action, especially for private repositories or regulated codebases.

## Limits and failure behavior

The action uses only the first `max_chars` characters of the diff. The default is 4,000 characters. The API client uses a 15-second timeout. Missing history, Docker failures, registry/network failures, API errors, or an invalid response cause the evaluation step to fail.

For a no-transmission check, set `dry_run: true`:

```yaml
- uses: varsamlewis/eval-diff@v1
  with:
    typesafe_api_key: ${{ secrets.TYPESAFE_API_KEY }}
    dry_run: true
```

## Documentation

- [Security, privacy, and data handling](docs/security_and_data_handling.md)
- [CLI guide](docs/cli_guide.md)
- [Contributing guide](docs/contributing_guide.md)
