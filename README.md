# Anthropoi
A simple accounts package and management tool.

## What is it
This package sets up and manages user accounts with multi-site support and per-site groups.

## Requirements
This package was made for use with PostgreSQL. CockroachDB probably won't work because of triggers.

## Installing
Run this command to get the package:

`go get -u github.com/Urethramancer/anthropoi`

And run this to compile and install the management command:

`go get -u github.com/Urethramancer/anthropoi/cmd/anthro`

## Using the package
(See built-in documentation for parameter information.)
- `New()` creates a DBM structure to use for all further calls.
- `ConnectionString()` rebuilds and returns a string based on the internally stored parameters.
- `Connect()` opens the connection to the specified host, or localhost.
- `DatabaseExists()` checks if there is a database of the specified name.
- `Create()` creates a new account database with the specified name.
- `InitDatabase()` sets up a new database with tables and triggers.

`New()` will use reasonable defaults for its connection string:
- `host`: localhost
- `port`: 5432
- `user`: postgres
- `password`: unused if blank
- `name`: unused if blank
- `mode`: disable if blank
