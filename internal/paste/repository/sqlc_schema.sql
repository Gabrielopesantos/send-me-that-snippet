CREATE TABLE "pastes" (
    "id" text NOT NULL,
    "content" text NOT NULL,
    "content_sha" text NOT NULL,
    "language" text NOT NULL,
    "created_at" timestamptz,
    "expires_in" bigint NOT NULL,
    "expired" boolean NOT NULL,
    CONSTRAINT "pastes_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
