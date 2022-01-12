
<a name="0.3.0"></a>
## [0.3.0] - 2022-01-12
### bug fixes
- Rework buildCommand structure to fix flaky test. (09f6303)
- Path in .terraform.d should contain normalized version. (7f31380)

### features
- Add basic list command. (c005aa9)


<a name="0.2.5"></a>
## [0.2.5] - 2022-01-08
### bug fixes
- Provide different build commands for different provider versions. (1a619ca)
- respect GOPATH environment variable (892628c)


<a name="0.2.4"></a>
## [0.2.4] - 2021-12-24
### bug fixes
- Run git commands in correct directory. (6cbfcf2)


<a name="0.2.3"></a>
## [0.2.3] - 2021-11-26
### bug fixes
 - Add building of linux binary to fix brew-tap issues.

<a name="0.2.2"></a>
## [0.2.2] - 2021-11-12
### features
 - Clean provider repo when running install command
 - Pull main branch when providing no version

<a name="0.2.1"></a>
## [0.2.1] - 2021-08-11
### bug fixes
- add linux os and amd64 arch for brew formula (ac07195)


<a name="0.2.0"></a>
## [0.2.0] - 2021-07-01
### features
- support for provider specific build commands. (2c34eb8)


<a name="0.1.3"></a>
## [0.1.3] - 2021-06-30
### code refactoring
- Move status application logic to app.go. (ab65a7e)
- Rewrite App.Config and make Out configurable. (ee65c8d)


<a name="0.1.2"></a>
## [0.1.2] - 2021-06-25
### bug fixes
- Add location of main func to goreleaser config. (ee0ac7c)


<a name="0.1.1"></a>
## [0.1.1] - 2021-06-25
### bug fixes
- Update gotestsum version to 1.6.4 for darwin_arm64 support (fc40532)


<a name="0.1.0"></a>
## 0.1.0 - 2021-06-25
### features
- Initial version (ae159b7)

