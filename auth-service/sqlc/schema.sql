CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username text NOT NULL UNIQUE,   
    email text NOT NULL UNIQUE,
    encoded_password text NOT NULL
);
