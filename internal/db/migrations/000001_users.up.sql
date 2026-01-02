create extension if not exists "pgcrypto";

CREATE TABLE if NOT EXISTS users (
  uuid uuid NOT NULL DEFAULT gen_random_uuid(),
  name varchar(50) NOT NULL,
  email varchar(100) UNIQUE NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
