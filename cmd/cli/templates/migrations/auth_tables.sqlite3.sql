DROP TABLE IF EXISTS "users";

CREATE TABLE "users" (
	"id"	TEXT NOT NULL,
	"first_name"	TEXT NOT NULL,
	"last_name"	TEXT NOT NULL,
	"user_active"	INTEGER NOT NULL DEFAULT '0',
	"access_level"	INTEGER NOT NULL DEFAULT '3',
	"email"	TEXT NOT NULL UNIQUE,
	"password"	TEXT NOT NULL,
	"deleted_at"	DATETIME,
	"created_at"	DATETIME DEFAULT current_timestamp,
	"updated_at"	DATETIME DEFAULT current_timestamp,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "remember_tokens";

CREATE TABLE "remember_tokens" (
	"id"	TEXT NOT NULL,
	"user_id"	TEXT NOT NULL,
	"remember_token"	TEXT NOT NULL,
	"created_at"	DATETIME DEFAULT current_timestamp,
	"updated_at"	DATETIME DEFAULT current_timestamp,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "tokens";

CREATE TABLE "tokens" (
	"id"	TEXT NOT NULL,
	"user_id"	TEXT NOT NULL,
	"email"	TEXT NOT NULL,
	"token"	TEXT NOT NULL,
	"token_hash"	BLOB,
	"created_at"	DATETIME DEFAULT current_timestamp,
	"updated_at"	DATETIME DEFAULT current_timestamp,
	"expiry"	DATETIME,
	PRIMARY KEY("id")
);
