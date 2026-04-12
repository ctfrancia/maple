package scraperdb

import (
	"context"
	"fmt"
	"time"

	"github.com/ctfrancia/maple/business/domain/scraperbus"
	"github.com/ctfrancia/maple/business/sdk/order"
	"github.com/ctfrancia/maple/business/sdk/page"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Store manages the set of APIs for scraper database access.
type Store struct {
	log *logger.Logger
	db  *gorm.DB
}

// NewStore constructs the API for data access.
func NewStore(log *logger.Logger, db *gorm.DB) (*Store, error) {
	psql, err := db.DB()
	if err != nil {
		return nil, err
	}

	psql.SetMaxOpenConns(10)
	psql.SetMaxIdleConns(5)
	psql.SetConnMaxLifetime(5 * time.Minute)

	if err := psql.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &Store{
		log: log,
		db:  db,
	}, nil
}

// CreateMigration creates the database schema for scraper tables.
func CreateMigration(db *gorm.DB) error {
	if err := db.AutoMigrate(&scraperbus.ScrapeJob{}, &scraperbus.ScrapedTournament{}); err != nil {
		return fmt.Errorf("failed to migrate scraper tables: %w", err)
	}
	return nil
}

// CreateJob creates a new scrape job.
func (s *Store) CreateJob(ctx context.Context, job scraperbus.ScrapeJob) error {
	return s.db.WithContext(ctx).Create(&job).Error
}

// UpdateJob updates an existing scrape job with pessimistic locking.
func (s *Store) UpdateJob(ctx context.Context, job scraperbus.ScrapeJob, update scraperbus.UpdateScrapeJob) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing scraperbus.ScrapeJob
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&existing, job.ID).Error; err != nil {
			return err
		}

		// Apply updates
		updates := map[string]interface{}{
			"date_updated": time.Now(),
		}

		if update.Status != nil {
			updates["status"] = *update.Status
		}
		if update.CompletedAt != nil {
			updates["completed_at"] = *update.CompletedAt
		}
		if update.Error != nil {
			updates["error"] = *update.Error
		}
		if update.Found != nil {
			updates["found"] = *update.Found
		}
		if update.Created != nil {
			updates["created"] = *update.Created
		}
		if update.Updated != nil {
			updates["updated"] = *update.Updated
		}

		return tx.Model(&existing).Updates(updates).Error
	})
}

// QueryJobs retrieves a list of scrape jobs based on the provided filter.
func (s *Store) QueryJobs(ctx context.Context, filter scraperbus.QueryFilter, orderBy order.By, pg page.Page) ([]scraperbus.ScrapeJob, error) {
	query := s.db.WithContext(ctx)

	if filter.ID != nil {
		if id, err := uuid.Parse(*filter.ID); err == nil {
			query = query.Where("id = ?", id)
		}
	}
	if filter.Federation != nil {
		query = query.Where("federation = ?", *filter.Federation)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.StartTime != nil {
		query = query.Where("started_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("started_at <= ?", *filter.EndTime)
	}

	// Apply ordering
	if orderBy.Field != "" {
		direction := "ASC"
		if orderBy.Direction == order.DESC {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", orderBy.Field, direction))
	} else {
		query = query.Order("date_created DESC")
	}

	// Apply pagination
	query = query.Offset(pg.Number() * pg.RowsPerPage()).Limit(pg.RowsPerPage())

	var jobs []scraperbus.ScrapeJob
	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

// QueryJobByID retrieves a scrape job by ID.
func (s *Store) QueryJobByID(ctx context.Context, jobID uuid.UUID) (scraperbus.ScrapeJob, error) {
	var job scraperbus.ScrapeJob
	if err := s.db.WithContext(ctx).First(&job, jobID).Error; err != nil {
		return scraperbus.ScrapeJob{}, err
	}
	return job, nil
}

// UpsertTournament creates or updates a scraped tournament.
func (s *Store) UpsertTournament(ctx context.Context, tournament scraperbus.ScrapedTournament) error {
	now := time.Now()
	tournament.DateUpdated = now
	tournament.LastScrapedAt = now

	// Use Clauses with OnConflict to upsert
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "url", "federation", "organizer", "arbiter", "location", "start_date", "players", "rounds", "time_control", "last_scraped_at", "date_updated"}),
	}).Create(&tournament).Error
}

// QueryTournaments retrieves a list of scraped tournaments based on the provided filter.
func (s *Store) QueryTournaments(ctx context.Context, filter scraperbus.QueryFilter, orderBy order.By, pg page.Page) ([]scraperbus.ScrapedTournament, error) {
	query := s.db.WithContext(ctx)

	if filter.ID != nil {
		if id, err := uuid.Parse(*filter.ID); err == nil {
			query = query.Where("id = ?", id)
		}
	}
	if filter.Federation != nil {
		query = query.Where("federation = ?", *filter.Federation)
	}

	// Apply ordering
	if orderBy.Field != "" {
		direction := "ASC"
		if orderBy.Direction == order.DESC {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", orderBy.Field, direction))
	} else {
		query = query.Order("last_scraped_at DESC")
	}

	// Apply pagination
	query = query.Offset(pg.Number() * pg.RowsPerPage()).Limit(pg.RowsPerPage())

	var tournaments []scraperbus.ScrapedTournament
	if err := query.Find(&tournaments).Error; err != nil {
		return nil, err
	}

	return tournaments, nil
}

// QueryTournamentByExternalID retrieves a scraped tournament by external ID.
func (s *Store) QueryTournamentByExternalID(ctx context.Context, externalID string) (scraperbus.ScrapedTournament, error) {
	var tournament scraperbus.ScrapedTournament
	if err := s.db.WithContext(ctx).Where("external_id = ?", externalID).First(&tournament).Error; err != nil {
		return scraperbus.ScrapedTournament{}, err
	}
	return tournament, nil
}
