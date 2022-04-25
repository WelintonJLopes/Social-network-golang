CREATE DATABASE IF NOT EXISTS social_network;
USE social_network;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS publications;

CREATE TABLE users (
    id int auto_increment primary key,
    name varchar(50) not null,
    nick varchar(50) not null unique,
    email varchar(50) not null unique,
    password varchar(100) not null unique,
    createdat timestamp default current_timestamp()
) ENGINE=INNODB;

CREATE TABLE followers (
    user_id int not null,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    follower_id int not null,
    FOREIGN KEY (follower_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    PRIMARY KEY (user_id, follower_id)
) ENGINE=INNODB;

CREATE TABLE publications (
    id int auto_increment primary key,
    title varchar(50) not null,
    content varchar(300) not null,

    author_id int not null,
    FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE,

    likes int default 0,
    createdat timestamp default current_timestamp
) ENGINE=INNODB;
