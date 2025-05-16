package media

import (
	"github.com/google/uuid"
	"path/filepath"
)

const (
	mediaFolder = "media"
)

func Key(courtID uuid.UUID, mediaID uuid.UUID) string {
	return filepath.Join(mediaFolder, courtID.String(), mediaID.String())
}
