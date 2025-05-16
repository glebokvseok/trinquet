package media

import (
	"github.com/google/uuid"
	"path/filepath"
)

const (
	mediaFolder = "media"
)

func Key(mediaID uuid.UUID) string {
	return filepath.Join(mediaFolder, mediaID.String())
}
