-- +goose Up
-- +goose StatementBegin
CREATE TABLE "roles" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "title" varchar NOT NULL,
  "description" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "venues" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "title" varchar NOT NULL,
  "slug" varchar NOT NULL,
  "venue_id" uuid,
  "logo_url" varchar,
  "status" boolean NOT NULL DEFAULT false,
  "description" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "role_id" uuid NOT NULL,
  "venue_id" uuid NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "phone_number" varchar NOT NULL,
  "middle_name" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "venues" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "roles" CASCADE;

DROP TABLE IF EXISTS "venues" CASCADE;

DROP TABLE IF EXISTS "users" CASCADE;

-- +goose StatementEnd
