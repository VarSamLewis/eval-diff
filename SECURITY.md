# Security Policy

## Supported versions

This project is pre-1.0 and should be treated as experimental. Only the latest released action version receives security fixes.

## Reporting a vulnerability

Do not open a public issue for a suspected vulnerability. Email the repository owner with a concise description, affected version or commit, reproduction steps, and potential impact. Do not include API keys, customer data, or proprietary code in the report.

The maintainer will acknowledge reports within five business days and work with the reporter on validation, remediation, and coordinated disclosure.

## Security boundaries

Unless dry-run mode is enabled, this action transmits the selected code diff and an API bearer key to the TypeSafe Jev API. Review [the data-handling guide](docs/security_and_data_handling.md) before use. The action pins its runtime image by SHA-256 digest and runs it without root privileges, but consumers remain responsible for approving the third-party service and protecting their workflow secrets.
