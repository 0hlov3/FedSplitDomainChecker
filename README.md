# FedSplitDomainChecker Documentation

**FedSplitDomainChecker** is a command-line tool designed to validate and manage Fediverse split-domain deployments. It helps ensure proper configuration of `host-meta`, `nodeinfo`, and `webfinger` endpoints for federated services such as GoToSocial.

## Features

- Validates `host-meta`, `nodeinfo`, and `webfinger` endpoints.
- Supports split-domain deployments where the account domain differs from the host domain.
- Provides detailed logging for each validation step.
- Offers verbose logging for debugging.

---

## Installation

### Prerequisites

- **Go**: Ensure you have Go version 1.20 or higher installed.

### Clone the Repository

```bash
git clone https://github.com/0hlov3/FedSplitDomainChecker.git
cd FedSplitDomainChecker
```

### Build the Binary

```bash
go build -o FedSplitDomainChecker
```

---

## Usage

### Command Syntax

```bash
./FedSplitDomainChecker [command] [flags]
```

### Commands

#### `checkSplitDomain`

Validates the `host-meta`, `nodeinfo`, and `webfinger` endpoints for a split-domain deployment.

**Example**:

```bash
./FedSplitDomainChecker checkSplitDomain --accountDomain example.org --hostDomain social.example.org --account admin@example.org
```

**Flags**:

| Flag            | Description                                   | Default Value                  |
|-----------------|-----------------------------------------------|--------------------------------|
| `--accountDomain` | The account domain to check.                  | `https://gotosocial.org`       |
| `--hostDomain`  | The host domain where the instance is hosted. | `https://gts.gotosocial.org`   |
| `--account`     | The user account to validate.                 | `admin@gotosocial.org`         |
| `--verbose`, `-v` | Enable verbose logging for debugging.         | `false`                        |
| `--debug`, `-d` | Enable debug mode logging.                    | `false`                        |

### Example Workflow

#### Input

- **Account Domain**: `example.org`
- **Host Domain**: `social.example.org`
- **Account**: `admin@example.org`

#### Command

```bash
./FedSplitDomainChecker checkSplitDomain \
  --accountDomain example.org \
  --hostDomain social.example.org \
  --account admin@example.org -v
```

#### Output

**Success**:

```plaintext
❓ Starting split-domain check
✅ Checking Host-Meta endpoint...
✅ Host-Meta endpoint validation passed!
✅ Checking Nodeinfo endpoint...
✅ Nodeinfo endpoint validation passed!
✅ Checking Webfinger endpoint...
✅ Webfinger endpoint validation passed!
✅ Split-domain check passed!
```

**Failure**:

```plaintext
❓ Starting split-domain check
❌ Webfinger validation failed: unexpected status code: 404
❌ Split-domain check failed: Webfinger validation failed
```

---

## Development

### Testing

Unit tests are available to validate the tool’s functionality. Run tests using:

```bash
go test ./... -v
```

**Example Output**:

```plaintext
=== RUN   TestMakeRequest
--- PASS: TestMakeRequest (0.00s)
=== RUN   TestValidateWebfinger
--- PASS: TestValidateWebfinger (0.00s)
=== RUN   TestValidateSelfLink
--- PASS: TestValidateSelfLink (0.00s)
PASS
ok      github.com/0hlov3/FedSplitDomainChecker/cmd 0.005s
```

### Adding New Features

1. Fork the repository.
2. Create a new branch for your feature:

   ```bash
   git checkout -b feature-name
   ```

3. Add your feature and write corresponding tests.
4. Submit a pull request for review.

---

## Design Details

### Split-Domain Deployments

Split-domain deployments allow usernames like `@me@example.org` while hosting the actual service on a subdomain, such as `social.example.org`. This tool validates the necessary configurations:

- **`host-meta`**: Ensures the account domain properly redirects or responds with valid metadata.
- **`nodeinfo`**: Confirms the presence and validity of the node information endpoint.
- **`webfinger`**: Verifies that user accounts resolve to correct `self` links and ActivityPub data.

### Logging

The tool uses `zap` for structured logging:

- **Info-level logs**: Standard messages indicating the status of checks.
- **Debug-level logs**: Detailed messages for verbose mode.

---

## Contributing

Contributions are welcome! Please follow the steps below:

1. Fork the repository.
2. Create a new branch for your changes.
3. Add your feature or fix and include tests.
4. Submit a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

---

## Contact

For questions or suggestions, feel free to open an issue.# FedSplitDomainChecker
