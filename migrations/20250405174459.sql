-- Create "system_connections" table
CREATE TABLE "public"."system_connections" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "system_from_id" text NULL,
  "system_to_id" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_systems_connections" FOREIGN KEY ("system_from_id") REFERENCES "public"."systems" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
