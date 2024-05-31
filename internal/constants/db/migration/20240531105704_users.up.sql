CREATE TABLE users (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "username" VARCHAR (50) UNIQUE NOT NULL,
    "email" VARCHAR (50) UNIQUE NOT NULL,
    "password" VARCHAR (200) NOT NULL,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now(),
    "deleted_at" timestamptz
);
