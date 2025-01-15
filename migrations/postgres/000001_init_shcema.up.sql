CREATE TABLE IF NOT EXISTS "crypto_exchange_rate" (
  "crypto_exchange_rate_id" SERIAL PRIMARY KEY,
  "crypto" text NOT NULL,
  "exchange_rate" text NOT NULL,
  "date" timestamp NOT NULL
);
