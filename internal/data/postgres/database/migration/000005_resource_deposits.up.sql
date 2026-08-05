CREATE TYPE location_type AS ENUM ('PLANET', 'WAYPOINT');

CREATE TABLE IF NOT EXISTS resource_deposit_states (
    system_x INT NOT NULL,
    system_y INT NOT NULL,
    loc_type location_type NOT NULL,
    loc_id INT NOT NULL,
    resource_type VARCHAR NOT NULL,

    remaining BIGINT NOT NULL CHECK (remaining >= 0),
    last_mined_at TIMESTAMP NOT NULL DEFAULT NOW(), 
    PRIMARY KEY (system_x, system_y, loc_type, loc_id, resource_type)
);