CREATE TABLE IF NOT EXISTS user_accounts(
  id SERIAL PRIMARY KEY,
  number TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS objects(
  id SERIAL PRIMARY KEY,
  code TEXT NOT NULL UNIQUE,
  name TEXT,
  delivered BOOLEAN NOT NULL DEFAULT false,
  last_sent_hash TEXT DEFAULT '',
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS events(
   id SERIAL PRIMARY KEY,
   hash TEXT NOT NULL UNIQUE,
   body json NOT NULL,
   object_id SERIAL NOT NULL UNIQUE,
   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP NOT NULL,
   FOREIGN KEY (object_id) REFERENCES objects (id)
);

CREATE TABLE IF NOT EXISTS object_user_account(
  id SERIAL PRIMARY KEY,
  object_id SERIAL NOT NULL,
  user_account_id SERIAL NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  FOREIGN KEY (object_id) REFERENCES objects (id),
  FOREIGN KEY (user_account_id) REFERENCES user_accounts (id)
);

CREATE TABLE IF NOT EXISTS logs(
    id SERIAL PRIMARY KEY,
    log TEXT,
    env TEXT,
    date TIMESTAMP,
    label TEXT
);