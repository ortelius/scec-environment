# Ortelius v11 Environment Microservice

> Version 11.0.0

RestAPI for the Environment Object
![Release](https://img.shields.io/github/v/release/ortelius/scec-environment?sort=semver)
![license](https://img.shields.io/github/license/ortelius/.github)

![Build](https://img.shields.io/github/actions/workflow/status/ortelius/scec-environment/build-push-chart.yml)
[![MegaLinter](https://github.com/ortelius/scec-environment/workflows/MegaLinter/badge.svg?branch=main)](https://github.com/ortelius/scec-environment/actions?query=workflow%3AMegaLinter+branch%3Amain)
![CodeQL](https://github.com/ortelius/scec-environment/workflows/CodeQL/badge.svg)
[![OpenSSF-Scorecard](https://api.securityscorecards.dev/projects/github.com/ortelius/scec-environment/badge)](https://api.securityscorecards.dev/projects/github.com/ortelius/scec-environment)

![Discord](https://img.shields.io/discord/722468819091849316)

## Path Table

| Method | Path | Description |
| --- | --- | --- |
| GET | [/msapi/environment](#getmsapienvironment) | Get a List of Environments |
| POST | [/msapi/environment](#postmsapienvironment) | Create a Environment |
| GET | [/msapi/environment/:key](#getmsapienvironmentkey) | Get a Environment |

## Reference Table

| Name | Path | Description |
| --- | --- | --- |

## Path Details

***

### [GET]/msapi/environment

- Summary  
Get a List of Environments

- Description  
Get a list of environments for the user.

#### Responses

- 200 OK

***

### [POST]/msapi/environment

- Summary  
Create a Environment

- Description  
Create a new Environment and persist it

#### Responses

- 200 OK

***

### [GET]/msapi/environment/:key

- Summary  
Get a Environment

- Description  
Get a environment based on the _key or name.

#### Responses

- 200 OK

## References
