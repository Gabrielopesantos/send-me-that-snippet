CREATE TABLE "pastes" (
    "id" text NOT NULL,
    "content" text,
    "content_sha" text,
    "language" text,
    "created_at" timestamptz,
    "expires_in" bigint,
    "expired" boolean,
    CONSTRAINT "pastes_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
