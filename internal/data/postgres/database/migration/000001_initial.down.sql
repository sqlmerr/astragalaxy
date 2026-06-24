DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS agents;

DROP INDEX IF EXISTS idx_agents_token_hash;

DROP TYPE IF EXISTS ship_type;

DROP TYPE IF EXISTS ship_status;

DROP TABLE IF EXISTS ships;

DROP INDEX IF EXISTS idx_ships_agent_id;

DROP EXTENSION IF EXISTS "uuid-ossp";
