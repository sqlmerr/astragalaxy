-- Create "inventories" table
CREATE TABLE "public"."inventories" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "holder" text NOT NULL,
  "holder_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Rename a column from "astral_id" to "inventory_id"
ALTER TABLE "public"."items" RENAME COLUMN "astral_id" TO "inventory_id";
-- Modify "items" table
ALTER TABLE "public"."items" DROP CONSTRAINT "fk_items_astral", ADD CONSTRAINT "fk_items_inventory" FOREIGN KEY ("inventory_id") REFERENCES "public"."inventories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
