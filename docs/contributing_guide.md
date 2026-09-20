# Contributing and Development Guide

This document describes how to make changes, build, test, and release the action.

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

Ensure all changes are merged into the main branch.

Update the version tag locally:

```
git switch main
git pull origin main
git push origin v1 --force

```

The release workflow in .github/workflows/release.yml triggers on the tag push, compiles the binary inside Docker, and pushes the image tag ghcr.io/varsamlewis/eval-diff:v1 to GitHub Container Registry.
