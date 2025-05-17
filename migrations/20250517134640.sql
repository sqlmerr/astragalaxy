-- Modify "spaceships" table
ALTER TABLE "public"."spaceships" DROP CONSTRAINT "fk_spaceships_planet", ADD CONSTRAINT "fk_spaceships_planet" FOREIGN KEY ("planet_id") REFERENCES "public"."planets" ("id") ON UPDATE CASCADE ON DELETE SET NULL;
