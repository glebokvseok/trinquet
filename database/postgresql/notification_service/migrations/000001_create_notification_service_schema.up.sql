CREATE SCHEMA IF NOT EXISTS notification_service;
GRANT USAGE ON SCHEMA notification_service to "service-user";
ALTER DEFAULT PRIVILEGES IN SCHEMA notification_service GRANT SELECT, UPDATE, INSERT, DELETE ON TABLES to "service-user";
