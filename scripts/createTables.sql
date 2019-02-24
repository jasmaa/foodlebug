CREATE TABLE users (
    id int NOT NULL UNIQUE,
    username varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    rating real NOT NULL
);

CREATE TABLE posts (
    id int NOT NULL UNIQUE,
    posterId int NOT NULL,
    photo bytea NOT NULL,
    content varchar(255),
    timePosted abstime,
    locationLat real,
    locationLon real,
    /*commentIds integer[]*/
    visible boolean,
);

CREATE TABLE comments (
    id int NOT NULL UNIQUE,
    postId int NOT NULL,
    posterId int NOT NULL,
    content varchar(255),
    visible boolean,
);
