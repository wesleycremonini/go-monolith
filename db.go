package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

func connectDB(path string) (*sql.DB, error) {
	// busy_timeout(10000): This pragma sets the busy timeout value to 10,000 milliseconds. The busy timeout determines how long SQLite will wait for a database lock to be released before returning an error.
	// journal_mode(WAL): This pragma sets the journal mode to Write-Ahead Logging (WAL). The WAL mode provides better concurrency and performance for most workloads compared to other journaling modes.
	// journal_size_limit(200000000): This pragma sets the maximum size limit for the write-ahead log file to 200,000,000 bytes. When the log file reaches this limit, it will be checkpointed and a new log file will be created.
	// synchronous(NORMAL): This pragma sets the synchronous mode to NORMAL. In this mode, SQLite will sync the database to disk at critical points but not necessarily for every transaction, which provides a balance between durability and performance.
	// foreign_keys(ON): This pragma enables the enforcement of foreign key constraints. When foreign keys are enabled, SQLite will enforce referential integrity by checking the validity of foreign key relationships between tables.
	pragmas := "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)"

	db, err := sql.Open("sqlite", path+pragmas)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}
