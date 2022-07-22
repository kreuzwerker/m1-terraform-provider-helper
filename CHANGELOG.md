
<a name="0.7.4"></a>
## [0.7.4] - 2022-07-22
### features
- Implement parsing of provider build output. ([#82](https://github.com/kreuzwerker/m1-terraform-provider-helper/issues/82)) ([ac1b248](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/ac1b248afb240726a3e6535973a934161a808f40))


<a name="0.7.3"></a>
## [0.7.3] - 2022-07-21

<a name="0.7.2"></a>
## [0.7.2] - 2022-07-19
### bug fixes
- Use raw `terraform version` for backwards compatibility. ([#76](https://github.com/kreuzwerker/m1-terraform-provider-helper/issues/76)) ([8ef73e7](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/8ef73e73adb2a26bc252fcf914ce2b0e867d7474))


<a name="0.7.1"></a>
## [0.7.1] - 2022-07-19
### bug fixes
- Correct plugin installation path for older tf versions. ([#74](https://github.com/kreuzwerker/m1-terraform-provider-helper/issues/74)) ([be46959](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/be469594a2156e99395f821d0a0cc46ced903840))


<a name="0.7.0"></a>
## [0.7.0] - 2022-07-07

<a name="0.6.2"></a>
## [0.6.2] - 2022-06-22
### bug fixes
- Fix linting error. ([1aa7344](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/1aa73441d83aea882a059894ad40b19940e9ee73))
- Improve command to checkout default branch ([bf729f8](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/bf729f8cfcbd9fd3ff3377cdf6e3a3337c760536))


<a name="0.6.1"></a>
## [0.6.1] - 2022-06-21
### features
- introduce go-version for semver handling. ([7c7a263](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/7c7a263376edde733275c3c801dde1df2bd7c3ec))


<a name="0.6.0"></a>
## [0.6.0] - 2022-05-27
### bug fixes
- error when GOPATH not set ([46aeaeb](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/46aeaeb367a018d595896680d7350651f5eb5f88))


<a name="0.5.1"></a>
## [0.5.1] - 2022-02-14

<a name="0.5.0"></a>
## [0.5.0] - 2022-02-10
### bug fixes
- Correct linting errors. ([11cbe6f](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/11cbe6ff29ab2c56fd800e743a02f60a69da5fe3))
- Use https for cloning. ([c55b20a](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/c55b20aa2b03edb986bd200c4f58e84b3065f0b3))


<a name="0.4.0"></a>
## [0.4.0] - 2022-01-27
### features
- Pass custom build command to install step. ([21a4651](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/21a46511541c24ecea0c2d3bb9d7b9916b41a8a9))


<a name="0.3.1"></a>
## [0.3.1] - 2022-01-12
### bug fixes
- Adapt buildCommand test to new function signature. ([f6e282d](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/f6e282dc6258f31cc5343035201993dd249c8920))
- Add build command for helm provider. ([9169996](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/91699965f1c01d68109edfa9a2fdbc48c6f4f9d3))


<a name="0.3.0"></a>
## [0.3.0] - 2022-01-12
### bug fixes
- Rework buildCommand structure to fix flaky test. ([09f6303](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/09f63039b501319d48a8c8765254735f6e2e4f1e))
- Path in .terraform.d should contain normalized version. ([7f31380](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/7f3138028e17f081c7d454291e31d81507af6461))

### features
- Add basic list command. ([c005aa9](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/c005aa90a2eb9c3aa44ce02c0e03756cc0c11a5b))


<a name="0.2.5"></a>
## [0.2.5] - 2022-01-08
### bug fixes
- Provide different build commands for different provider versions. ([1a619ca](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/1a619caa969e6909e87d07749c112ee866c1bbf3))
- respect GOPATH environment variable ([892628c](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/892628cd3e8aed66fc1a3142ce15687d53f798a6))


<a name="0.2.4"></a>
## [0.2.4] - 2021-12-24
### bug fixes
- Run git commands in correct directory. ([6cbfcf2](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/6cbfcf23d386c61841d760f040abbfc44b91b9b7))


<a name="0.2.3"></a>
## [0.2.3] - 2021-11-26

<a name="0.2.2"></a>
## [0.2.2] - 2021-11-12

<a name="0.2.1"></a>
## [0.2.1] - 2021-08-11
### bug fixes
- add linux os and amd64 arch for brew formula ([ac07195](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/ac0719515d3cd66d6585c7aa5c74d7416fc0d88c))


<a name="0.2.0"></a>
## [0.2.0] - 2021-07-01
### features
- support for provider specific build commands. ([2c34eb8](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/2c34eb86756d5c2208258e0cadabf273bd64cba5))


<a name="0.1.3"></a>
## [0.1.3] - 2021-06-30
### code refactoring
- Move status application logic to app.go. ([ab65a7e](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/ab65a7e2a9f8a0c977dc3a2c70951155016a5a90))
- Rewrite App.Config and make Out configurable. ([ee65c8d](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/ee65c8d7cf7bc409612dd37547569af0c953f724))


<a name="0.1.2"></a>
## [0.1.2] - 2021-06-25
### bug fixes
- Add location of main func to goreleaser config. ([ee0ac7c](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/ee0ac7c7463d7db0efd10e7d84ab6e1f99f728e2))


<a name="0.1.1"></a>
## [0.1.1] - 2021-06-25
### bug fixes
- Update gotestsum version to 1.6.4 for darwin_arm64 support ([fc40532](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/fc40532e2bd793a412244f7661f255c45c8d5bf7))


<a name="0.1.0"></a>
## 0.1.0 - 2021-06-25
### features
- Initial version ([ae159b7](https://github.com/kreuzwerker/m1-terraform-provider-helper/commit/ae159b754681990f2339b95835f79452cb5fc455))

