package filters

type OperatorValue struct {
	Op    string
	Value interface{}
}

type FilterValue[T string | int] struct {
	Op    string `json:"op,omitempty"`
	Value T      `json:"value,omitempty"`
}
