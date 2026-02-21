package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	models "github.com/ctfrancia/maple/api/services/scrape/domain"
	"github.com/ctfrancia/maple/foundation/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Store interface {
	// KnownTournamentIDs returns all tournament IDs we've already seen for a federation.
	KnownTournamentIDs(ctx context.Context, federation string) (map[string]bool, error)

	// UpsertTournament inserts or updates a tournament record. Returns true if it was new.
	UpsertTournament(ctx context.Context, t *models.Tournament) (isNew bool, err error)

	// UpdateStatus sets the tournament status (discovered/scraped/completed/error).
	UpdateStatus(ctx context.Context, tournamentID, status string) error

	// UpsertPlayers replaces the player list for a tournament.
	UpsertPlayers(ctx context.Context, tournamentID string, players []models.Player) error

	// TournamentsToScrape returns tournaments that need detail scraping.
	// (status = "discovered" or last_scraped older than maxAge)
	TournamentsToScrape(ctx context.Context, maxAge time.Duration) ([]models.Tournament, error)

	// Close releases resources.
	Close() error
}

// ---------- In-memory implementation (for dev/testing) ----------

// MemoryStore is a simple in-memory store for development and tests.
type MemoryStore struct {
	mu          sync.RWMutex
	tournaments map[string]*models.Tournament
	players     map[string][]models.Player
}

// NewMemoryStore returns a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tournaments: make(map[string]*models.Tournament),
		players:     make(map[string][]models.Player),
	}
}

func (m *MemoryStore) KnownTournamentIDs(_ context.Context, federation string) (map[string]bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ids := make(map[string]bool)
	for id, t := range m.tournaments {
		if federation == "" || t.Federation == federation {
			ids[id] = true
		}
	}
	return ids, nil
}

func (m *MemoryStore) UpsertTournament(_ context.Context, t *models.Tournament) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.tournaments[t.ID]
	now := time.Now()
	if !exists {
		t.CreatedAt = now
		t.Status = "discovered"
	}
	t.UpdatedAt = now
	m.tournaments[t.ID] = t
	return !exists, nil
}

func (m *MemoryStore) UpdateStatus(_ context.Context, id, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t, ok := m.tournaments[id]; ok {
		t.Status = status
		t.UpdatedAt = time.Now()
	}
	return nil
}

func (m *MemoryStore) UpsertPlayers(_ context.Context, tournamentID string, players []models.Player) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.players[tournamentID] = players
	if t, ok := m.tournaments[tournamentID]; ok {
		t.LastScraped = time.Now()
		t.Status = "scraped"
	}
	return nil
}

