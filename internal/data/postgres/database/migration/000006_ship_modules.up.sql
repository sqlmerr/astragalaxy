CREATE TABLE IF NOT EXISTS ship_modules (
    ship_id UUID NOT NULL REFERENCES ships(id) ON DELETE CASCADE,
    module_type VARCHAR(255) NOT NULL,
    
    PRIMARY KEY (ship_id, module_type)
);