k8S Community Bundle Validator
==

## Overview

It is an external validator which can be used to ensure that an OLM bundle is respecting
the specific criteria to publish in the [Kubernetes Community Operator](https://github.com/k8s-operatorhub/community-operators). 

Note that a distribution must attend the common requirements defied to integrate the Operator project
with OLM and then, this validator will check the bundle which is intended to be published on [OperatorHub.io](https://operatorhub.io/) 
catalogs.

The common criteria to publish are defined and implemented in the [https://github.com/operator-framework/api](https://github.com/operator-framework/api) 
which is also used by this project in order to try to ensure and respect the defined standards. Users are able to test their
bundles with the common and criteria to distributed in OLM by running `operator-sdk bundle validate ./bundle --select-optional suite=operatorframework` 
and using [Operator-SDK][operator-sdk].

> **The purpose of this validator is to ensure 
and centralize any rule and criteria which is specific to publish in 
K8s Community Operator Catalog ([OperatorHub.io](https://operatorhub.io/)).**

**NOTE** We have an [EP in WIP](https://github.com/operator-framework/enhancements/pull/98). The idea is in the future
[Operator-SDK][operator-sdk] also be able to run this validator. 

## Install

Download the binary from the release page. 
Following the steps to allow you do that via command line:

1. Set platform information:

```sh
export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac)
export OS=$(uname | awk '{print tolower($0)}')
```

2. Download the binary:

```sh
export VALIDADOR_URL=https://github.com/k8s-operatorhub/bundle-validator/releases/download/{tagVersion}
curl -LO ${VALIDADOR_URL}/${OS}-${ARCH}-amd64-k8s-community-bundle-validator
chmod +x amd64-k8s-community-bundle-validator
```

### From source-code

Run `make install` to be able to build and install the binary locally. 

## Usage

You can test this validator by running:

```sh
$ k8s-community-bundle-validator <bundle-path>
```

**NOTE** You can use the option `--output=json-alpha1` to output the format in `JSON` format.

## How to check what is validated with this project?

The documentation ought to get done in this project source code in order to generate the Golang docs. 

## Release

Create a new tag and publish in the repository. It will call the GitHub action release and the
artifacts will be built and publish in the release page automatically after few minutes. 

[operator-sdk]: https://github.com/operator-framework/operator-sdk