package rules

type Trigger string

const (
	TriggerEvent Trigger = "event"
)

type Operator string

const (
	OpEquals    Operator = "equals"
	OpNotEquals Operator = "not_equals"
	OpContains  Operator = "contains"
)

type Condition struct {
	Field    string  `json:"field"`
	Operator Operator `json:"operator"`
	Value    string  `json:"value"`
}

type Action string

const (
	ActionNotify   Action = "notify"
	ActionExecute  Action = "execute"
)

type Rule struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Enabled     bool         `json:"enabled"`
	Trigger     Trigger      `json:"trigger"`
	Conditions  []Condition  `json:"conditions"`
	Actions     []Action     `json:"actions"`
}

type Match struct {
	Rule Rule
	Event Event
}

type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Severity  string                 `json:"severity"`
	Source    string                 `json:"source"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt interface{}            `json:"createdAt"`
}
