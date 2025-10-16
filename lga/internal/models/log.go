package models

import "time"

type Log struct {
	StatusCode string
	DateTime time.Time
	Path string
}