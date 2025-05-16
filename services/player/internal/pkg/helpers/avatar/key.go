package avatar

import (
	"github.com/google/uuid"
	"path/filepath"
)

const (
	avatarsFolder = "avatars"
)

func Key(userID uuid.UUID, avatarID uuid.UUID) string {
	return filepath.Join(avatarsFolder, userID.String(), avatarID.String())
}
