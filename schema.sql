--create types
DO $$ BEGIN
    CREATE TYPE ROLE AS ENUM ('admin', 'member');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS user_accounts(
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  email TEXT,
  name TEXT NOT NULL,
  city TEXT NOT NULL,
  pix_key TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS buckets(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  due_date DATE NOT NULL,
  monthly_total_value TEXT NOT NULL,
  monthly_member_value TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS bucket_user_account(
  id SERIAL PRIMARY KEY,
  role ROLE NOT NULL,
  bucket_id SERIAL NOT NULL,
  user_account_id SERIAL NOT NULL,
  FOREIGN KEY (bucket_id) REFERENCES buckets (id),
  FOREIGN KEY (user_account_id) REFERENCES user_accounts (id)
);

CREATE TABLE IF NOT EXISTS transactions(
  id SERIAL PRIMARY KEY,
  attachment TEXT NOT NULL,
  bucket_id SERIAL NOT NULL,
  date DATE NOT NULL,
  user_account_id SERIAL NOT NULL,
  FOREIGN KEY (bucket_id) REFERENCES buckets (id),
  FOREIGN KEY (user_account_id) REFERENCES user_accounts (id)
);

CREATE TABLE IF NOT EXISTS logs (
    id SERIAL PRIMARY KEY,
    log TEXT NOT NULL,
    env TEXT,
    date TIMESTAMP,
    label TEXT
);
