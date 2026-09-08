\connect persons

GRANT USAGE, CREATE ON SCHEMA public TO program;
SET ROLE program;

\ir schema/01-persons.sql

RESET ROLE;
