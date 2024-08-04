-- message table
CREATE TABLE
    IF NOT EXISTS "message" (
        "id" BIGSERIAL PRIMARY KEY,
        "from_user_id" BIGINT NOT NULL REFERENCES "users" ("id"),
        "to_user_id" BIGINT NOT NULL REFERENCES "users" ("id"),
        "is_sent" BOOLEAN NOT NULL,
        "content" TEXT NOT NULL,
        "created_at" TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_message_from_to ON message (from_user_id, to_user_id);