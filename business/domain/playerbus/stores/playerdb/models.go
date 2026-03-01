package playerdb

//		if _, err := stmt.ExecContext(ctx, tournamentID, p.Rank, p.Title, p.Name, p.FideID, p.Federation, p.Rating, p.Points, p.RatingPerf); err != nil {

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Player represents a player in the database.
type Player struct {
	gorm.Model
	PublicID  uuid.UUID `gorm:"type:uuid"`
	FirstName string
	LastName  string
	Email     string
}
