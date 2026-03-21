# Terraform Provider Pangolin

A Terraform provider for managing [Pangolin](https://pangolin.net) resources.

## Architecture Decisions

### 1. Framework Selection
This provider is built using the **Terraform Plugin Framework** (instead of the older SDKv2). 
- **Why**: It provides better type safety, improved validation capabilities, and is the current standard recommended by HashiCorp for new providers.

### 2. Internal Client Decoupling
The API logic is isolated in `internal/client/`.
- **Why**: This separates the HTTP/JSON concerns from the Terraform state management. It makes the code more maintainable and allows for easier unit testing of the API client independent of the Terraform lifecycle.

### 3. Authentication
The provider uses Bearer Token authentication as required by the Pangolin Integration API. The token is marked as `sensitive` in the schema to ensure it doesn't leak into logs.

## Requirements

- [Go](https://golang.org/doc/install) >= 1.24
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0

## Compatibility

| Pangolin Version | Provider Version |
| ---------------- | ---------------- |
| >= v1.16.2       | >= 1.0.0         |

## Building The Provider

To compile the provider locally, run:

```bash
go build -o terraform-provider-pangolin
```

## Using the Provider Locally

To test the provider without publishing it, you can use Terraform's `dev_overrides` feature. Create or edit your `~/.terraformrc` file:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/gmpinder/pangolin" = "/path/to/your/project/pangolin-tf"
  }
  direct {}
}
```

## Configuration

This provider communicates with the **Pangolin Integration API**.

### Enabling Integration API in Pangolin

You must enable the Integration API in your Pangolin server configuration (`config.yml`).

```yaml
flags:
  enable_integration_api: true

server:
  integration_port: 3003 # Optional, defaults to 3003
```

Ensure the port `3003` (or your configured port) is exposed and accessible to the machine running Terraform.

```hcl
provider "pangolin" {
  token    = "your-api-token"
  base_url = "https://your-pangolin-instance.com:3003/v1" # Point to the Integration API
}
```

### Provider Configuration

```hcl
provider "pangolin" {
  token    = "your-api-token"
  base_url = "https://api.pangolin.net/v1" # Optional
}
```

## Examples

See the [examples/](examples/) directory for a full configuration.

## Contributing

### Commit Messages

This project enforces [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/). All commit messages must follow this specification (e.g., `feat: add new resource`, `fix: resolve crash`).

This format is used to automate the semantic versioning and changelog generation via the Release Please workflow.

### Release Process

Releases are automated using [Release Please](https://github.com/google-github-actions/release-please).

1.  **Commit**: Push changes using Conventional Commits.
2.  **Automated PR**: The Release Please bot analyzes new commits and automatically creates a "Release PR" with the updated `CHANGELOG.md` and calculated version bump (e.g., `chore(main): release v0.2.0`).
3.  **Merge**: When the Release PR is merged, the bot:
    *   Creates a GitHub Release.
    *   Creates a Git Tag.
4.  **Artifacts**: The creation of the tag triggers GoReleaser to build and upload the provider binaries to the release.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
