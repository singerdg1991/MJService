package integration_tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hoitek/Maja-Service/config"
	aConfig "github.com/hoitek/Maja-Service/internal/address/config"
	"github.com/hoitek/Maja-Service/internal/address/handlers"
	aRepo "github.com/hoitek/Maja-Service/internal/address/repositories"
	aService "github.com/hoitek/Maja-Service/internal/address/service"
	cRepo "github.com/hoitek/Maja-Service/internal/city/repositories"
	cService "github.com/hoitek/Maja-Service/internal/city/service"
	csRepo "github.com/hoitek/Maja-Service/internal/customer/repositories"
	csService "github.com/hoitek/Maja-Service/internal/customer/service"
	nRepo "github.com/hoitek/Maja-Service/internal/staff/repositories"
	nService "github.com/hoitek/Maja-Service/internal/staff/service"
	uRepo "github.com/hoitek/Maja-Service/internal/user/repositories"
	uService "github.com/hoitek/Maja-Service/internal/user/service"
	"github.com/hoitek/Maja-Service/router"
)

func TestAddressQueryIntegration(t *testing.T) {
	// Load default config
	config.LoadDefault()
	cfg := config.LoadAddressConfig()
	aConfig.AddressConfig = &cfg

	// Create a new instance of address service
	addressService := aService.NewAddressService(aRepo.NewAddressRepositoryStub(), nil)
	cityService := cService.NewCityService(cRepo.NewCityRepositoryStub(), nil)
	staffService := nService.NewStaffService(nRepo.NewStaffRepositoryStub(), nRepo.NewStaffRepositoryMongoDBStub(), nil)
	customerService := csService.NewCustomerService(csRepo.NewCustomerRepositoryStub(), nil, nil)
	userService := uService.NewUserService(uRepo.NewUserRepositoryStub(), nil, nil)

	// Create a new instance of HTTP server handler
	handler, err := handlers.NewAddressHandler(
		router.Init(),
		addressService,
		&cityService,
		&staffService,
		&customerService,
		userService,
	)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	// Create a test server based on your handler
	server := httptest.NewServer(handler.Query())
	defer server.Close()

	// Make a request to the test server
	response, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("failed to make GET request: %v", err)
	}

	// Ensure the response status code is as expected
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, response.StatusCode)
	}

	// Read body from response
	bodyBuffer, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("failed to read response body: %v", err)
	}

	// Close response body
	response.Body.Close()

	// Unmarshal response body
	var body map[string]interface{}
	if err := json.Unmarshal(bodyBuffer, &body); err != nil {
		t.Errorf("failed to unmarshal response body: %v", err)
	}

	// Validate response body
	if _, ok := body["statusCode"]; !ok {
		t.Errorf("response body does not contain statusCode")
	}
	if _, ok := body["data"]; !ok {
		t.Errorf("response body does not contain data")
	}
	responseBodyData, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Errorf("response body data is not an object")
	}
	if _, ok := responseBodyData["limit"]; !ok {
		t.Errorf("response body data does not contain limit")
	}
	if _, ok := responseBodyData["offset"]; !ok {
		t.Errorf("response body data does not contain offset")
	}
	if _, ok := responseBodyData["page"]; !ok {
		t.Errorf("response body data does not contain page")
	}
	if _, ok := responseBodyData["totalRows"]; !ok {
		t.Errorf("response body data does not contain totalRows")
	}
	if _, ok := responseBodyData["totalPages"]; !ok {
		t.Errorf("response body data does not contain totalPages")
	}
	if _, ok := responseBodyData["items"]; !ok {
		t.Errorf("response body data does not contain items")
	}
}
