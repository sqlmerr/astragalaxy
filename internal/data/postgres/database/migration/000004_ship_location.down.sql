ALTER TABLE ships DROP COLUMN location;
ALTER TABLE ships DROP COLUMN location_id;

DROP TYPE IF EXISTS ship_location;