# Security Policy

## Reporting a vulnerability

Please **do not** open a public issue for security vulnerabilities.

Report privately to **livio@action.tec.br** with:

- a description of the issue
- steps to reproduce
- the affected version

We acknowledge reports within **3 business days**.

## Supported versions

Only the latest published version receives security fixes.

## Scope

This repository contains the Go SDK for the FalaAI API. Vulnerabilities in the hosted
API itself should also be reported to the address above.

## Handling of secrets

Never commit API keys (`fai_...`). Pass them at runtime through an environment variable
or a secret manager.
