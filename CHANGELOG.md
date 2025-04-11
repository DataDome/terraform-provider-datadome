# DataDome Terraform Provider

## 2.2.1 (2025-04-11)

### BUG FIXES:

- Fix the `Update` method of the `ClientCustomRule`

## 2.2.0 (2025-04-02)

### BUG FIXES:

- Update dependencies to fix vulnerabilities
- Update `Endpoint` structure to fix update issue when fields were omitted from the payload

### ENHANCEMENTS:

- Add [`query`](https://docs.datadome.co/docs/endpoints#1-traffic-query) field for `endpoint` management
- Upgrade `go` version to `1.23`

## 2.1.0 (2025-02-03)

### BUG FIXES:

- Update dependencies to fix vulnerabilities

### ENHANCEMENTS:

- Add support for [endpoints management](https://docs.datadome.co/docs/endpoints)

## 2.0.2 (2025-01-15)

### BUG FIXES:

- Update list of allowed `endpoint_type` for [adding](https://docs.datadome.co/reference/post_1-1-protection-custom-rules), [updating](https://docs.datadome.co/reference/put_1-1-protection-custom-rules-customruleid), and importing custom rule resources

## 2.0.1 (2024-12-18)

### BUG FIXES:

- Update dependencies to fix vulnerabilities

## 2.0.0 (2024-07-31)

### BREAKING CHANGES:

- Change `whitelist` response to `allow`

### ENHANCEMENTS:

- Upgrade `go` version and dependencies
- Improve error handling
- Improve `godoc` for each functions and constants
- Update CI to include linting and static code checking
- Add unit and acceptance tests

## 1.1.0 (2023-12-15)

### BUG FIXES:

- Fix the authentication to DataDome API from query param to header

### ENHANCEMENTS:

- Upgrade `go` version and dependencies
