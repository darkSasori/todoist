package todoist

// Task is a model of todoist project entity
type Task struct {
	ID           int    `json:"id"`
	CommentCount int    `json:"comment_count"`
	Completed    bool   `json:"completed"`
	Content      string `json:"content"`
	Indent       int    `json:"indent"`
	LabelIDs     []int  `json:"label_ids"`
	Order        int    `json:"order"`
	Priority     int    `json:"priority"`
	ProjectID    int    `json:"project_id"`
	Due          Due    `json:"due"`
}

// Due is a model of todoist project entity
type Due struct {
	String   string     `json:"string"`
	Date     string     `json:"date"`
	Datetime CustomTime `json:"datetime"`
	Timezone string     `json:"timezone"`
}
