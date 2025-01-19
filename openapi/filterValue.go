package openapi

/*
 * @apiDefine: FilterValueString
 */
type FilterValueString struct {
	Op    string `json:"op,omitempty" openapi:"example:eq;enum:eq,neq,gt,gte,lt,lte,like,notlike,ilike,notilike,in,notin,notnull,null;"`
	Value string `json:"value,omitempty" openapi:"example:1;"`
}

/*
 * @apiDefine: FilterValueInt
 */
type FilterValueInt struct {
	Op    string `json:"op,omitempty" openapi:"example:eq;enum:eq,neq,gt,gte,lt,lte,like,notlike,ilike,notilike,in,notin,notnull,null;"`
	Value int    `json:"value,omitempty" openapi:"example:1;"`
}
