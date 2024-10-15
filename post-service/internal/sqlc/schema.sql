CREATE TABLE posts
(
    id SERIAL PRIMARY KEY,
    owner_id int NOT NULL,   
    body text NOT NULL,
    media_id text
);
