package models

import (
	"encoding/json"
	"net/http"
	"strings"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: ReportsQueryShiftsFilterType
 */
type ReportsQueryShiftsFilterType struct {
}

/*
 * @apiDefine: ReportsQueryShiftsSortType
 */
type ReportsQueryShiftsSortType struct {
}

/*
 * @apiDefine: ReportsQueryShiftsTableAggregationType
 */
type ReportsQueryShiftsTableAggregationType struct {
	OperationType string `json:"operationType" openapi:"example:equal"`
	Field         string `json:"field" openapi:"example:cs.id"`
	Value         string `json:"value" openapi:"example:1"`
}

/*
 * @apiDefine: ReportsQueryShiftsTableRequestParams
 */
type ReportsQueryShiftsTableRequestParams struct {
	ID                     int                                      `json:"id,string,omitempty" openapi:"example:1"`
	SectionIDs             interface{}                              `json:"sectionIds,omitempty" openapi:"example:[1,2,3];type:array;"`
	Aggregations           interface{}                              `json:"aggregations,omitempty" openapi:"example:[{\"operationType\":\"equal\",\"field\":\"cs.id\",\"value\":\"1\"}];type:array;"`
	FilterIDs              interface{}                              `json:"filterIds,omitempty" openapi:"example:[1,2,3];type:array;"`
	Page                   int                                      `json:"page,string,omitempty" openapi:"example:1"`
	Limit                  int                                      `json:"limit,string,omitempty" openapi:"example:10"`
	Filters                ReportsQueryShiftsFilterType             `json:"filters,omitempty" openapi:"$ref:ReportsQueryShiftsFilterType;in:query"`
	Sorts                  ReportsQueryShiftsSortType               `json:"sorts,omitempty" openapi:"$ref:ReportsQueryShiftsSortType;in:query"`
	SectionIDsAsInt64Slice []int64                                  `json:"-" openapi:"ignored"`
	AggregationsValue      []ReportsQueryShiftsTableAggregationType `json:"-" openapi:"ignored"`
	FilterIDsAsInt64Slice  []int64                                  `json:"-" openapi:"ignored"`
}

func (data *ReportsQueryShiftsTableRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":           govalidity.New("id").Int().Optional(),
		"sectionIds":   govalidity.New("sectionIds").Optional(),
		"aggregations": govalidity.New("aggregations").Optional(),
		"filterIds":    govalidity.New("filterIds").Optional(),
		"page":         govalidity.New("page").Int().Default("1"),
		"limit":        govalidity.New("limit").Int().Default("10"),
		"filters":      govalidity.Schema{},
		"sorts":        govalidity.Schema{},
	}

	// Validate the data
	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	// Convert the SectionIDs to int slice
	if data.SectionIDs != nil {
		sIdsString, ok := data.SectionIDs.(string)
		if !ok {
			return govalidity.ValidityResponseErrors{
				"sectionIds": []string{"Invalid sectionIds, must be an array of integers"},
			}
		}
		var sIdsInt []int
		err := json.Unmarshal([]byte(sIdsString), &sIdsInt)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"sectionIds": []string{"Invalid sectionIds, must be an array of integers"},
			}
		}
		for _, sId := range sIdsInt {
			data.SectionIDsAsInt64Slice = append(data.SectionIDsAsInt64Slice, int64(sId))
		}
	}

	// Convert the FilterIDs to int slice
	if data.FilterIDs != nil {
		fIdsString, ok := data.FilterIDs.(string)
		if !ok {
			return govalidity.ValidityResponseErrors{
				"filterIds": []string{"Invalid filterIds, must be an array of integers"},
			}
		}
		var fIdsInt []int
		err := json.Unmarshal([]byte(fIdsString), &fIdsInt)
		if err != nil {
			return govalidity.ValidityResponseErrors{
				"filterIds": []string{"Invalid filterIds, must be an array of integers"},
			}
		}
		for _, fId := range fIdsInt {
			data.FilterIDsAsInt64Slice = append(data.FilterIDsAsInt64Slice, int64(fId))
		}
	}

	// Convert the Aggregations to AggregationsValue
	if data.Aggregations != nil {
		var (
			aggInterface interface{}
		)

		aggString, ok := data.Aggregations.(string)
		if ok {
			err := json.Unmarshal([]byte(aggString), &aggInterface)
			if err != nil {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid aggregations, must be an array of objects and each object must have operationType, field and value"},
				}
			}
			data.Aggregations = aggInterface
		}
		aggregationsInterface, ok := data.Aggregations.([]interface{})
		if !ok {
			return govalidity.ValidityResponseErrors{
				"aggregations": []string{"Invalid aggregations, must be an array of objects and each object must have operationType, field and value"},
			}
		}
		var aggregations []ReportsQueryShiftsTableAggregationType
		for _, aggregationInterface := range aggregationsInterface {
			aggregationMap, ok := aggregationInterface.(map[string]interface{})
			if !ok {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid aggregations, must be an array of objects and each object must have operationType, field and value"},
				}
			}
			operationType, ok := aggregationMap["operationType"].(string)
			if !ok {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid aggregations, operationType must be a string"},
				}
			}
			var acceptableOperationTypes = []string{
				"equal",
				"notEqual",
				"greaterThan",
				"lessThan",
				"contains",
			}
			var found = false
			for _, acceptableOperationType := range acceptableOperationTypes {
				if operationType == acceptableOperationType {
					found = true
					break
				}
			}
			if !found {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid aggregations, operationType must be one of the following: " + strings.Join(acceptableOperationTypes, ", ")},
				}
			}

			field, ok := aggregationMap["field"].(string)
			if !ok {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid aggregations, field must be a string"},
				}
			}
			var acceptableFields = []string{
				"cs.id",
				"cs.exchangeKey",
				"cs.cycleId",
				"cs.staffTypeIds",
				"cs.shiftName",
				"cs.vehicleType",
				"cs.startLocation",
				"cs.datetime",
				"cs.status",
				"cs.created_at",
				"cs.updated_at",
				"cs.deleted_at",
			}
			found = false
			for _, acceptableField := range acceptableFields {
				if field == acceptableField {
					found = true
					break
				}
			}
			if !found {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid aggregations, field must be one of the following: " + strings.Join(acceptableFields, ", ")},
				}
			}

			value, ok := aggregationMap["value"].(string)
			if !ok {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid aggregations, value must be a string"},
				}
			}
			aggregation := ReportsQueryShiftsTableAggregationType{
				OperationType: operationType,
				Field:         field,
				Value:         value,
			}
			aggregations = append(aggregations, aggregation)
		}
		data.AggregationsValue = aggregations
	}

	return nil
}
