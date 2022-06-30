CREATE DATABASE crypto_forum;
USE crypto_forum;

CREATE TABLE `users` (
    id                  INT             NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username            VARCHAR(32)     NOT NULL UNIQUE,
    old_usernames       JSON            NOT NULL,
    email               VARCHAR(64)     NOT NULL UNIQUE,
    last_ip             VARCHAR(64)     NOT NULL,
    user_groups         JSON            NOT NULL,
    last_login          TIMESTAMP       NOT NULL,
    join_date           TIMESTAMP       NOT NULL,
    wallet_address      VARCHAR(128)    NOT NULL UNIQUE,
    signature           TEXT            NOT NULL,   
    posts               INT             NOT NULL DEFAULT 0,
    threads             INT             NOT NULL DEFAULT 0,
    likes               BIGINT          NOT NULL DEFAULT 0,
    reputation          INT             NOT NULL DEFAULT 0,
    vouches             INT             NOT NULL DEFAULT 0,
    banned              BOOLEAN         NOT NULL DEFAULT false,
    ban_reason          TEXT            NOT NULL,
    ban_expiry          TIMESTAMP           NULL DEFAULT '1970-01-02',
    muted               BOOLEAN         NOT NULL DEFAULT false,
    mute_reason         TEXT            NOT NULL,
    mute_expiry         TIMESTAMP           NULL DEFAULT '1970-01-02',
    discord_id          VARCHAR(128)        NULL,
    discord_tag         VARCHAR(64)         NULL,
    telegram            VARCHAR(64)         NULL,
    change_username     BOOLEAN         NOT NULL
);

CREATE TABLE `ranks` (
    id                  INT             NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name                VARCHAR(32)     NOT NULL UNIQUE,
    price               INT             NOT NULL,
    min_vouches         INT             NOT NULL,
    min_rep             INT             NOT NULL,
    max_modifier        INT             NOT NULL,
);


CREATE TABLE `notifications` (
    id                  INT             NOT NULL AUTO_INCREMENT PRIMARY KEY,
    content             TEXT            NOT NULL,
    read                BOOLEAN         NOT NULL,
    link                TEXT            NOT NULL,
    creation_date       TIMESTAMP       NOT NULL
);

CREATE TABLE `messages` (
    id                  INT             NOT NULL AUTO_INCREMENT PRIMARY KEY,
    sender_id           INT             NOT NULL,
    recipient_id        INT             NOT NULL,
    content             TEXT            NOT NULL,  
    creation_date       TIMESTAMP       NOT NULL
);

CREATE TABLE `threads` (
    id                  INT             NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title               VARCHAR(128)    NOT NULL,
    section_id          INT             NOT NULL,
    author_id           INT             NOT NULL,
    creation_date       TIMESTAMP       NOT NULL,
    last_post           TIMESTAMP       NOT NULL,
    post_count          BIGINT          NOT NULL DEFAULT 0,
    hidden              BOOLEAN         NOT NULL DEFAULT false
);

CREATE TABLE `vouches` (
    id                  INT             NOT NULL AUTO_INCREMENT PRIMARY KEY,
    sender_id           INT             NOT NULL,
    recipient_id        INT             NOT NULL,
    message             TEXT            NOT NULL,
    deal_amount         INT             NOT NULL,
    show_amount         BOOLEAN         NOT NULL,
    creation_date       TIMESTAMP       NOT NULL
);

CREATE TABLE `feedbacks` (
    id                  INT             NOT NULL AUTO_INCREMENT PRIMARY KEY,
    sender_id           INT             NOT NULL,
    recipient_id        INT             NOT NULL,
    message             TEXT            NOT NULL,
    modifier            INT             NOT NULL,
    creation_date       TIMESTAMP       NOT NULL
);


CREATE TABLE `ip_bans` (
    id                  INT             NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ip_address          VARCHAR(64)     NOT NULL,
    ban_expiry          TIMESTAMP       NULL DEFAULT '1970-01-02',
    banned_by           VARCHAR(32)     NOT NULL,
    ban_reason          TEXT            NOT NULL
);

CREATE TABLE `thread_sections` (
    id                  INT             NOT NULL AUTO_INCREMENT PRIMARY KEY,
    parent_category     VARCHAR(32)     NOT NULL,
    parent_section      VARCHAR(128)    NOT NULL,
    section_name        VARCHAR(128)    NOT NULL UNIQUE
);

CREATE TABLE `statistics` (
    id                      INT         NOT NULL AUTO_INCREMENT PRIMARY KEY,
    total_threads           BIGINT      NOT NULL DEFAULT 0,
    total_posts             BIGINT      NOT NULL DEFAULT 0,
    total_users             BIGINT      NOT NULL DEFAULT 0
);


# TODO: On server side make a channel for actions that involve changing something but don't require feedback, that way its thread-safe

ALTER DATABASE crypto_forum CHARACTER SET utf8 COLLATE utf8_general_ci;

