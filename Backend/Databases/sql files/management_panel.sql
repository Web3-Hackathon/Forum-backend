CREATE DATABASE forum_management;

CREATE TABLE `users` (
    id              INT     NOT NULL    PRIMARY KEY AUTO_INCREMENT,
    username        TEXT    NOT NULL,
    password        TEXT    NOT NULL,
    permission      TINYINT NOT NULL    DEFAULT 0
);

CREATE TABLE `logs` (
    id      INT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
    action  TEXT NOT NULL,
    user    TEXT NOT NULL
);