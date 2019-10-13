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
- `mode`: enable if set to "enable", disable otherwise

## Using the CLI tool
The command line tool, `anthro`, can be used to manage users, profiles, groups and permissions, saving you from having to write your own manager.

### init
Running `anthro init` will set up the initial database. If you want to drop the current database, supply the `-D` flag. For more advanced features, such as dumping a backup of the database, use the `pgsql` command from the relevant PostgreSQL package.

### user
The `user`command has subcommands for user account management. A user account has just the bare minimum details about a user, such as a display username, primary e-mail, login password, name and special data for some sites. There are two JSON fields, data and tokens, available to use for whatever you need. Methods to make use of them aren't currently implemented, but using PostgreSQL's JSON lookup is fairly straightforward.

#### user add
The `user add` subcommand takes a username at minimum. A password and salt will be generated and stored, and the password will be displayed in the terminal. Write it down or lose it!

Optional arguments are e-mail, first and last name and a cost, which is the complexity to use for hashing the password. The current minimum amount is 10, which gives a decent amount for testing. 11+ is recommended for production use, especially on very fast server hardware. The time roughly doubles for each increase by 1.

#### user remove
The `user remove` (or `user rm`) subcommand removes a user by ID or name.

#### user edit
Line by line editing of user fields, except password.

#### user list
This lists all users. More flags to filter on will be added in the near future.

### profile
The `profile` command has subcommands for per-site profiles. Profiles are useful when you want one system to handle many domains with different profiles, containing different access rights, but which should share common logins. This is useful for blogging systems where different subdomains are used for different subjects, or to create a domain admin system for e-mail, for example.

#### profile add
The `profile add` subcommand adds a profile to a user. Permissions and groups are handled separately.

#### profile remove
Removes a profile, effectively removing access to a domain for a user.

#### profile setgroups
This manages groups for a profile, i.e. per-site permissions. Access rights are handled via groups, while profiles can contain collections of groups.

#### profile copy
This allows you to copy the non-personal parts of a profile from one user to another to quickly set permissions.

#### profile list
Lists profiles in the database, with optional filtering by site and by user.
