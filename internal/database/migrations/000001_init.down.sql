DROP TRIGGER IF EXISTS set_updated_at_assignments ON assignments;
DROP TRIGGER IF EXISTS set_updated_at_job_dispatches ON job_dispatches;
DROP TRIGGER IF EXISTS set_updated_at_jobs ON jobs;
DROP TRIGGER IF EXISTS set_updated_at_contractors ON contractors;
DROP TRIGGER IF EXISTS set_updated_at_users ON users;
DROP TRIGGER IF EXISTS set_updated_at_organizations ON organizations;
DROP FUNCTION IF EXISTS update_updated_at();

DROP TABLE IF EXISTS ratings;
DROP TABLE IF EXISTS contractor_blacklist;
DROP TABLE IF EXISTS communication_logs;
DROP TABLE IF EXISTS assignments;
DROP TABLE IF EXISTS job_offers;
DROP TABLE IF EXISTS job_dispatches;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS contractors;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS organizations;

DROP TYPE IF EXISTS rating_author_type;
DROP TYPE IF EXISTS communication_direction;
DROP TYPE IF EXISTS assignment_status;
DROP TYPE IF EXISTS offer_status;
DROP TYPE IF EXISTS dispatch_status;
DROP TYPE IF EXISTS job_status;
DROP TYPE IF EXISTS contractor_type;
DROP TYPE IF EXISTS user_role;
