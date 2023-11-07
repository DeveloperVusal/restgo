CREATE TABLE IF NOT EXISTS "users" (
  "id" SERIAL,
  "email" VARCHAR(150) NOT NULL,
  "password" VARCHAR(255)  NOT NULL,
  "activation" BOOL DEFAULT false,
  "name" VARCHAR(200),
  "surname" VARCHAR(200) ,
  "token_secret_key" text NOT NULL,
  "updated_at" TIMESTAMP(0),
  "created_at" TIMESTAMP(0) NOT NULL,
  "confirm_code" CHAR(6) DEFAULT NULL::bpchar,
  "confirmed_at" TIMESTAMP(0) DEFAULT NULL::TIMESTAMP without time zone,
  "confirm_status" VARCHAR(20) DEFAULT NULL::character varying,
  CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "users_email_key" UNIQUE ("email")
);

-- ALTER TABLE "users" OWNER TO "butago_accounts";

COMMENT ON COLUMN "users"."created_at" IS 'Время создания пользователя';
COMMENT ON COLUMN "users"."confirm_code" IS 'Код подтверждения';
COMMENT ON COLUMN "users"."confirmed_at" IS 'Время последнего подтверждения';
COMMENT ON COLUMN "users"."confirm_status" IS 'Статус кода подтверждения';