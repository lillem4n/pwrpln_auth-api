-- migrate:up

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY,
  "created" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "name" text NOT NULL,
  "apiKey" text,
  "password" text
);
CREATE UNIQUE INDEX idx_accountname ON accounts ("name");

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

CREATE TABLE "renewalTokens" (
  "accountId" uuid NOT NULL,
  "exp" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP + '24 hours',
  "token" char(60) NOT NULL
);
ALTER TABLE "renewalTokens"
  ADD FOREIGN KEY ("accountId") REFERENCES "accounts" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE INDEX idx_renewaltokensaccountid ON "renewalTokens" ("accountId");
CREATE INDEX idx_renewaltokensexp ON "renewalTokens" ("exp");
CREATE INDEX idx_renewaltokenstoken ON "renewalTokens" ("token");

-- migrate:down