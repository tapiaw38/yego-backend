package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
	"wappi/internal/platform/config"
)

var (
	db   *sql.DB
	once sync.Once
)

// GetInstance returns the singleton database instance
func GetInstance() *sql.DB {
	once.Do(func() {
		cfg := config.GetInstance()
		var err error
		db, err = sql.Open("postgres", cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

		log.Println("Database connection established")
	})
	return db
}

// RunMigrations creates the necessary tables
func RunMigrations(db *sql.DB) error {
	schema := `
	-- Profile locations (must be created first)
	CREATE TABLE IF NOT EXISTS profile_locations (
		id UUID PRIMARY KEY,
		longitude DOUBLE PRECISION NOT NULL,
		latitude DOUBLE PRECISION NOT NULL,
		address TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	-- Profiles (depends on profile_locations)
	CREATE TABLE IF NOT EXISTS profiles (
		id UUID PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL UNIQUE,
		phone_number VARCHAR(50) NOT NULL,
		location_id UUID REFERENCES profile_locations(id),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id);

	-- Orders (depends on profiles)
	CREATE TABLE IF NOT EXISTS orders (
		id UUID PRIMARY KEY,
		profile_id UUID REFERENCES profiles(id),
		user_id VARCHAR(255),
		status VARCHAR(50) NOT NULL DEFAULT 'CREATED',
		eta VARCHAR(255),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	-- Add user_id column to orders if it doesn't exist (migration for existing tables)
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'orders' AND column_name = 'user_id') THEN
			ALTER TABLE orders ADD COLUMN user_id VARCHAR(255);
		END IF;
	END $$;

	-- Make profile_id nullable if it's NOT NULL (migration for existing tables)
	DO $$
	BEGIN
		IF EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_name = 'orders' AND column_name = 'profile_id' AND is_nullable = 'NO'
		) THEN
			ALTER TABLE orders ALTER COLUMN profile_id DROP NOT NULL;
		END IF;
	END $$;

	-- Add data column to orders (JSONB for order items)
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'orders' AND column_name = 'data') THEN
			ALTER TABLE orders ADD COLUMN data JSONB;
		END IF;
	END $$;

	CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
	CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
	CREATE INDEX IF NOT EXISTS idx_orders_profile_id ON orders(profile_id);
	CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);

	-- Order tokens (for claiming orders via link)
	CREATE TABLE IF NOT EXISTS order_tokens (
		id UUID PRIMARY KEY,
		order_id UUID NOT NULL REFERENCES orders(id),
		token VARCHAR(255) NOT NULL UNIQUE,
		phone_number VARCHAR(50),
		claimed_at TIMESTAMP WITH TIME ZONE,
		claimed_by_user_id VARCHAR(255),
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_order_tokens_token ON order_tokens(token);
	CREATE INDEX IF NOT EXISTS idx_order_tokens_order_id ON order_tokens(order_id);

	-- Profile tokens
	CREATE TABLE IF NOT EXISTS profile_tokens (
		id UUID PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		token VARCHAR(255) NOT NULL UNIQUE,
		used BOOLEAN DEFAULT FALSE,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_profile_tokens_token ON profile_tokens(token);
	CREATE INDEX IF NOT EXISTS idx_profile_tokens_user_id ON profile_tokens(user_id);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	log.Println("Database migrations completed")
	return nil
}
