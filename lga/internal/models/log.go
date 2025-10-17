package models

import "time"

type Path string
type Log struct {
	StatusCode string
	DateTime time.Time
	Path Path
}