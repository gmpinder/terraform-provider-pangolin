# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0](https://github.com/gmpinder/terraform-provider-pangolin/compare/v1.0.0...v1.1.0) (2026-03-21)


### ⚠ BREAKING CHANGES

* Complete refactor to create sites and fix issues
* Make the resource and target resources more extensible

### Bug Fixes

* add missing os import in target acceptance test ([3f34515](https://github.com/gmpinder/terraform-provider-pangolin/commit/3f345151cf6b814fb7d74c13c96efc47faecf517))
* lint errors and update provider version handling ([35fb230](https://github.com/gmpinder/terraform-provider-pangolin/commit/35fb2300495f89aa81e77af1556750264b5beea3))
* **pipeline:** Fix linting and tests ([92a69de](https://github.com/gmpinder/terraform-provider-pangolin/commit/92a69defc03c40793fef3697e5556d5ac76747cc))


### Miscellaneous Chores

* release 1.1.0 ([40fe14d](https://github.com/gmpinder/terraform-provider-pangolin/commit/40fe14d7a88be263c42a8419e281719907deff78))


### Code Refactoring

* Complete refactor to create sites and fix issues ([ed8a8d5](https://github.com/gmpinder/terraform-provider-pangolin/commit/ed8a8d5e762b213a9e1adbd5a3e3ef4a1d919fd0))
* Make the resource and target resources more extensible ([1dec3e1](https://github.com/gmpinder/terraform-provider-pangolin/commit/1dec3e187db633c2fc5ecc8b147b2ae52845299e))

## 1.0.0 (2026-02-08)


### Bug Fixes

* add missing os import in target acceptance test ([3f34515](https://github.com/groteck/pangolin-tf/commit/3f345151cf6b814fb7d74c13c96efc47faecf517))
* lint errors and update provider version handling ([35fb230](https://github.com/groteck/pangolin-tf/commit/35fb2300495f89aa81e77af1556750264b5beea3))
* **pipeline:** Fix linting and tests ([92a69de](https://github.com/groteck/pangolin-tf/commit/92a69defc03c40793fef3697e5556d5ac76747cc))

## [0.1.0] - 2026-02-08

### Added
- Initial release of the Pangolin Terraform provider.
- Supported Resources:
    - `pangolin_role`: Manage organization roles.
    - `pangolin_site_resource`: Manage Host/CIDR resources.
    - `pangolin_resource`: Manage App resources (HTTP/TCP/UDP).
    - `pangolin_target`: Manage backend targets for resources.
- Supported Data Sources:
    - `pangolin_role`: Fetch role ID by name.
    - `pangolin_site`: Fetch site ID by name.
- Provider authentication via token and environment variables.
- Automated documentation generation via `tfplugindocs`.
- Multi-architecture release automation via GoReleaser.
- Docker-based acceptance testing suite.
