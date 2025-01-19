package restypes

type QueryResponse struct {
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
	Page       int         `json:"page"`
	TotalRows  int64       `json:"totalRows"`
	TotalPages int         `json:"totalPages"`
	Items      interface{} `json:"items"`
}

type DeleteResponse struct {
	IDs []uint `json:"ids"`
}
