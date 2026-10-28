CREATE TABLE IF NOT EXISTS planet_discoveries (
    system_x INT NOT NULL,
    system_y INT NOT NULL,
    orbit INT NOT NULL,
    discovered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    discovered_by UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    PRIMARY KEY (system_x, system_y, orbit)
);

CREATE TABLE IF NOT EXISTS system_discoveries (
    system_x INT NOT NULL,
    system_y INT NOT NULL,
    discovered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    discovered_by UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    PRIMARY KEY (system_x, system_y)
);