BEGIN;

-- Create Company table
CREATE TABLE IF NOT EXISTS company (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP   
);

-- Create Contractor table
CREATE TABLE IF NOT EXISTS contractor (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

-- Create Job table
CREATE TABLE IF NOT EXISTS job (
    id VARCHAR(255) NOT NULL,
    version INTEGER NOT NULL,
    uid VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    rate DECIMAL(10, 2) NOT NULL,
    title VARCHAR(255) NOT NULL,
    company_id VARCHAR(255) NOT NULL,
    contractor_id VARCHAR(255) NOT NULL,
    is_latest BOOLEAN NOT NULL,
    PRIMARY KEY (id, version)
);

-- Create index on uid for faster lookup
CREATE INDEX IF NOT EXISTS job_uid_idx ON job(uid);

-- Create Timelog table
CREATE TABLE IF NOT EXISTS timelog (
    id VARCHAR(255) NOT NULL,
    version INTEGER NOT NULL,
    uid VARCHAR(255) NOT NULL,
    duration BIGINT NOT NULL,
    time_start BIGINT NOT NULL,
    time_end BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,
    job_uid VARCHAR(255) NOT NULL,
    is_latest BOOLEAN NOT NULL,
    PRIMARY KEY (id, version)
);

-- Create index on uid for faster lookup
CREATE INDEX IF NOT EXISTS timelog_uid_idx ON timelog(uid);
CREATE INDEX IF NOT EXISTS timelog_job_uid_idx ON timelog(job_uid);

-- Create PaymentLineItems table
CREATE TABLE IF NOT EXISTS payment_line_items (
    id VARCHAR(255) NOT NULL,
    version INTEGER NOT NULL,
    uid VARCHAR(255) NOT NULL,
    job_uid VARCHAR(255) NOT NULL,
    timelog_uid VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    is_latest BOOLEAN NOT NULL,
    PRIMARY KEY (id, version)
);

-- Create index on uid for faster lookup
CREATE INDEX IF NOT EXISTS payment_line_items_uid_idx ON payment_line_items(uid);
CREATE INDEX IF NOT EXISTS payment_line_items_job_uid_idx ON payment_line_items(job_uid);
CREATE INDEX IF NOT EXISTS payment_line_items_timelog_uid_idx ON payment_line_items(timelog_uid);

COMMIT;
