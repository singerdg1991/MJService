package quilder

import (
	"errors"
	"fmt"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Go-Quilder/utils"
	"strings"
)

type Statement struct {
	Condition string
	Values    Values
	Enable    bool
}

type (
	QueryFields  = []string
	Statements   = []Statement
	QueriesGroup = map[string]*Statements
	Values       = []interface{}
)

type Query struct {
	RawQuery     string
	QueriesGroup *QueriesGroup
	Limit        int
	Offset       int
	Fields       *QueryFields
}

type QueryResult struct {
	Where  string
	Limit  int
	Offset int
	Fields *QueryFields
}

func (qo *Query) sanitizeValues(values Values) []string {
	strValues := []string{}
	for _, v := range values {
		strValue := ""
		switch v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			strValue = fmt.Sprintf("%v", v)
		case string:
			strValue = v.(string)
		default:
			strValue = fmt.Sprintf("%v", v)
		}
		strValues = append(strValues, strValue)
	}
	return strValues
}

func (q *Query) sanitizeStatements(operator string, statements Statements) string {
	qs := []string{}
	for _, statement := range statements {
		if !statement.Enable {
			continue
		}
		stSlice := strings.Split(statement.Condition, "?")
		sanitizedValues := q.sanitizeValues(statement.Values)
		if len(stSlice) != len(sanitizedValues)+1 {
			continue
		}
		for i, sv := range sanitizedValues {
			stSlice[i] = stSlice[i] + "'" + sv + "'"
		}
		st := strings.TrimSpace(strings.Join(stSlice, " "))
		qs = append(qs, st)
	}
	return strings.Join(qs, " "+operator+" ")
}

func (q *Query) makeQueries() string {
	queryResult := ""
	if *q.QueriesGroup != nil {
		queries := []string{}
		for k, v := range *q.QueriesGroup {
			q := q.sanitizeStatements(k, *v)
			queries = append(queries, q)
		}
		queryResult = strings.Join(queries, " "+operators.AND+" ")
	}
	return queryResult
}

func (q *Query) Build() *QueryResult {
	where := strings.TrimSpace(q.RawQuery)

	if where == "" {
		where = q.makeQueries()
	}

	return &QueryResult{
		Where:  where,
		Limit:  q.Limit,
		Offset: q.Offset,
		Fields: q.Fields,
	}
}

func CreateQueriesGroup[T interface{}](queries T) (*QueriesGroup, error) {
	queriesHashMap := utils.ToMap(queries)

	fMap, ok := queriesHashMap["filters"]
	if !ok {
		return nil, errors.New("filters not found")
	}
	filtersHashMap := fMap.(map[string]interface{})
	statements := []Statement{}

	for k, v := range filtersHashMap {
		field := v.(map[string]interface{})
		statement := utils.ParseSQLOperator(fmt.Sprintf("%v", field["op"]), field["value"])
		filtersHashMap[k] = filters.FilterValue[string]{Op: statement.Op, Value: statement.Value.(string)}
		queriesHashMap["filters"].(map[string]interface{})[k] = filtersHashMap[k]
		number, _ := utils.ToNumber(statement.Value)
		statements = append(statements, Statement{
			Condition: fmt.Sprintf("%s %s ?", k, statement.Op),
			Values:    Values{statement.Value},
			Enable:    statement.Value != nil && len(fmt.Sprintf("%v", statement.Value)) > 0 || (number != nil && *number > 0),
		})
	}

	group := &QueriesGroup{
		operators.AND: &statements,
	}

	return group, nil
}
