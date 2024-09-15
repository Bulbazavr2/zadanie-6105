CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tenders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    organization_id UUID NOT NULL,
    creator_username VARCHAR(50) NOT NULL,
    service_type VARCHAR(100) NOT NULL,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    FOREIGN KEY (creator_username) REFERENCES employee(username) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tender_versions (
    version_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tender_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    organization_id UUID NOT NULL,
    creator_username VARCHAR(50) NOT NULL,
    service_type VARCHAR(100) NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tender_id) REFERENCES tenders(id) ON DELETE CASCADE
);


-- Удаляем триггер, если он существует
DROP TRIGGER IF EXISTS trigger_save_tender_version ON tenders;

-- Удаляем функцию, если она существует
DROP FUNCTION IF EXISTS save_tender_version();

CREATE OR REPLACE FUNCTION save_tender_version() RETURNS TRIGGER AS $$
BEGIN
    -- Сохранение текущей версии тендера в таблицу tender_versions перед обновлением
    INSERT INTO tender_versions (tender_id, name, description, status, organization_id, creator_username, service_type, version, created_at)
    SELECT OLD.id, OLD.name, OLD.description, OLD.status, OLD.organization_id, OLD.creator_username, OLD.service_type, OLD.version, OLD.created_at;
    -- Увеличиваем версию на 1 при каждом обновлении
    NEW.version := OLD.version + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER  trigger_save_tender_version
BEFORE UPDATE ON tenders
FOR EACH ROW
EXECUTE FUNCTION save_tender_version();
