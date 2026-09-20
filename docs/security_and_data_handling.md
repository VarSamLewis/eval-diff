# Security, Privacy, and Data Handling

## What the evaluator sends

With normal evaluation enabled, the CLI packages the supplied content in the `git_diff` field of a JSON request and sends it over HTTPS to the third-party TypeSafe Jev API at `https://api.typesafe.ai/v1/systemone`.

For the GitHub Action, that content is `git diff origin/main...HEAD`. For the CLI, it is either standard input or the full contents of every file path passed as an argument. The request also includes the `TYPESAFE_API_KEY` as a bearer token.

The action runs a container image from GitHub Container Registry: `ghcr.io/varsamlewis/eval-diff`. It therefore depends on both that registry and the TypeSafe API being reachable from the runner.

## Risks to assess before use

- **Third-party code disclosure:** a diff can contain proprietary code, unreleased features, vulnerability fixes, internal file paths, customer-related data, or accidentally included credentials. The evaluator is unsuitable for content that you are not authorized to share with TypeSafe and its service providers.
- **Secret exposure:** truncation does not remove secrets. A secret near the beginning of the selected content can still be transmitted. Scan or redact input before evaluation if there is any risk of sensitive material.
- **API-key handling:** the API key is a bearer credential. Keep it in GitHub Actions secrets or an approved secret manager. Do not commit it, pass it through a command line likely to be logged, or keep it in a shared `.env` file.
- **Incomplete analysis:** the GitHub Action evaluates at most 4,000 characters by default; the CLI default is 100,000. Content after the configured limit is omitted, so a low-risk result does not mean the complete change is safe.
- **Non-deterministic advice:** impact, confidence, and breaking-change values are model-generated signals. They can be inaccurate, incomplete, or change as the external service changes. Continue to use code review, tests, security review, deployment controls, and rollback plans.
- **Availability and supply chain:** evaluation requires Docker, outbound network access, the configured container image, and the external API. Registry outages, network restrictions, API errors, or service changes can fail the workflow. Pin the action to a trusted version or commit SHA where your supply-chain policy requires it.

## Safer operating practices

1. Confirm that your organization has approved TypeSafe for the class of code or data in the repository.
2. Use a narrowly scoped API key and store it in a secret manager.
3. Inspect or scan the diff for credentials and sensitive data before evaluation.
4. Set a suitable `max_chars` limit, knowing it reduces volume rather than sensitivity.
5. Use `dry_run: true` in GitHub Actions or `--dry-run=true` in the CLI to validate payload size without sending the content. The CLI does not require an API key in dry-run mode; the GitHub Action still declares `typesafe_api_key` as a required input, even though dry-run does not use it for an API call.
6. Treat failed evaluations according to your delivery policy; do not rely on this action as the only quality or security gate.

## Retention and compliance

This repository does not define TypeSafe's retention, training, subprocessors, regional processing, incident handling, or contractual terms. Obtain those details directly from TypeSafe and complete your organization's privacy, legal, procurement, and security reviews before sending sensitive or regulated content.
