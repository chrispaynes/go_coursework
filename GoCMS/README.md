#GoCMS

##Features
* Creating a Go Library
* Postgres DB Connectivity
* Dynamic Content Rendering
* HTML Templating
* Routing HTTP Requests
* Testing

##TODO
*Fix SQL DB initialize access privileges with gocms user

##Installation
1. Install Postgres: [Detailed Postgres Installation Guides](https://wiki.postgresql.org/wiki/Detailed_installation_guides)
2. Start the Postgres Server
3. Run init.sql script to create and initialize the Postgres database setup
```
$ cd ./src/GoCMS
$ psql
gocms=# \i init.sql
CREATE ROLE
CREATE DATABASE
GRANT
CREATE TABLE
CREATE TABLE
CREATE TABLE
```

##Usage
* After installing Postgres, starting the postgres and creating the database...
`$ cd go run ./src/GoCMS/cmd/main.go`

* Allow for incoming network connections on Port 3000
* Visit the application on localhost:3000
* Visit  `localhost:3000/`  to view the index.
* Visit  `localhost:3000/new`  to add new pages and posts
* Visit  `localhost:3000/page/`  to view and index of pages
* Visit  `localhost:3000/page/1`  to visit a specific page

##Running Tests
