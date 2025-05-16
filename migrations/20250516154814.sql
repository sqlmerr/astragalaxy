-- Modify "system_connections" table
ALTER TABLE "public"."system_connections" DROP CONSTRAINT "fk_systems_connections", ADD CONSTRAINT "fk_systems_connections" FOREIGN KEY ("system_from_id") REFERENCES "public"."systems" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
