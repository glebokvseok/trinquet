CREATE SCHEMA IF NOT EXISTS player_service;
GRANT USAGE ON SCHEMA player_service to "service-user";
ALTER DEFAULT PRIVILEGES IN SCHEMA player_service GRANT SELECT, UPDATE, INSERT, DELETE ON TABLES to "service-user";
