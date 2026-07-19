package ecdict

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	FullSourceURL = "https://raw.githubusercontent.com/skywind3000/ECDICT/master/ecdict.csv"
	// The official source currently yields 768,739 usable translated entries.
	// Reject substantially truncated downloads instead of activating a partial
	// database that merely happens to exceed the compact edition's 50k rows.
	minimumFullRows = 700_000
	downloadTimeout = 10 * time.Minute
)

type Status struct {
	EntryCount int    `json:"entryCount"`
	Edition    string `json:"edition"`
}

func (e *ECDict) Status() (Status, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.db == nil {
		return Status{}, errors.New("ecdict: database closed")
	}
	var status Status
	if err := e.db.QueryRow(`SELECT COUNT(*) FROM ecdict`).Scan(&status.EntryCount); err != nil {
		return Status{}, err
	}
	status.Edition = "compact"
	var edition string
	if err := e.db.QueryRow(`SELECT value FROM medict_meta WHERE key = 'edition'`).Scan(&edition); err == nil && edition != "" {
		status.Edition = edition
	}
	return status, nil
}

func (e *ECDict) InstallFull(ctx context.Context, client *http.Client) (int, error) {
	if client == nil {
		client = &http.Client{Timeout: downloadTimeout}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, FullSourceURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "Medict/3 ECDICT installer")
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("download ECDICT: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return 0, fmt.Errorf("download ECDICT: HTTP %d", resp.StatusCode)
	}
	return e.installFromCSV(ctx, resp.Body, minimumFullRows)
}

func (e *ECDict) installFromCSV(ctx context.Context, source io.ReadCloser, minimumRows int) (int, error) {
	defer source.Close()
	tmpPath := e.dbPath + ".download"
	backupPath := e.dbPath + ".compact-backup"
	_ = os.Remove(tmpPath)
	defer os.Remove(tmpPath)

	count, err := buildDatabaseFromCSV(ctx, source, tmpPath)
	if err != nil {
		return 0, err
	}
	if count < minimumRows {
		return 0, fmt.Errorf("ECDICT source only contained %d usable rows; need at least %d", count, minimumRows)
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	if e.db == nil {
		return 0, errors.New("ecdict: database closed")
	}
	if err := e.db.Close(); err != nil {
		return 0, fmt.Errorf("close compact ECDICT: %w", err)
	}
	e.db = nil

	_ = os.Remove(backupPath)
	if err := os.Rename(e.dbPath, backupPath); err != nil {
		e.db, _ = openDatabase(e.dbPath)
		return 0, fmt.Errorf("backup compact ECDICT: %w", err)
	}
	if err := os.Rename(tmpPath, e.dbPath); err != nil {
		_ = os.Rename(backupPath, e.dbPath)
		e.db, _ = openDatabase(e.dbPath)
		return 0, fmt.Errorf("activate full ECDICT: %w", err)
	}

	newDB, err := openDatabase(e.dbPath)
	if err != nil {
		_ = os.Remove(e.dbPath)
		_ = os.Rename(backupPath, e.dbPath)
		e.db, _ = openDatabase(e.dbPath)
		return 0, fmt.Errorf("verify full ECDICT: %w", err)
	}
	e.db = newDB
	_ = os.Remove(backupPath)
	return count, nil
}

func openDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM ecdict`).Scan(&count); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func buildDatabaseFromCSV(ctx context.Context, source io.Reader, dbPath string) (count int, err error) {
	_ = os.Remove(dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return 0, err
	}
	defer func() {
		db.Close()
		if err != nil {
			_ = os.Remove(dbPath)
		}
	}()
	if _, err = db.Exec(`
		CREATE TABLE ecdict (word TEXT PRIMARY KEY, phonetic TEXT, definition TEXT, translation TEXT, pos TEXT, frq INTEGER);
		CREATE TABLE medict_meta (key TEXT PRIMARY KEY, value TEXT NOT NULL);
	`); err != nil {
		return 0, err
	}

	reader := csv.NewReader(source)
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("read ECDICT header: %w", err)
	}
	columns := make(map[string]int, len(header))
	for i, name := range header {
		columns[strings.TrimSpace(name)] = i
	}
	for _, required := range []string{"word", "phonetic", "definition", "translation", "pos", "frq"} {
		if _, ok := columns[required]; !ok {
			return 0, fmt.Errorf("ECDICT source missing %q column", required)
		}
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `INSERT OR REPLACE INTO ecdict(word, phonetic, definition, translation, pos, frq) VALUES(?,?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	value := func(record []string, name string) string {
		idx := columns[name]
		if idx >= len(record) {
			return ""
		}
		return strings.TrimSpace(record[idx])
	}
	for {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return 0, fmt.Errorf("read ECDICT row: %w", readErr)
		}
		word, translation := value(record, "word"), value(record, "translation")
		if word == "" || translation == "" {
			continue
		}
		frequency, _ := strconv.Atoi(value(record, "frq"))
		if _, err := stmt.ExecContext(ctx, word, value(record, "phonetic"), value(record, "definition"), translation, value(record, "pos"), frequency); err != nil {
			return 0, fmt.Errorf("insert ECDICT row %q: %w", word, err)
		}
		count++
	}
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM ecdict`).Scan(&count); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO medict_meta(key, value) VALUES
		('edition', 'full'),
		('entry_count', ?),
		('source_url', ?),
		('license', 'MIT'),
		('copyright', 'Copyright (c) 2025 Linwei')`, count, FullSourceURL); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return count, nil
}
