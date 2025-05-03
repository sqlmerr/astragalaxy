-- Create "wallets" table
CREATE TABLE "public"."wallets" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "units" bigint NOT NULL,
  "quarks" bigint NOT NULL,
  "locked" boolean NOT NULL DEFAULT false,
  "astral_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_wallets_astral" FOREIGN KEY ("astral_id") REFERENCES "public"."astrals" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
