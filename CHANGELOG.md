# DataDome Terraform Provider

## 2.4.0 (2026-06-30)

- Add support of `rate_limit` and `time_box` policy options for custom rules
- Add new `endpoint_type` values for custom rules:
  - `agentic-general`
  - `agentic-account-creation`
  - `agentic-login`
  - `agentic-cart`
  - `agentic-forms`
  - `agentic-payment`
- Add new `response` values for custom rules:
  - device_check
  - intent_based
  - monetize
- Add `overridden_bot` field to specify which verified model is overridden by the custom rule
- Add `Agentic Protocol` source value for endpoints

## 2.3.2 (2026-03-25)

- Update dependencies to fix vulnerabilities

## 2.3.1 (2025-10-27)

- Fix issue when updating `endpoint` resource with optional fields set to empty strings

## 2.3.0 (2025-09-04)

- Update dependencies to fix vulnerabilities
- Add support of activation and expiration dates for custom rules

## 2.2.2 (2025-05-19)

- Update dependencies to fix vulnerabilities

## 2.2.1 (2025-04-11)

- Fix the `Update` method of the `ClientCustomRule` structure

## 2.2.0 (2025-04-02)

- Update dependencies to fix vulnerabilities
- Update `Endpoint` structure to fix update issue when fields were omitted from the payload
- Add [`query`](https://docs.datadome.co/docs/endpoints#1-traffic-query) field for `endpoint` management
- Upgrade `go` version to `1.23`

## 2.1.0 (2025-02-03)

- Update dependencies to fix vulnerabilities
- Add support for [endpoints management](https://docs.datadome.co/docs/endpoints)

## 2.0.2 (2025-01-15)

- Update list of allowed `endpoint_type` for [adding](https://docs.datadome.co/reference/post_1-1-protection-custom-rules), [updating](https://docs.datadome.co/reference/put_1-1-protection-custom-rules-customruleid), and importing custom rule resources

## 2.0.1 (2024-12-18)

- Update dependencies to fix vulnerabilities

## 2.0.0 (2024-07-31)

### Breaking changes

- Change `whitelist` response to `allow`

### Other changes

- Upgrade `go` version and dependencies
- Improve error handling
- Improve `godoc` for each functions and constants
- Update CI to include linting and static code checking
- Add unit and acceptance tests

## 1.1.0 (2023-12-15)

- Fix the authentication to DataDome API from query param to header
- Upgrade `go` version and dependencies
