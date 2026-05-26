package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(databaseURL string) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for range 15 {
		db, err = sqlx.Connect("postgres", databaseURL)
		if err == nil && db != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	if db == nil {
		return nil, fmt.Errorf("failed to connect to postgres: connection is nil")
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	const bootstrapDDL = `
	CREATE TABLE IF NOT EXISTS organizations (
	    id SERIAL PRIMARY KEY,
	    name VARCHAR(255) NOT NULL,
	    short_name VARCHAR(100) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS week_templates (
	    id SERIAL PRIMARY KEY,
	    group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
	    type VARCHAR(10) NOT NULL DEFAULT 'EVEN' CHECK (type IN ('EVEN', 'ODD', 'OTHER'))
	);
	CREATE INDEX IF NOT EXISTS idx_week_templates_group ON week_templates(group_id);

	CREATE TABLE IF NOT EXISTS days (
	    id SERIAL PRIMARY KEY,
	    week_template_id INT NOT NULL REFERENCES week_templates(id) ON DELETE CASCADE,
	    day_index INT NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_days_template ON days(week_template_id);

	CREATE TABLE IF NOT EXISTS day_slots (
	    id SERIAL PRIMARY KEY,
	    day_id INT NOT NULL REFERENCES days(id) ON DELETE CASCADE,
	    start_minute INT NOT NULL,
	    end_minute INT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_day_slots_day ON day_slots(day_id);

	ALTER TABLE lessons ADD COLUMN IF NOT EXISTS day_id INT REFERENCES days(id) ON DELETE CASCADE;
	CREATE INDEX IF NOT EXISTS idx_lessons_day_id ON lessons(day_id);
	`
	_, _ = db.ExecContext(ctx, bootstrapDDL)

	return db, nil
}
