# Security Policy

## Reporting a Vulnerability

Please report security vulnerabilities privately via GitHub Security Advisories
for this repository. Do not open public issues for undisclosed vulnerabilities.

## Supported Versions

Only the latest release on `main` is supported.

## Security Practices

- Dependencies are monitored by Dependabot.
- CI runs `gosec` and `govulncheck`.
- The CLI limits input size to 10 MiB.
