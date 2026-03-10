# Change log

All notable changes will be documented in this file. This project adheres to [Semantic Versioning](http://semver.org).

## [1.2.1](https://github.com/launchdarkly/go-configtypes/compare/v1.2.0...v1.2.1) (2026-03-10)


### Bug Fixes

* **deps:** bump supported Go versions to 1.26 and 1.25 ([#17](https://github.com/launchdarkly/go-configtypes/issues/17)) ([817a8ca](https://github.com/launchdarkly/go-configtypes/commit/817a8cafd39c1bcfaee84b2327d41735e8815cb6))

## [1.2.0](https://github.com/launchdarkly/go-configtypes/compare/1.1.0...v1.2.0) (2024-07-03)


### Features

* Add OptBase2Bytes type ([#7](https://github.com/launchdarkly/go-configtypes/issues/7)) ([f439419](https://github.com/launchdarkly/go-configtypes/commit/f4394199962d51066318cf86018c307a940a5e9b))


### Bug Fixes

* **deps:** Bump supported Go versions to 1.21 and 1.22 ([#9](https://github.com/launchdarkly/go-configtypes/issues/9)) ([f97f175](https://github.com/launchdarkly/go-configtypes/commit/f97f1750e1c8b5d5c27818d61d92bbc439545539))

## [1.1.0] - 2020-09-18
### Added:
- `OptFloat64` type, and support for `float64` in `VarReader`.

### Fixed:
- The error message for invalid duration strings had an incorrect list of supported formats.

## [1.0.0] - 2020-07-17
Initial release.
