package models

import (
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
)

type (
	Avatar s3.Object
)

type PlayerPreview struct {
	Name    string
	Surname string
	Avatar  *Avatar
}
