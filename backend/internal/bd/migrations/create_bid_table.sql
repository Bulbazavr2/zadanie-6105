CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- Уникальный идентификатор предложения
    name VARCHAR(100) NOT NULL, -- Полное название предложения
    description VARCHAR(500) NOT NULL, -- Описание предложения
    status VARCHAR(50) NOT NULL, -- Статус предложения
    tender_id UUID NOT NULL, -- Уникальный идентификатор тендера
    author_type VARCHAR(50) NOT NULL, -- Тип автора (Enum возможных значений)
    author_id UUID NOT NULL, -- Уникальный идентификатор автора предложения
    version INTEGER NOT NULL DEFAULT 1 CHECK (version >= 1), -- Номер версии предложения
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Дата и время создания предложения
    FOREIGN KEY (tender_id) REFERENCES tenders(id) ON DELETE CASCADE, -- Связь с таблицей тендеров
    FOREIGN KEY (author_id) REFERENCES employee(id) ON DELETE CASCADE -- Связь с таблицей сотрудников
);

CREATE TABLE IF NOT EXISTS bid_versions (
    version_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    status VARCHAR(50) NOT NULL,
    tender_id UUID NOT NULL,
    author_type VARCHAR(50) NOT NULL,
    author_id UUID NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bid_id) REFERENCES bid(id) ON DELETE CASCADE
);

-- Удаляем триггер, если он существует
DROP TRIGGER IF EXISTS trigger_save_bid_version ON bid;

-- Удаляем функцию, если она существует
DROP FUNCTION IF EXISTS save_bid_version();

CREATE OR REPLACE FUNCTION save_bid_version() RETURNS TRIGGER AS $$
BEGIN
    -- Сохранение текущей версии предложения в таблицу bid_versions перед обновлением
    INSERT INTO bid_versions (bid_id, name, description, status, tender_id, author_type, author_id, version, created_at)
    SELECT OLD.id, OLD.name, OLD.description, OLD.status, OLD.tender_id, OLD.author_type, OLD.author_id, OLD.version, OLD.created_at;
    -- Увеличиваем версию на 1 при каждом обновлении
    NEW.version := OLD.version + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_save_bid_version
BEFORE UPDATE ON bid
FOR EACH ROW
EXECUTE FUNCTION save_bid_version();

