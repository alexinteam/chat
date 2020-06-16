-- DB

CREATE DATABASE chat;
GRANT ALL PRIVILEGES ON DATABASE db_gochat TO pguser;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO pguser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO pguser;

-- Tables

CREATE TABLE auth_user (
    id serial NOT NULL,
    full_name varchar(60),
    username varchar(60) NOT NULL,
    email varchar(60) NOT NULL,
    password varchar(60) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE auth_session (
    id serial NOT NULL,
    key varchar(64) NOT NULL,
    user_id integer NOT NULL,
    create_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expire_date timestamp NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE room (
    id serial NOT NULL,
    name varchar(60) NOT NULL,
    create_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE message (
    id serial NOT NULL,
    room_id integer NOT NULL,
    sender_id integer NOT NULL,
    recipient_id integer,
    text text NOT NULL,
    send_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    recieve_date timestamp,
    PRIMARY KEY (id)
);

CREATE TABLE db_version (
    id serial NOT NULL,
    version integer NOT NULL,
    PRIMARY KEY (id)
);

-- Populate

-- password: 123
INSERT INTO auth_user (id, full_name, username, email, password)
    VALUES (1, 'Administrator', 'admin', 'admin@gochat.local', 'admin');
INSERT INTO auth_user (id, full_name, username, email, password)
    VALUES (2, 'Moderator', 'moder', 'moder@gochat.local', 'moder');
INSERT INTO auth_user (id, full_name, username, email, password)
    VALUES (3, 'User #1', 'user1', 'user1@gochat.local', 'user1');
INSERT INTO auth_user (id, full_name, username, email, password)
    VALUES (4, 'User #2', 'user2', 'user2@gochat.local', 'user2');

INSERT INTO room (id, name, create_date) VALUES (1, 'Room #1', CURRENT_TIMESTAMP);
INSERT INTO room (id, name, create_date) VALUES (2, 'Room #2', CURRENT_TIMESTAMP);

INSERT INTO db_version (id, version) VALUES (1, 1);
