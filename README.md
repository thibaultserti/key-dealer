# Key dealer

## Badges

[![Build Status](https://github.com/thibaultserti/key-dealer/actions/workflows/release.yaml/badge.svg)](https://github.com/thibaultserti/key-dealer/actions/workflows/release.yaml)
[![License](https://img.shields.io/github/license/thibaultserti/key-dealer)](/LICENSE)
[![Release](https://img.shields.io/github/release/thibaultserti/key-dealer.svg)](https://github.com/thibaultserti/key-dealer/releases/latest)
[![GitHub Releases Stats of key-dealer](https://img.shields.io/github/downloads/thibaultserti/key-dealer/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=thibaultserti&repository=key-dealer)

[![Maintainability](https://api.codeclimate.com/v1/badges/4133d7da3d73fa0c0884/maintainability)](https://codeclimate.com/github/thibaultserti/key-dealer/maintainability)
[![codecov](https://codecov.io/gh/thibaultserti/key-dealer/branch/main/graph/badge.svg?token=5BO47LR632)](https://codecov.io/gh/thibaultserti/key-dealer)
[![Go Report Card](https://goreportcard.com/badge/github.com/thibaultserti/test-saas-ci)](https://goreportcard.com/report/github.com/thibaultserti/key-dealer)

## How to

Key dealer exposes static json GCP keys via HTTP. It automatically revokes and renew them once a day.
The api endpoint to request a key is `/key/<project>/<sa_name>`
