package todoist

import "time"

type NullableTime struct {
	time.Time
}

type SyncResponse struct {
	User  *User  `json:"user,omitempty"`
	Items []Item `json:"items,omitempty"`
}

type CompletedTasksResponse struct {
	Items    []CompletedItem          `json:"items"`
	Projects map[string]ProjectDetail `json:"projects"`
}

type User struct {
	FullName string `json:"full_name"`
}

func (nt *NullableTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	t, err := time.Parse(`"2006-01-02T15:04:05.999999Z"`, string(b))
	if err != nil {
		return err
	}
	nt.Time = t
	return nil
}

type Item struct {
	Id          string        `json:"id"`
	Checked     bool          `json:"checked"`
	Content     string        `json:"content"`
	CompletedAt *NullableTime `json:"completed_at"`
}

type CompletedItem struct {
	Id          string        `json:"id"`
	TaskId      string        `json:"task_id"`
	ProjectId   string        `json:"project_id"`
	Content     string        `json:"content"`
	CompletedAt *NullableTime `json:"completed_at"`
}

type ProjectDetail struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ParsedItem struct {
	Id          string
	Content     string
	CompletedAt *NullableTime
	Project     string
}
