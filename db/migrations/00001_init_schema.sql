-- +goose Up
-- +goose StatementBegin
CREATE TABLE "roles" (
  "id" bigserial PRIMARY KEY,
  "uuid" uuid DEFAULT (gen_random_uuid()),
  "title" varchar NOT NULL,
  "description" text,
  "created_at" timestamptz DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "restaurants" (
  "id" bigserial PRIMARY KEY,
  "uuid" uuid DEFAULT (gen_random_uuid()),
  "title" varchar NOT NULL,
  "slug" varchar NOT NULL,
  "logo_url" varchar,
  "status" BOOLEAN NOT NULL DEFAULT false,
  "description" text,
  "created_at" timestamptz DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "uuid" uuid DEFAULT (gen_random_uuid()),
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "role_id" INT NOT NULL,
  "restaurant_id" INT NOT NULL,
  "full_name" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "deleted_at" timestamptz
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("restaurant_id") REFERENCES "restaurants" ("id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "roles" CASCADE;

DROP TABLE IF EXISTS "restaurants" CASCADE;

DROP TABLE IF EXISTS "users" CASCADE;

-- +goose StatementEnd
