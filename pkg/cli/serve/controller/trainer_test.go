package serve

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrainerController_GetAllRequests(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haGRpZWhtb2doaXNlaDgxdHJhaW5lcjJAZ21haWwuY29tIiwiZXhwIjoxNzE2NzU1MzI3LCJ1c2VyX2lkIjoyfQ.EhOIyG5fUDM3aP7g11aUu98QDiEsrVs9yCWKi8EklD8"

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:1234/trainer/requests", nil)
	req.Header.Set("Authorization", token)
	resp := httptest.NewRecorder()

	assert.Equal(t, http.StatusOK, resp.Code)
}
