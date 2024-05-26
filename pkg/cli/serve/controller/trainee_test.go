package serve

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateProgramRequest(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haGRpZWhtb2doaXNlaDgxdHJhaW5lZTJAZ21haWwuY29tIiwiZXhwIjoxNzE2NzU1Mjk2LCJ1c2VyX2lkIjoxfQ.IizpOZ5At3WcwYdR6jK4tGE8eDeRSKvJlT8FqBkcPnE"

	requestBody := dto.ProgramRequest{
		Description: "Sample program description",
		ActiveDays:  []bool{true, false, true, false, false, true, true},
	}

	requestBodyJSON, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "http://127.0.0.1:1234/trainee/request", bytes.NewBuffer(requestBodyJSON))
	req.Header.Set("Authorization", token)
	resp := httptest.NewRecorder()

	assert.Equal(t, http.StatusOK, resp.Code)
}
