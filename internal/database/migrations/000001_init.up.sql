-- Enums

CREATE TYPE user_role AS ENUM ('ADMIN', 'MANAGER', 'OBSERVER');
CREATE TYPE contractor_type AS ENUM ('INDIVIDUAL', 'IP', 'LLC');
CREATE TYPE job_status AS ENUM ('DRAFT', 'ACTIVE', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED');
CREATE TYPE dispatch_status AS ENUM ('PENDING', 'SENDING', 'COMPLETED', 'PAUSED', 'CANCELLED');
CREATE TYPE offer_status AS ENUM ('SENT', 'ACCEPTED', 'DECLINED', 'EXPIRED');
CREATE TYPE assignment_status AS ENUM ('ASSIGNED', 'CONFIRMED', 'IN_WORK', 'COMPLETED', 'CANCELLED');
CREATE TYPE communication_direction AS ENUM ('OUTGOING', 'INCOMING');
CREATE TYPE rating_author_type AS ENUM ('CUSTOMER', 'CONTRACTOR');

-- Tables

CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    inn TEXT UNIQUE,
    contact_email TEXT,
    contact_phone TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'MANAGER',
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE contractors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    phone TEXT NOT NULL UNIQUE,
    email TEXT,
    type contractor_type NOT NULL DEFAULT 'INDIVIDUAL',
    regions TEXT[] NOT NULL DEFAULT '{}',
    equipment_types TEXT[] NOT NULL DEFAULT '{}',
    price_expectations TEXT,
    is_available BOOLEAN NOT NULL DEFAULT true,
    consent_given BOOLEAN NOT NULL DEFAULT false,
    consent_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    region TEXT NOT NULL,
    equipment_type TEXT,
    volume TEXT,
    deadline TIMESTAMPTZ,
    price DECIMAL(12, 2),
    conditions TEXT,
    status job_status NOT NULL DEFAULT 'DRAFT',
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    created_by_user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE job_dispatches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    wave INT NOT NULL DEFAULT 1,
    status dispatch_status NOT NULL DEFAULT 'PENDING',
    "limit" INT NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE job_offers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    dispatch_id UUID NOT NULL REFERENCES job_dispatches(id) ON DELETE CASCADE,
    contractor_id UUID NOT NULL REFERENCES contractors(id),
    status offer_status NOT NULL DEFAULT 'SENT',
    token UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    responded_at TIMESTAMPTZ,
    UNIQUE(job_id, contractor_id)
);

CREATE TABLE assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    contractor_id UUID NOT NULL REFERENCES contractors(id),
    status assignment_status NOT NULL DEFAULT 'ASSIGNED',
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(job_id, contractor_id)
);

CREATE TABLE communication_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID REFERENCES jobs(id) ON DELETE SET NULL,
    contractor_id UUID NOT NULL REFERENCES contractors(id),
    channel TEXT NOT NULL,
    message TEXT NOT NULL,
    direction communication_direction NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE contractor_blacklist (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    contractor_id UUID NOT NULL REFERENCES contractors(id),
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(organization_id, contractor_id)
);

CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    contractor_id UUID NOT NULL REFERENCES contractors(id),
    author_contractor_id UUID REFERENCES contractors(id),
    organization_id UUID,
    score INT NOT NULL,
    comment TEXT,
    author_type rating_author_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes

CREATE INDEX idx_users_organization_id ON users(organization_id);
CREATE INDEX idx_jobs_organization_id ON jobs(organization_id);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_job_dispatches_job_id ON job_dispatches(job_id);
CREATE INDEX idx_job_offers_job_id ON job_offers(job_id);
CREATE INDEX idx_job_offers_contractor_id ON job_offers(contractor_id);
CREATE INDEX idx_job_offers_token ON job_offers(token);
CREATE INDEX idx_assignments_job_id ON assignments(job_id);
CREATE INDEX idx_assignments_contractor_id ON assignments(contractor_id);
CREATE INDEX idx_communication_logs_contractor_id ON communication_logs(contractor_id);
CREATE INDEX idx_contractor_blacklist_organization_id ON contractor_blacklist(organization_id);
CREATE INDEX idx_ratings_contractor_id ON ratings(contractor_id);
CREATE INDEX idx_ratings_job_id ON ratings(job_id);

-- Updated_at trigger function

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_organizations BEFORE UPDATE ON organizations FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_users BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_contractors BEFORE UPDATE ON contractors FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_jobs BEFORE UPDATE ON jobs FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_job_dispatches BEFORE UPDATE ON job_dispatches FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER set_updated_at_assignments BEFORE UPDATE ON assignments FOR EACH ROW EXECUTE FUNCTION update_updated_at();
