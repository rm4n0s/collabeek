CREATE TYPE roles AS ENUM ('admin', 'moderator', 'member');

create table
    if not exists member (
        id SERIAL PRIMARY KEY,
        username VARCHAR(60) NOT NULL DEFAULT '',
        fullname VARCHAR(60) NOT NULL DEFAULT '',
        password VARCHAR(73) NOT NULL DEFAULT '',
        email VARCHAR(60) NOT NULL,
        registration_secret VARCHAR(60) NOT NULL,
        email_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
        role roles NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at TIMESTAMP
    );

create table
    if not exists team (
        id SERIAL PRIMARY KEY,
        name VARCHAR(24) NOT NULL,
        description VARCHAR(200) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at TIMESTAMP
    );

create table
    if not exists team_member (
        id SERIAL PRIMARY KEY,
        member_id SERIAL NOT NULL,
        team_id SERIAL NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at TIMESTAMP
    );