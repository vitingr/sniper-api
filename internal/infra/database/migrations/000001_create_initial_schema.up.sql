-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  profile_picture VARCHAR(500) DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NULL,

  CONSTRAINT users_email_check CHECK (
    email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
  )
);

CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

COMMENT ON TABLE users IS 'Application users';
COMMENT ON COLUMN users.email IS 'Unique user email address';
COMMENT ON COLUMN users.password IS 'Hashed password (bcrypt)';

CREATE TABLE IF NOT EXISTS companies (
  id TEXT PRIMARY KEY,
  name VARCHAR(500) NOT NULL,
  address TEXT,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  business_status VARCHAR(50),
  types TEXT [],
  website VARCHAR(500),
  phone VARCHAR(50),
  summary TEXT,
  employee_count VARCHAR(50),
  revenue VARCHAR(100),
  industry VARCHAR(255),
  linkedin_url VARCHAR(500),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  CONSTRAINT companies_coordinates_check CHECK (
    latitude BETWEEN -90 AND 90 AND longitude BETWEEN -180 AND 180
  ),
  CONSTRAINT companies_website_check CHECK (
    website IS NULL OR website ~* '^https?://'
  )
);


CREATE INDEX idx_companies_name ON companies(name);
CREATE INDEX idx_companies_coordinates ON companies(latitude, longitude);
CREATE INDEX idx_companies_created_at ON companies(created_at DESC);
CREATE INDEX idx_companies_industry ON companies(industry) WHERE industry IS NOT NULL;

COMMENT ON TABLE companies IS 'Companies discovered from Google Places and enriched with third-party data';
COMMENT ON COLUMN companies.types IS 'Array of business types (e.g., restaurant, cafe)';
COMMENT ON COLUMN companies.business_status IS 'Operational status (e.g., OPERATIONAL, CLOSED_TEMPORARILY)';

CREATE TABLE IF NOT EXISTS contacts (
  id SERIAL PRIMARY KEY,
  company_id TEXT NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255),
  full_name VARCHAR(500) GENERATED ALWAYS AS (
    CASE
      WHEN last_name IS NOT NULL THEN first_name || ' ' || last_name
      ELSE first_name
    END
  ) STORED,
  title VARCHAR(255),
  seniority VARCHAR(100),
  department VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(50),
  direct_phone VARCHAR(50),
  linkedin_url VARCHAR(500),
  twitter_url VARCHAR(500),
  photo_url VARCHAR(500),
  email_status VARCHAR(50) DEFAULT 'unknown',
  source VARCHAR(100) NOT NULL DEFAULT 'manual',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_contacts_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,

  CONSTRAINT contacts_email_check CHECK (
    email IS NULL OR email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
  ),
  CONSTRAINT contacts_email_status_check CHECK (
    email_status IN (
      'verified',
      'risky',
      'unknown',
      'invalid',
      'catch_all'
    )
  ),
  CONSTRAINT contacts_source_check CHECK (
    source IN (
      'apollo',
      'clearbit',
      'hunter',
      'rocketreach',
      'manual',
      'linkedin',
      'import'
    )
  )
);

CREATE INDEX idx_contacts_company_id ON contacts(company_id);
CREATE INDEX idx_contacts_email ON contacts(email) WHERE email IS NOT NULL;
CREATE INDEX idx_contacts_full_name ON contacts(full_name);
CREATE INDEX idx_contacts_seniority ON contacts(seniority) WHERE seniority IS NOT NULL;
CREATE INDEX idx_contacts_department ON contacts(department) WHERE department IS NOT NULL;
CREATE INDEX idx_contacts_email_status ON contacts(email_status);
CREATE INDEX idx_contacts_source ON contacts(source);
CREATE INDEX idx_contacts_created_at ON contacts(created_at DESC);
CREATE INDEX idx_contacts_company_email_status ON contacts(company_id, email_status);

COMMENT ON TABLE contacts IS 'Individual contacts associated with companies';
COMMENT ON COLUMN contacts.full_name IS 'Computed column: first_name + last_name';
COMMENT ON COLUMN contacts.email_status IS 'Email verification status from enrichment services';
COMMENT ON COLUMN contacts.source IS 'Data source/enrichment provider';

-- TRIGERS:
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_companies_updated_at
  BEFORE UPDATE ON companies
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_contacts_updated_at
  BEFORE UPDATE ON contacts
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

ALTER TABLE users SET (autovacuum_vacuum_scale_factor = 0.1);
ALTER TABLE companies SET (autovacuum_vacuum_scale_factor = 0.1);
ALTER TABLE contacts SET (autovacuum_vacuum_scale_factor = 0.1);

ANALYZE users;
ANALYZE companies;
ANALYZE contacts;
