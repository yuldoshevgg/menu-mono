-- +goose Up
-- +goose StatementBegin
CREATE TABLE "roles" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "title" varchar NOT NULL,
  "description" text,
  "created_at" timestamptz DEFAULT (now()),
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
  "created_at" timestamptz DEFAULT (now()),
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
  "created_at" timestamptz DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "categories" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" varchar NOT NULL,
  "venue_id" uuid NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "menu_items" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "title" varchar NOT NULL,
  "category_id" uuid NOT NULL,
  "venue_id" uuid NOT NULL,
  "price" decimal NOT NULL DEFAULT 0,
  "image_url" varchar NOT NULL DEFAULT '',
  "is_available" boolean NOT NULL DEFAULT false,
  "description" text,
  "created_at" timestamptz DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "qr_links" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "venue_id" uuid NOT NULL,
  "table_number" integer NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "deleted_at" timestamptz
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "categories" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "menu_items" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "menu_items" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "qr_links" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

CREATE UNIQUE INDEX ON "roles" ("title");

CREATE UNIQUE INDEX ON "venues" ("slug");

CREATE UNIQUE INDEX ON "users" ("username");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "roles" CASCADE;

DROP TABLE IF EXISTS "venues" CASCADE;

DROP TABLE IF EXISTS "users" CASCADE;

DROP TABLE IF EXISTS "categories" CASCADE;

DROP TABLE IF EXISTS "menu_items" CASCADE;

DROP INDEX IF EXISTS "roles_title_idx";

DROP INDEX IF EXISTS "venues_slug_idx";

DROP INDEX IF EXISTS "users_username_idx";

-- +goose StatementEnd
