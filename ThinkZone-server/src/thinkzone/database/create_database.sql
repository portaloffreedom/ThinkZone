CREATE TABLE t_user (
  id INT PRIMARY KEY,
  username CHAR(64) NOT NULL UNIQUE,
  password BYTEA
);

CREATE TABLE conversation (
  id INT PRIMARY KEY
);

CREATE TABLE archive (
  t_user INT,
  conversation INT,
  
  PRIMARY KEY (t_user,conversation),
  FOREIGN KEY (t_user) REFERENCES t_user(id),
  FOREIGN KEY (conversation) REFERENCES conversation(id)
);

CREATE TABLE post (
  id INT,
  conversation INT,
  text TEXT,
  
  father INT,
  first_response INT,
  second_response INT,
  
  PRIMARY KEY (id, conversation),
  FOREIGN KEY (conversation)    REFERENCES conversation(id)
--   FOREIGN KEY (pather)          REFERENCES post(id),
--   FOREIGN KEY (first_response)  REFERENCES post(id),
--   FOREIGN KEY (second_response) REFERENCES post(id)
);

CREATE TABLE author (
  t_user INT,
  post INT,
  conversation INT,
  
  PRIMARY KEY (t_user,conversation, post),
  FOREIGN KEY (t_user)         REFERENCES t_user(id),
  FOREIGN KEY (post, conversation)         REFERENCES post(id, conversation)
);