func (m *MemoryStore) TournamentsToScrape(_ context.Context, maxAge time.Duration) ([]models.Tournament, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []models.Tournament
	cutoff := time.Now().Add(-maxAge)
	for _, t := range m.tournaments {
		if t.Status == "discovered" || t.LastScraped.Before(cutoff) {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (m *MemoryStore) Close() error { return nil }

// ---------- PostgreSQL implementation ----------

// PostgresStore persists data in PostgreSQL.
type PostgresStore struct {
	//db     *sql.DB
	db     *gorm.DB
	logger *logger.Logger
}

func NewPostgresStore(databaseURL string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	s := &PostgresStore{db: db, logger: slog.Default()}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}
	return s, nil
}

func (s *PostgresStore) migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS tournaments (
			id          TEXT PRIMARY KEY,
			name        TEXT NOT NULL DEFAULT '',
			location    TEXT NOT NULL DEFAULT '',
			federation  TEXT NOT NULL DEFAULT '',
			start_date  TEXT NOT NULL DEFAULT '',
			end_date    TEXT NOT NULL DEFAULT '',
			players     INTEGER NOT NULL DEFAULT 0,
			rounds      INTEGER NOT NULL DEFAULT 0,
			time_control TEXT NOT NULL DEFAULT '',
			arbiter     TEXT NOT NULL DEFAULT '',
			organizer   TEXT NOT NULL DEFAULT '',
			url         TEXT NOT NULL DEFAULT '',
			status      TEXT NOT NULL DEFAULT 'discovered',
			last_scraped TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01T00:00:00Z',
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tournaments_federation ON tournaments(federation)`,
		`CREATE INDEX IF NOT EXISTS idx_tournaments_status ON tournaments(status)`,
		`CREATE TABLE IF NOT EXISTS players (
			tournament_id TEXT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
			rank          INTEGER NOT NULL DEFAULT 0,
			title         TEXT NOT NULL DEFAULT '',
			name          TEXT NOT NULL DEFAULT '',
			fide_id       TEXT NOT NULL DEFAULT '',
			federation    TEXT NOT NULL DEFAULT '',
			rating        INTEGER NOT NULL DEFAULT 0,
			points        DOUBLE PRECISION NOT NULL DEFAULT 0,
			rating_perf   INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (tournament_id, rank)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_players_fide_id ON players(fide_id) WHERE fide_id != ''`,
	}
	for _, m := range migrations {
		if _, err := s.db.Exec(m); err != nil {
			return fmt.Errorf("executing migration: %w", err)
		}
	}
	return nil
}

func (s *PostgresStore) KnownTournamentIDs(ctx context.Context, federation string) (map[string]bool, error) {
	query := `SELECT id FROM tournaments WHERE federation = $1`
	rows, err := s.db.QueryContext(ctx, query, federation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make(map[string]bool)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids[id] = true
	}
	return ids, rows.Err()
}

func (s *PostgresStore) UpsertTournament(ctx context.Context, t *models.Tournament) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1)`, t.ID).Scan(&exists)
	if err != nil {
		return false, err
	}

	query := `INSERT INTO tournaments (id, name, location, federation, start_date, end_date, players, rounds, time_control, arbiter, organizer, url, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		ON CONFLICT (id) DO UPDATE SET
			name=EXCLUDED.name, location=EXCLUDED.location, federation=EXCLUDED.federation,
			start_date=EXCLUDED.start_date, end_date=EXCLUDED.end_date, players=EXCLUDED.players,
			rounds=EXCLUDED.rounds, time_control=EXCLUDED.time_control, arbiter=EXCLUDED.arbiter,
			organizer=EXCLUDED.organizer, url=EXCLUDED.url, updated_at=NOW()`

	status := t.Status
	if status == "" {
		status = "discovered"
	}
	_, err = s.db.ExecContext(ctx, query,
		t.ID, t.Name, t.Location, t.Federation, t.StartDate, t.EndDate,
		t.Players, t.Rounds, t.TimeControl, t.Arbiter, t.Organizer, t.URL, status,
	)
	return !exists, err
}

func (s *PostgresStore) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE tournaments SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (s *PostgresStore) UpsertPlayers(ctx context.Context, tournamentID string, players []models.Player) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM players WHERE tournament_id = $1`, tournamentID); err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO players (tournament_id, rank, title, name, fide_id, federation, rating, points, rating_perf)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range players {
		if _, err := stmt.ExecContext(ctx, tournamentID, p.Rank, p.Title, p.Name, p.FideID, p.Federation, p.Rating, p.Points, p.RatingPerf); err != nil {
			return err
		}
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE tournaments SET status='scraped', last_scraped=NOW(), updated_at=NOW() WHERE id=$1`,
		tournamentID); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostgresStore) TournamentsToScrape(ctx context.Context, maxAge time.Duration) ([]models.Tournament, error) {
	cutoff := time.Now().Add(-maxAge)
	query := `SELECT id, name, location, federation, url, status
		FROM tournaments
		WHERE status = 'discovered' OR last_scraped < $1
		ORDER BY created_at DESC
		LIMIT 50`

	rows, err := s.db.QueryContext(ctx, query, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Tournament
	for rows.Next() {
		var t models.Tournament
		if err := rows.Scan(&t.ID, &t.Name, &t.Location, &t.Federation, &t.URL, &t.Status); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}
