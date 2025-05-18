-- Create "bundles" table
CREATE TABLE "public"."bundles" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "resources" jsonb NULL,
  "inventory_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_bundles_inventory" FOREIGN KEY ("inventory_id") REFERENCES "public"."inventories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
