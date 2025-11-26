CREATE TYPE message_status AS ENUM ('pending', 'sent', 'failed');

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    "to" VARCHAR(20) NOT NULL,
    content VARCHAR(160) NOT NULL,
    status message_status NOT NULL DEFAULT 'pending',
    sent_at TIMESTAMPTZ NULL,
    external_id VARCHAR(100) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_status ON messages(status);