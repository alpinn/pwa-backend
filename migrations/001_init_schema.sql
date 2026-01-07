CREATE SCHEMA "public";

CREATE TABLE "users" (
	"id" varchar(36) PRIMARY KEY,
	"username" varchar(50) NOT NULL CONSTRAINT "users_username_key" UNIQUE,
	"password" varchar(255) NOT NULL,
	"name" varchar(100) NOT NULL,
	"role" varchar(20) DEFAULT 'user' NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);

CREATE TABLE "products" (
	"id" varchar(36) PRIMARY KEY,
	"name" varchar(100) NOT NULL,
	"description" text,
	"price" numeric(10, 2) NOT NULL,
	"stock" integer DEFAULT 0 NOT NULL,
	"image_url" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);

CREATE TABLE "transactions" (
	"id" varchar(36) PRIMARY KEY,
	"user_id" varchar(36) NOT NULL,
	"total_amount" numeric(10, 2) NOT NULL,
	"status" varchar(20) DEFAULT 'pending' NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "fk_transactions_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);

CREATE TABLE "user_sessions" (
	"id" varchar(36) PRIMARY KEY,
	"user_id" varchar(36) NOT NULL,
	"refresh_token" varchar(255) NOT NULL,
	"device_id" varchar(100),
	"user_agent" text,
	"expires_at" timestamp NOT NULL,
	"revoked_at" timestamp,
	"created_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "fk_user_sessions_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);

CREATE TABLE "transaction_items" (
	"id" varchar(36) PRIMARY KEY,
	"transaction_id" varchar(36) NOT NULL,
	"product_id" varchar(36) NOT NULL,
	"product_name" varchar(100) NOT NULL,
	"quantity" integer NOT NULL,
	"price" numeric(10, 2) NOT NULL,
	"subtotal" numeric(10, 2) NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"user_id" varchar(36),
	CONSTRAINT "fk_transaction_items_transaction" FOREIGN KEY ("transaction_id") REFERENCES "transactions"("id") ON DELETE CASCADE,
	CONSTRAINT "fk_transaction_items_product" FOREIGN KEY ("product_id") REFERENCES "products"("id") ON DELETE RESTRICT,
	CONSTRAINT "fk_transaction_items_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE SET NULL
);

CREATE TABLE "stock_events" (
	"id" varchar(36) PRIMARY KEY,
	"product_id" varchar(36) NOT NULL,
	"qty" integer NOT NULL,
	"type" varchar(20) NOT NULL,
	"source" varchar(20) NOT NULL,
	"transaction_id" varchar(36),
	"user_id" varchar(36),
	"device_id" varchar(100),
	"note" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "fk_stock_events_product" FOREIGN KEY ("product_id") REFERENCES "products"("id") ON DELETE CASCADE,
	CONSTRAINT "fk_stock_events_transaction" FOREIGN KEY ("transaction_id") REFERENCES "transactions"("id") ON DELETE SET NULL,
	CONSTRAINT "fk_stock_events_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE SET NULL,
	CONSTRAINT "stock_events_source_check" CHECK (source IN ('pos', 'dashboard', 'online')),
	CONSTRAINT "stock_events_type_check" CHECK (type IN ('sale', 'restock', 'reject', 'adjustment', 'opening_stock'))
);

-- Users indexes
CREATE UNIQUE INDEX "users_pkey" ON "users" ("id");
CREATE UNIQUE INDEX "users_username_key" ON "users" ("username");

-- Products indexes
CREATE UNIQUE INDEX "products_pkey" ON "products" ("id");
CREATE INDEX "idx_products_name" ON "products" ("name");

-- Transactions indexes
CREATE UNIQUE INDEX "transactions_pkey" ON "transactions" ("id");
CREATE INDEX "idx_transactions_user_id" ON "transactions" ("user_id");
CREATE INDEX "idx_transactions_status" ON "transactions" ("status");

-- User Sessions indexes
CREATE INDEX "idx_user_sessions_user_id" ON "user_sessions" ("user_id");
CREATE INDEX "idx_user_sessions_refresh_token" ON "user_sessions" ("refresh_token");

-- Transaction Items indexes
CREATE UNIQUE INDEX "transaction_items_pkey" ON "transaction_items" ("id");
CREATE INDEX "idx_transaction_items_transaction_id" ON "transaction_items" ("transaction_id");
CREATE INDEX "idx_transaction_items_user_id" ON "transaction_items" ("user_id");

-- Stock Events indexes
CREATE UNIQUE INDEX "stock_events_pkey" ON "stock_events" ("id");
CREATE INDEX "idx_stock_events_product_id" ON "stock_events" ("product_id");
CREATE INDEX "idx_stock_events_transaction_id" ON "stock_events" ("transaction_id");
CREATE INDEX "idx_stock_events_device_id" ON "stock_events" ("device_id");
CREATE INDEX "idx_stock_events_type" ON "stock_events" ("type");
CREATE INDEX "idx_stock_events_created_at" ON "stock_events" ("created_at");
CREATE INDEX "idx_stock_events_product_created" ON "stock_events" ("product_id", "created_at");