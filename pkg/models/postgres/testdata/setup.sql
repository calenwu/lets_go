CREATE TABLE snippets (
  id SERIAL NOT NULL PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created TIMESTAMP(0) NOT NULL,
  expires TIMESTAMP(0) NOT NULL
);

CREATE TABLE users (
   id SERIAL NOT NULL PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   email VARCHAR(255) NOT NULL,
   hashed_password CHAR(60) NOT NULL,
   created TIMESTAMP(0) NOT NULL
);

ALTER TABLE users ADD UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created) VALUES (
 'Alice Jones',
 'alice@example.com',
 '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
 '2018-12-23 17:25:22'
 );
