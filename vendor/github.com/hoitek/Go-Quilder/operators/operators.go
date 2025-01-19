package operators

const (
	AND = "AND"
	OR  = "OR"
)

const (
	EQUALS                     = "equals"
	CONTAINS                   = "contains"
	STARTS_WITH                = "startsWith"
	ENDS_WITH                  = "endsWith"
	IS_EMPTY                   = "isEmpty"
	IS_NOT_EMPTY               = "isNotEmpty"
	IS_ANY_OF                  = "isAnyOf"
	NUMBER_EQUALS              = "="
	NUMBER_NOT_EQUALS          = "!="
	NUMBER_GREATER_THAN        = ">"
	NUMBER_GREATER_THAN_EQUALS = ">="
	NUMBER_LESS_THAN           = "<"
	NUMBER_LESS_THAN_EQUALS    = "<="
	DATE_IS                    = "is"
	DATE_IS_NOT                = "isNot"
	DATE_IS_AFTER              = "isAfter"
	DATE_IS_BEFORE             = "isBefore"
	DATE_IS_ON_OR_AFTER        = "isOnOrAfter"
	DATE_IS_ON_OR_BEFORE       = "isOnOrBefore"
)

var SQL = map[string]string{
	EQUALS:                     "=",
	CONTAINS:                   "LIKE",
	STARTS_WITH:                "LIKE",
	ENDS_WITH:                  "LIKE",
	IS_EMPTY:                   "IS NULL",
	IS_NOT_EMPTY:               "IS NOT NULL",
	IS_ANY_OF:                  "IN",
	NUMBER_EQUALS:              "=",
	NUMBER_NOT_EQUALS:          "!=",
	NUMBER_GREATER_THAN:        ">",
	NUMBER_GREATER_THAN_EQUALS: ">=",
	NUMBER_LESS_THAN:           "<",
	NUMBER_LESS_THAN_EQUALS:    "<=",
	DATE_IS:                    "=",
	DATE_IS_NOT:                "!=",
	DATE_IS_AFTER:              ">",
	DATE_IS_BEFORE:             "<",
	DATE_IS_ON_OR_AFTER:        ">=",
	DATE_IS_ON_OR_BEFORE:       "<=",
}
