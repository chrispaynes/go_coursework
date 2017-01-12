-- Create new user and password for development. 
CREATE USER gocms WITH PASSWORD 'gocms';

-- Create new development database.
CREATE DATABASE gocms;

-- Grant all gocms DB privileges to gocms user.
GRANT ALL PRIVILEGES ON DATABASE gocms to gocms;

-- Create new table to store pages.
CREATE TABLE IF NOT EXISTS pages(
  id            SERIAL     PRIMARY KEY,
  title         TEXT       NOT NULL,
  content       TEXT       NOT NULL
);

-- Create new table to store posts.
CREATE TABLE IF NOT EXISTS posts(
  id            SERIAL     PRIMARY KEY,
  title         TEXT       NOT NULL,
  content       TEXT       NOT NULL,
  date_created  DATE       NOT NULL
);

-- Create new table to store comments.
CREATE TABLE IF NOT EXISTS comments(
  id            SERIAL     PRIMARY KEY,
  author        TEXT       NOT NULL,
  content       TEXT       NOT NULL,
  date_created  DATE       NOT NULL,
  post_id       INT        references posts(id)
);

