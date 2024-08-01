-- users table
CREATE TABLE
    IF NOT EXISTS "users" (
        "id" BIGSERIAL PRIMARY KEY,
        "email" VARCHAR(255) NOT NULL UNIQUE,
        "username" VARCHAR(255) NOT NULL UNIQUE,
        "password" VARCHAR(255) NOT NULL,
        "users_photo_link" TEXT,
        "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE UNIQUE INDEX "user_email_index" ON "users" ("email");

CREATE UNIQUE INDEX "username_index" ON "users" ("username");