CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username text NOT NULL UNIQUE,   
    email text NOT NULL UNIQUE,
    first_name text NOT NULL,
    last_name text NOT NULL,
    dob text NOT NULL
);