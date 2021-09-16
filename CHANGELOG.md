## [0.15.0](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.14.0...v0.15.0) (2021-09-16)


### Features

* **group:** add support for snmp notifications ([058cb2c](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/058cb2c0a362956c0e519e337c479c38801cd0ff)), closes [#66](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/66)


### Bug Fixes

* **scheduler:** remove seconds from excluded time interval format ([7193fba](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/7193fba4957df09883c0983e52e55dbbf1a51097)), closes [#67](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/67)

## [0.14.0](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.13.0...v0.14.0) (2021-07-26)


### Features

* **filter:** adding support for ignored error code ranges ([f796de5](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/f796de5030c7ac62e6d5240d1e6d99896d204ee6)), closes [#64](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/64)

## [0.13.0](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.12.0...v0.13.0) (2021-07-22)


### Features

* **filter:** implementing filter resource and data source ([8231227](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/8231227e4a6d1c4130e805508ce807be858b20dd)), closes [#14](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/14) [#15](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/15)
* **platform:** adding platform id validation ([6c13f1a](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/6c13f1a46fdc10272293dfadb8642ade28e94850)), closes [#10](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/10)

## [0.12.0](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.11.1...v0.12.0) (2021-07-14)


### Features

* **group:** adding support for slack, teams and alertops ([6a06e4f](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/6a06e4f6cce9e9b936e56833bb2565cf578c9f25)), closes [#54](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/54) [#55](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/55) [#59](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/59)

### [0.11.1](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.11.0...v0.11.1) (2021-07-14)


### Bug Fixes

* **docs:** scheduler docs formatting ([d8d44be](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/d8d44bea0a0d98ea58d4c69054099a47676f46e9))

## [0.11.0](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.10.2...v0.11.0) (2021-07-14)


### Features

* **group:** deprecating pager notifications ([377fda9](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/377fda9018c40e1ef94d14b226ebbe2e7912960b)), closes [#56](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/56)

### [0.10.2](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.10.1...v0.10.2) (2021-07-14)


### Bug Fixes

* **ci:** only run release on tag ref ([08fefd6](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/08fefd6c46a9de73868b07f785f38301e6dec1aa))
* **docs:** fixing incorrect resource documentation ([afd2919](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/afd2919381bcae6b3bc2ea5c5fe5cab7d4bf907f))

### [0.10.1](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.10.0...v0.10.1) (2021-07-13)

## [0.10.0](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.9.0...v0.10.0) (2021-07-13)


### Features

* **scheduler:** refactoring intervals to be user friendly ([9b8ca21](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/9b8ca21c1b29ce485dad69be86b0867f0d7ba330)), closes [#39](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/39)

## [0.9.0](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.8.8...v0.9.0) (2021-07-13)


### Features

* changing inputs from lists to sets ([c0b87de](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/c0b87de4ad0490f481263da3fb682d2ca4c21f11)), closes [#12](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/12)


### Bug Fixes

* **ci:** adding release notes to goreleaser ([7a3eaf2](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/7a3eaf208e1a8482887ad6616c450cf8cc4fea45))
* **ci:** workflow tag trigger ([990d425](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/990d4255fc0ed448d5b9e56d1112f4e8c2bc79d7))
* **device:** marking locations as required input ([c0304a6](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/c0304a69f73651bd7e43ada0f82f5213d5ec51b5))
* **task:** dns resolve mode default ([31a9a5e](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/31a9a5ede70a2d3391f135fbfbe75341c8cfe9de))
* **task:** incorrectly setting header params ([9485fb2](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/9485fb234b722f6fed6f69fc4765930d08fe7ae9))

## [0.8.0](https://github.com/rymancl/terraform-provider-dotcommonitor/compare/v0.7.0...v0.8.0) (2021-07-08)


### Features

* **ci:** automated tagging and changelog ([2e1873f](https://github.com/rymancl/terraform-provider-dotcommonitor/commit/2e1873f5af1b4915f008477b48c8fe08a19c7973)), closes [#28](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/28)
