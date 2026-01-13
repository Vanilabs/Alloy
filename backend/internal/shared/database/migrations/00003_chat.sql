-- +goose Up
-- Channels table
CREATE TABLE IF NOT EXISTS channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(10) NOT NULL,
    name VARCHAR(100),
    description VARCHAR(400),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Channel members table
CREATE TABLE IF NOT EXISTS channel_members (
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (channel_id, user_id)
);

-- Conversations table
CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Conversation reads table
CREATE TABLE IF NOT EXISTS conversation_reads (
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_message_id UUID NOT NULL,
    last_read_at TIMESTAMP NOT NULL,
    last_delivered_at TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (conversation_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_conversation_reads_user_last_read ON conversation_reads(user_id, last_read_at DESC);
CREATE INDEX IF NOT EXISTS idx_channel_members_user ON channel_members(user_id);

-- +goose Down

DROP TABLE IF EXISTS conversation_reads;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS channel_members;
DROP TABLE IF EXISTS channels;
