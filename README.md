# Maja Service
![Build Status](https://travis-ci.org/nock/nock.svg)
![Coverage Status](http://img.shields.io/badge/coverage-100%25-brightgreen.svg)


## Installation
Run the following command in your project

~~~bash
  make install
~~~
If you have docker installed, you can run the following command to run the project:

~~~bash
  make docker-start
~~~

Otherwise, you can run the following command to run the project:

~~~bash
  make start
~~~

## List of the make commands

~~~bash
# Run test convey analyzer
make convey

# Run test with docker
make docker-test

# Run lint with docker
make docker-lint

# Run test
make test

# Run lint
make lint

# Run test coverage
make testcov

# Run cli
make cli

# Run cli with docker
make docker-cli

# Make swagger docs
make swagger

# Run project with watch
make start

# Run project with watch with docker
make docker-start

# Install project
make install

# Stop project with docker
make docker-stop

# Run project without watch
make serve

# Run project without watch with docker
make docker-serve

# Stop project without watch with docker
make docker-drop

# Show logs of service with docker
make docker-logs

# Prepare the last updated git hooks for the project
make prepare-hooks
~~~

## Tech Stack
**Server:** Golang

## Licence
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
