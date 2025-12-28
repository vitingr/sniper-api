-- DROP TRIGGERS
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_companies_updated_at ON companies;
DROP TRIGGER IF EXISTS update_contacts_updated_at ON contacts;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- DROP TABLES (in reverse order due to foreign keys)
DROP TABLE IF EXISTS contacts CASCADE;
DROP TABLE IF EXISTS companies CASCADE;
DROP TABLE IF EXISTS users CASCADE;
