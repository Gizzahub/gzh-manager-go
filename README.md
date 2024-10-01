# gzh-manager

종합 CLI 툴

<div style="text-align: center;">
A general purpose project for easy dev.
<br>
<br>
This template serves as a starting point for golang commandline applications it is based on golang projects that I consider high quality and various other useful blog posts that helped me understanding golang better.
<br>
<br>
<img src="https://github.com/gizzahub/gzh-manager-go/actions/workflows/test.yml/badge.svg" alt="drawing"/>
<img src="https://github.com/gizzahub/gzh-manager-go/actions/workflows/lint.yml/badge.svg" alt="drawing"/>
<img src="https://pkg.go.dev/badge/github.com/gizzahub/gzh-manager-go.svg" alt="drawing"/>
<img src="https://codecov.io/gh/Gizzahub/gzh-manager-go/branch/main/graph/badge.svg" alt="drawing"/>
<img src="https://img.shields.io/github/v/release/Gizzahub/gzh-manager-go" alt="drawing"/>
<img src="https://img.shields.io/docker/pulls/Gizzahub/gzh-manager-go" alt="drawing"/>
<img src="https://img.shields.io/github/downloads/Gizzahub/gzh-manager-go/total.svg" alt="drawing"/>
</div>

# Table of Contents
<!--ts-->
   * [gzh-manager-go](#gzh-manager-go)
   * [Features](#features)
   * [Project Layout](#project-layout)
   * [How to use this template](#how-to-use-this-template)
   * [Demo Application](#demo-application)
   * [Makefile Targets](#makefile-targets)
   * [Contribute](#contribute)

<!-- Added by: morelly_t1, at: Tue 10 Aug 2021 08:54:24 AM CEST -->

<!--te-->

# Features
- [goreleaser](https://goreleaser.com/) with `deb.` and `.rpm` packer and container (`docker.hub` and `ghcr.io`) releasing including `manpages` and `shell completions` and grouped Changelog generation.
- [golangci-lint](https://golangci-lint.run/) for linting and formatting
- [Github Actions](.github/worflows) Stages (Lint, Test (`windows`, `linux`, `mac-os`), Build, Release) 
- [Gitlab CI](.gitlab-ci.yml) Configuration (Lint, Test, Build, Release)
- [cobra](https://cobra.dev/) example setup including tests
- [Makefile](Makefile) - with various useful targets and documentation (see Makefile Targets)
- [Github Pages](_config.yml) using [jekyll-theme-minimal](https://github.com/pages-themes/minimal) (checkout [https://Gizzahub.github.io/gzh-manager-go/](https://Gizzahub.github.io/gzh-manager-go/))
- Useful `README.md` badges
- [pre-commit-hooks](https://pre-commit.com/) for formatting and validating code before committing

# Project Layout
* [assets/](https://pkg.go.dev/github.com/gizzahub/gzh-manager-go/assets) => docs, images, etc
* [cmd/](https://pkg.go.dev/github.com/gizzahub/gzh-manager-go/cmd)  => commandline configurartions (flags, subcommands)
* [pkg/](https://pkg.go.dev/github.com/gizzahub/gzh-manager-go/pkg)  => packages that are okay to import for other projects
* [internal/](https://pkg.go.dev/github.com/gizzahub/gzh-manager-go/pkg)  => packages that are only for project internal purposes
- [`tools/`](tools/) => for automatically shipping all required dependencies when running `go get` (or `make bootstrap`) such as `golang-ci-lint` (see: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)
)
- [`scripts/`](scripts/) => build scripts 

# How to use this template
```sh```

In order to make the CI work you will need to have the following Secrets in your repository defined:

Repository  -> Settings -> Secrets & variables -> `CODECOV_TOKEN`, `DOCKERHUB_TOKEN` & `DOCKERHUB_USERNAME`

# Demo Application

```sh
$> gzh-manager -h
golang-cli cli application by managing gzh-manager

Usage:
  gzh [flags]
  gzh [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  setclone    세트로클론 한다는 뜻의 순수한 지구말입니다
  help        Help about any command
  version     gzh-manager version

Flags:
  -h, --help   help for gzh-manager

Use "gzh-manager [command] --help" for more information about a command.
```

먼저.. .원하는 경로에 설정파일 만들기. 참고
[setclone.yaml](pkg/setclone/setclone.yaml)

```sh
$> gzh setclone -t $HOME/mywork

이건... 안됨
$> gzh setclone -t ./mywork
$> gzh setclone -t $HOME/mywork
$> gzh setclone -t ~/mywork
```

# Makefile Targets
```sh
$> make
bootstrap                      install build deps
build                          build golang binary
clean                          clean up environment
cover                          display test coverage
docker-build                   dockerize golang application
fmt                            format go files
help                           list makefile targets
install                        install golang binary
lint                           lint go files
pre-commit                     run pre-commit hooks
run                            run the app
test                           display test coverage
```

# Contribute
If you find issues in that setup or have some nice features / improvements, I would welcome an issue or a PR :)

# 기록용
템플릿을 .. `https://github.com/FalcoSuessgott/golang-cli-template` 를 썼는데.... `https://github.com/create-go-app/cli` 가 더 좋은듯? 심플해서..
뭔가 복잡해졌다. 
