# Decode JSON Web Tokens

A simple CLI tool to decode JWTs.  
Currently only supports JSON Web Signatures (JWSs).

## Prerequisites

- Golang

## Installation

```
go install github.com/HarikrishnanBalagopal/jwtdecode@v1
```

## Usage

```
jwtdecode "${MY_JWT}"
```

or

```
echo "${MY_JWT}" | jwtdecode -
```
