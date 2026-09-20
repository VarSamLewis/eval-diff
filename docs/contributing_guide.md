# Contributing and Development Guide

This document describes how to make changes, build, test, and release the action.

## Safe Testing

Use `--dry-run=true` for local development unless you have approval to transmit the test input to the third-party TypeSafe Jev API. A normal evaluation sends the supplied diff or files and an API bearer token to that service. Do not use production code, credentials, customer data, or unreviewed security fixes as test input. See the [README security considerations](../README.md#security-privacy-and-risk-considerations) for the full risk summary.

Prerequisites

- Go 1.22 or higher
- Git
- Docker
- Make

Local Setup

Clone the repository and set up the repository git hooks path:

```
git config core.hooksPath .githooks
```

Ensure the pre-commit script is executable:

```
chmod +x .githooks/pre-commit
```

Build and Test Commands

A Makefile is provided to run common operations using standard Go tools.

Run code formatting, static checks, tests, and build the binary:

```
make

```

Run tests only:

```
make test

```

Format source code:

```
make fmt
```

Run static analysis:

```
make vet

```

Build the static Linux binary:

```
make build

```

Clean build artifacts:

```
make clean

```

Testing the Action Locally

To test the container locally, build the Docker image using the provided Dockerfile:

```
docker build -t eval-diff:local .

```

Generate a git diff and pipe it to the container:

```
git diff main...HEAD | docker run --rm -i \
  -e TYPESAFE_API_KEY="your-api-key" \
  eval-diff:local --dry-run=true

```

Release Process

The action uses a pre-built Docker image hosted on GitHub Container Registry.

Git action releases and container releases are versioned separately. GitHub release tags (for example, `v0.0.3`) determine the version users reference with `uses: varsamlewis/eval-diff@...`. The `CONTAINER_VERSION` file determines the container image tag the action pulls.

### Publish a container image

Update `CONTAINER_VERSION`, merge the change to `main`, then manually run `.github/workflows/release.yml`. It publishes only this version-specific image tag:

```
ghcr.io/varsamlewis/eval-diff:v<CONTAINER_VERSION>

```

The workflow does not create, update, or force-push Git tags.

After publishing, retrieve the image's SHA-256 digest and record it in `CONTAINER_DIGEST`. Action releases pull the digest, not the mutable version tag. Publish the action release only after the digest change is merged and tested.

### Publish an action release

Create a GitHub release with a new immutable action tag, such as `v0.0.3`, and select the Marketplace publishing option when appropriate. Update the rolling `v1` tag only after verifying that release. The action source at a release tag includes its `CONTAINER_VERSION` file, which makes its selected container image explicit.
