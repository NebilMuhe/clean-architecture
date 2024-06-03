CREATE TABLE "sessions" (
     "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     "username" VARCHAR (50) NOT NULL UNIQUE,
     "refresh_token" VARCHAR (500) NOT NULL,
     "created_at" timestamptz NOT NULL default now(),
     "deleted_at" timestamptz
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE;