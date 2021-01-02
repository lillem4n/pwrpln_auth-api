-- migrate:up

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY,
  "created" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "accountName" text NOT NULL,
  "apiKey" text,
  "password" text
);
CREATE UNIQUE INDEX idx_accountname ON accounts ("accountName");

CREATE TABLE "accountsFields" (
  "id" uuid PRIMARY KEY,
  "created" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "accountId" uuid NOT NULL,
  "name" text NOT NULL,
  "value" text[] NOT NULL
);
ALTER TABLE "accountsFields"
  ADD FOREIGN KEY ("accountId") REFERENCES "accounts" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_accountsfields ON "accountsFields" ("accountId", "name");

-- migrate:down