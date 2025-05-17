-- Create "exploration_infos" table
CREATE TABLE "public"."exploration_infos" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "exploring" boolean NULL,
  "type" text NULL,
  "started_at" bigint NULL,
  "required_time" bigint NULL,
  "astral_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_exploration_infos_astral" FOREIGN KEY ("astral_id") REFERENCES "public"."astrals" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
