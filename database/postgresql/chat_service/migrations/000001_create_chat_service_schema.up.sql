CREATE SCHEMA IF NOT EXISTS chat_service;
GRANT USAGE ON SCHEMA chat_service to "service-user";
ALTER DEFAULT PRIVILEGES IN SCHEMA chat_service GRANT SELECT, UPDATE, INSERT, DELETE ON TABLES to "service-user";
