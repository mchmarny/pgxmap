package pgxmap

// Item represents a single item in a ConfigMap that's persisted into DB.
type Item struct {
	Key         string      `json:"k"`
	Value       interface{} `json:"v"`
	Type        string      `json:"t"`
	Transformed bool        `json:"b"`
}
