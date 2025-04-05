-- Create "systems" table
CREATE TABLE "public"."systems" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "username" text NOT NULL,
  "password" text NULL,
  "in_spaceship" boolean NULL DEFAULT false,
  "location" text NULL,
  "system_id" uuid NULL,
  "token" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_username" UNIQUE ("username"),
  CONSTRAINT "fk_users_system" FOREIGN KEY ("system_id") REFERENCES "public"."systems" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "items" table
CREATE TABLE "public"."items" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NULL,
  "code" text NULL,
  "durability" bigint NULL DEFAULT 100,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_items_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "item_data_tags" table
CREATE TABLE "public"."item_data_tags" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "item_id" uuid NULL,
  "key" text NOT NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_item_data_tags_item" FOREIGN KEY ("item_id") REFERENCES "public"."items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "planets" table
CREATE TABLE "public"."planets" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NULL,
  "system_id" uuid NULL,
  "threat" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_planets_system" FOREIGN KEY ("system_id") REFERENCES "public"."systems" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "flight_infos" table
CREATE TABLE "public"."flight_infos" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "flying" boolean NOT NULL DEFAULT false,
  "destination" text NULL,
  "destination_id" text NULL,
  "flown_out_at" bigint NULL,
  "flying_time" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create "spaceships" table
CREATE TABLE "public"."spaceships" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL,
  "user_id" uuid NULL,
  "location" text NULL,
  "flight_id" uuid NULL,
  "system_id" uuid NULL,
  "planet_id" uuid NULL,
  "player_sit_in" boolean NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_spaceships_flight" FOREIGN KEY ("flight_id") REFERENCES "public"."flight_infos" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_spaceships_planet" FOREIGN KEY ("planet_id") REFERENCES "public"."planets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_spaceships_system" FOREIGN KEY ("system_id") REFERENCES "public"."systems" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_users_spaceships" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
