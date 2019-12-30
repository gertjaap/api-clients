package kimai

import "time"

type TimeEntry struct {
	ID          int    `json:"id"`
	Project     int    `json:"project"`
	StartText   string `json:"begin"`
	EndText     string `json:"end"`
	Description string `json:"description"`
	Start       time.Time
	End         time.Time
}

func (c *Client) GetTimeEntries() ([]TimeEntry, error) {
	entries := make([]TimeEntry, 0)
	req, err := c.Get("timesheets")
	if err != nil {
		return entries, err
	}
	req.Limit(10000)
	err = req.Json(&entries)
	if err != nil {
		return entries, err
	}

	for i, e := range entries {
		entries[i].Start, err = time.Parse("2006-01-02T15:04:05-0700", e.StartText)
		entries[i].End, err = time.Parse("2006-01-02T15:04:05-0700", e.EndText)
		if entries[i].End.Year() == 1 {
			entries[i].End = time.Now()
		}
	}

	return entries, err
}
