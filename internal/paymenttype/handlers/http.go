package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/paymenttype/models"
	"github.com/hoitek/Maja-Service/internal/paymenttype/service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/paymenttype/config"
)

type PaymentTypeHandler struct {
	PaymentTypeService service.PaymentTypeService
}

func NewPaymentTypeHandler(r *mux.Router, s service.PaymentTypeService) (PaymentTypeHandler, error) {
	paymentTypeHandler := PaymentTypeHandler{
		PaymentTypeService: s,
	}
	if r == nil {
		return PaymentTypeHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.PaymentTypeConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.PaymentTypeConfig.ApiVersion1).Subrouter()

	rv1.Handle("/paymenttypes", paymentTypeHandler.Query()).Methods(http.MethodGet)
	return paymentTypeHandler, nil
}

/*
* @apiTag: paymenttype
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: paymenttypes
* @apiResponseRef: PaymentTypesQueryResponse
* @apiSummary: Query Payment Types
* @apiParametersRef: PaymentTypesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: PaymentTypesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: PaymentTypesQueryNotFoundResponse
 */
func (h *PaymentTypeHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.PaymentTypesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		paymentTypes, err := h.PaymentTypeService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(paymentTypes)
	})
}
