-- Modify "planets" table
ALTER TABLE "public"."planets" ALTER COLUMN "name" SET DEFAULT 'undefined', ALTER COLUMN "system_id" SET NOT NULL;
