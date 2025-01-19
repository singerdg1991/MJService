package models

import (
	"encoding/json"
	"log"
	"net/http"

	govalidity "github.com/hoitek/Govalidity"
)

/*
 * @apiDefine: ReportsQueryCustomerAgeGroupChartAggregationType
 */
type ReportsQueryCustomerAgeGroupChartAggregationType struct {
	OperationType string `json:"operationType" openapi:"example:equal"`
	Field         string `json:"field" openapi:"example:s.userId"`
	Value         string `json:"value" openapi:"example:1"`
}

/*
 * @apiDefine: ReportsQueryCustomerAgeGroupChartRequestParams
 */
type ReportsQueryCustomerAgeGroupChartRequestParams struct {
	SectionIDs             interface{}                                        `json:"sectionIds,omitempty" openapi:"example:[1,2,3];type:array;"`
	Aggregations           interface{}                                        `json:"aggregations,omitempty" openapi:"example:[{\"operationType\":\"equal\",\"field\":\"s.userId\",\"value\":\"1\"}];type:array;"`
	FilterIDs              interface{}                                        `json:"filterIds,omitempty" openapi:"example:[1,2,3];type:array;"`
	SectionIDsAsInt64Slice []int64                                            `json:"-" openapi:"ignored"`
	AggregationsValue      []ReportsQueryCustomerAgeGroupChartAggregationType `json:"-" openapi:"ignored"`
	FilterIDsAsInt64Slice  []int64                                            `json:"-" openapi:"ignored"`
}

// ValidateQueries validates the query parameters
func (data *ReportsQueryCustomerAgeGroupChartRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"sectionIds":   govalidity.New("sectionIds").Optional(),
		"aggregations": govalidity.New("aggregations").Optional(),
		"filterIds":    govalidity.New("filterIds").Optional(),
	}

	// Validate the data
	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		return govalidity.DumpErrors(errs)
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
					"aggregations": []string{"Invalid aggregations"},
				}
			}
			data.Aggregations = aggInterface
		}

		aggregationsInterface, ok := data.Aggregations.([]interface{})
		if !ok {
			return govalidity.ValidityResponseErrors{
				"aggregations": []string{"Invalid aggregations"},
			}
		}

		for _, aggregationInterface := range aggregationsInterface {
			aggregationMap, ok := aggregationInterface.(map[string]interface{})
			if !ok {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid aggregation format"},
				}
			}

			operationType, ok := aggregationMap["operationType"].(string)
			if !ok {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid operationType"},
				}
			}

			field, ok := aggregationMap["field"].(string)
			if !ok {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid field"},
				}
			}

			value, ok := aggregationMap["value"].(string)
			if !ok {
				return govalidity.ValidityResponseErrors{
					"aggregations": []string{"Invalid value"},
				}
			}

			data.AggregationsValue = append(data.AggregationsValue, ReportsQueryCustomerAgeGroupChartAggregationType{
				OperationType: operationType,
				Field:         field,
				Value:         value,
			})
		}
		log.Printf("Aggregations: %v\n", data.AggregationsValue)
	}

	return nil
}
