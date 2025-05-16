package responses

import (
	"github.com/google/uuid"
)

type CreateMatchResponse struct {
	MatchID uuid.UUID `json:"match_id"`
}
