package people

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"

	"github.com/ystv/web-api/services/people"
	"github.com/ystv/web-api/utils"
)

type (
	// Repo stores our dependencies
	Repo struct {
		people *people.Store
		access *utils.Accesser
	}
)

// NewRepo creates our data store
func NewRepo(db *sqlx.DB, cdn *s3.S3, access *utils.Accesser) *Repo {
	return &Repo{
		people: people.NewStore(db, cdn),
		access: access,
	}
}
