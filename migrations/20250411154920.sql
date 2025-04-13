-- Modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "in_spaceship", DROP COLUMN "location", DROP COLUMN "system_id";
-- Create "astrals" table
CREATE TABLE "public"."astrals" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "code" text NOT NULL,
  "in_spaceship" boolean NULL DEFAULT false,
  "location" text NULL,
  "system_id" text NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_astrals_code" UNIQUE ("code"),
  CONSTRAINT "fk_astrals_system" FOREIGN KEY ("system_id") REFERENCES "public"."systems" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_astrals_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Rename a column from "user_id" to "astral_id"
ALTER TABLE "public"."items" RENAME COLUMN "user_id" TO "astral_id";
-- Modify "items" table
ALTER TABLE "public"."items" DROP CONSTRAINT "fk_items_user", ADD CONSTRAINT "fk_items_astral" FOREIGN KEY ("astral_id") REFERENCES "public"."astrals" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Rename a column from "user_id" to "astral_id"
ALTER TABLE "public"."spaceships" RENAME COLUMN "user_id" TO "astral_id";
-- Modify "spaceships" table
ALTER TABLE "public"."spaceships" DROP CONSTRAINT "fk_users_spaceships", ADD CONSTRAINT "fk_astrals_spaceships" FOREIGN KEY ("astral_id") REFERENCES "public"."astrals" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "system_connections" table
ALTER TABLE "public"."system_connections" DROP CONSTRAINT "fk_systems_connections", ADD CONSTRAINT "fk_systems_connections" FOREIGN KEY ("system_from_id") REFERENCES "public"."systems" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
