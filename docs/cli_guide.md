# CLI Tool Documentation

The evaluation engine is packaged as a Go command-line tool. It reads a git diff from standard input (stdin) and evaluates the content.

## Data Handling and Security

Unless `--dry-run=true` is used, the CLI sends the supplied standard input or full contents of specified files to the third-party TypeSafe Jev API. Review the [security, privacy, and data-handling guide](security_and_data_handling.md) before sending proprietary or sensitive material. Dry-run mode reports payload statistics without transmitting content or requiring an API key.

Installation

Compile the binary using make:

```
make build
```

The compiled executable is placed at bin/eval-diff.

Environment Variables

```
TYPESAFE_API_KEY (required unless --dry-run=true): The API key used to authenticate against the TypeSafe Jev API.

GITHUB_OUTPUT (optional): File path used when running inside GitHub Actions to set step output parameters.

```

Command-Line Flags

```
--format: Set output detail. Accepted values are standard, score, or full. Default is standard.

--max-chars: Maximum number of characters to process from the diff payload before truncation. Default is 100000.

--dry-run: When set to true, outputs payload statistics and skips the network request to the API. Default is false.

```

Examples

Run Dry Run on Current Diff

```
git diff main...HEAD | ./bin/eval-diff --dry-run=true

```

Evaluate with Full Output Format

```
export TYPESAFE_API_KEY="your-api-key-here"
git diff main...HEAD | ./bin/eval-diff --format=full

```

Truncate Large Diffs

```
export TYPESAFE_API_KEY="your-api-key-here"
git diff main...HEAD | ./bin/eval-diff --max-chars=50000
```
