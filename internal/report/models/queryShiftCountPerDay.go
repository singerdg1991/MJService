package models

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: ReportsQueryShiftCountPerDayChartRequestParams
 */
type ReportsQueryShiftCountPerDayChartRequestParams struct {
	SectionIDs             interface{}                                        `json:"sectionIds,omitempty" openapi:"example:[1,2,3];type:array;"`
	Aggregations           interface{}                                        `json:"aggregations,omitempty" openapi:"example:[{\"operationType\":\"equal\",\"field\":\"s.userId\",\"value\":\"1\"}];type:array;"`
	FilterIDs              interface{}                                        `json:"filterIds,omitempty" openapi:"example:[1,2,3];type:array;"`
	SectionIDsAsInt64Slice []int64                                            `json:"-" openapi:"ignored"`
	AggregationsValue      []ReportsQueryShiftCountPerDayChartAggregationType `json:"-" openapi:"ignored"`
	FilterIDsAsInt64Slice  []int64                                            `json:"-" openapi:"ignored"`
}

/*
 * @apiDefine: ReportsQueryShiftCountPerDayChartAggregationType
 */
type ReportsQueryShiftCountPerDayChartAggregationType struct {
	OperationType string `json:"operationType" openapi:"example:equal"`
	Field         string `json:"field" openapi:"example:s.userId"`
	Value         string `json:"value" openapi:"example:1"`
}

func (data *ReportsQueryShiftCountPerDayChartRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"sectionIds":   govalidity.New("sectionIds").Optional(),
		"aggregations": govalidity.New("aggregations").Optional(),
		"filterIds":    govalidity.New("filterIds").Optional(),
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
		var aggInterface interface{}

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
		var aggregations []ReportsQueryShiftCountPerDayChartAggregationType
		for _, aggregationInterface := range aggregationsInterface {
			aggregationMap, ok := aggregationInterface.(map[string]interface{})
			if !ok {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid aggregations, must be an array of objects and each object must have operationType, field and value"},
				}
			}

			// Check operation type
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

			// Check field
			field, ok := aggregationMap["field"].(string)
			if !ok {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid aggregations, field must be a string"},
				}
			}
			var acceptableFields = []string{
				"cps.id",
				"c.sectionId",
				"cps.cycleId",
				"s.id",
				"s.userId",
				"u.firstName",
				"u.lastName",
				"u.avatarUrl",
				"cst.id",
				"cst.shiftName",
				"cst.startHour",
				"cst.endHour",
				"cst.isUnplanned",
				"cst.datetime",
				"cps.status",
				"cps.prevStatus",
				"cps.startKilometer",
				"cps.reasonOfTheCancellation",
				"cps.reasonOfTheReactivation",
				"cps.reasonOfTheResume",
				"cps.reasonOfThePause",
				"cps.isUnplanned",
				"cps.datetime",
				"cps.created_at",
				"cps.updated_at",
				"cps.deleted_at",
				"cps.started_at",
				"cps.ended_at",
				"cps.cancelled_at",
				"cps.delayed_at",
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
			aggregation := ReportsQueryShiftCountPerDayChartAggregationType{
				OperationType: operationType,
				Field:         field,
				Value:         value,
			}
			aggregations = append(aggregations, aggregation)
		}
		data.AggregationsValue = aggregations
		log.Printf("Aggregations: %v\n", data.AggregationsValue)
	}

	return nil
}
