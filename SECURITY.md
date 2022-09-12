# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| v1.x    | :white_check_mark: |

## Reporting a Vulnerability

Please submit any security vulnerabilities via [our public bug bounty program](https://app.intigriti.com/programs/intigriti/intigriti) so we can reward you accordingly.

## Security features

1. We do *CodeQL code scanning* [in our CI/CD pipeline](https://github.com/intigriti/sdk-go/blob/main/.github/workflows/securityscan.yml#L13).
2. We execute *Go vulnerability scanning* [in our CI/CD pipeline](https://github.com/intigriti/sdk-go/blob/main/.github/workflows/securityscan.yml#L33).
3. We run *trivy vulnerability scanning* [in our CI/CD pipeline](https://github.com/intigriti/sdk-go/blob/main/.github/workflows/securityscan.yml#57).
4. Our dependencies are automatically kept up-to-date via [GitHub Dependabot](https://github.com/intigriti/sdk-go/blob/main/.github/dependabot.yml).
5. Pull requests are required and need approval from [the code owners](https://github.com/intigriti/sdk-go/blob/main/.github/CODEOWNERS).
6. We utilize a release process [that produces SBOMs and is SLSA 3 compliance](https://github.com/intigriti/sdk-go/blob/main/.github/goreleaser.yml).
7. Tests need to complete before a merge [in our CI/CD pipeline](https://github.com/intigriti/sdk-go/blob/main/.github/workflows/test.yml).
