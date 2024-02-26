package files

import "time"

type FileDetails struct {
	Name string
	Path string
}

type ParsedDailyNote struct {
	Path      string
	CreatedAt time.Time
}
