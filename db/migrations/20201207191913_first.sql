-- migrate:up

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "created" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "username" text NOT NULL,
  "password" text NOT NULL
);

CREATE TABLE "usersFields" (
  "id" uuid PRIMARY KEY,
  "created" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "userId" uuid NOT NULL,
  "name" text NOT NULL,
  "value" text[] NOT NULL
);
ALTER TABLE "usersFields"
  ADD FOREIGN KEY ("userId") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_usersfields ON "usersFields" ("userId", "name");

-- migrate:down