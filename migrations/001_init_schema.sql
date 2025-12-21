CREATE SCHEMA "public";
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
	CONSTRAINT "stock_events_source_check" CHECK (CHECK (((source)::text = ANY ((ARRAY['pos'::character varying, 'dashboard'::character varying, 'online'::character varying])::text[])))),
	CONSTRAINT "stock_events_type_check" CHECK (CHECK (((type)::text = ANY ((ARRAY['sale'::character varying, 'restock'::character varying, 'reject'::character varying, 'adjustment'::character varying, 'opening_stock'::character varying])::text[]))))
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
	"user_id" varchar(36)
);
CREATE TABLE "transactions" (
	"id" varchar(36) PRIMARY KEY,
	"user_id" varchar(36) NOT NULL,
	"total_amount" numeric(10, 2) NOT NULL,
	"status" varchar(20) DEFAULT 'pending' NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
CREATE TABLE "users" (
	"id" varchar(36) PRIMARY KEY,
	"username" varchar(50) NOT NULL CONSTRAINT "users_username_key" UNIQUE,
	"password" varchar(255) NOT NULL,
	"name" varchar(100) NOT NULL,
	"role" varchar(20) DEFAULT 'user' NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
ALTER TABLE "stock_events" ADD CONSTRAINT "fk_stock_events_product" FOREIGN KEY ("product_id") REFERENCES "products"("id") ON DELETE CASCADE;
ALTER TABLE "stock_events" ADD CONSTRAINT "fk_stock_events_transaction" FOREIGN KEY ("transaction_id") REFERENCES "transactions"("id") ON DELETE SET NULL;
ALTER TABLE "stock_events" ADD CONSTRAINT "fk_stock_events_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE SET NULL;
ALTER TABLE "transaction_items" ADD CONSTRAINT "transaction_items_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "products"("id");
ALTER TABLE "transaction_items" ADD CONSTRAINT "transaction_items_transaction_id_fkey" FOREIGN KEY ("transaction_id") REFERENCES "transactions"("id") ON DELETE CASCADE;
ALTER TABLE "transactions" ADD CONSTRAINT "transactions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id");
CREATE INDEX "idx_products_name" ON "products" ("name");
CREATE UNIQUE INDEX "products_pkey" ON "products" ("id");
CREATE INDEX "idx_stock_events_created_at" ON "stock_events" ("created_at");
CREATE INDEX "idx_stock_events_device_id" ON "stock_events" ("device_id");
CREATE INDEX "idx_stock_events_product_created" ON "stock_events" ("product_id","created_at");
CREATE INDEX "idx_stock_events_product_id" ON "stock_events" ("product_id");
CREATE INDEX "idx_stock_events_transaction_id" ON "stock_events" ("transaction_id");
CREATE INDEX "idx_stock_events_type" ON "stock_events" ("type");
CREATE UNIQUE INDEX "stock_events_pkey" ON "stock_events" ("id");
CREATE INDEX "idx_transaction_items_transaction_id" ON "transaction_items" ("transaction_id");
CREATE INDEX "idx_transaction_items_user_id" ON "transaction_items" ("user_id");
CREATE UNIQUE INDEX "transaction_items_pkey" ON "transaction_items" ("id");
CREATE INDEX "idx_transactions_status" ON "transactions" ("status");
CREATE INDEX "idx_transactions_user_id" ON "transactions" ("user_id");
CREATE UNIQUE INDEX "transactions_pkey" ON "transactions" ("id");
CREATE UNIQUE INDEX "users_pkey" ON "users" ("id");
CREATE UNIQUE INDEX "users_username_key" ON "users" ("username");