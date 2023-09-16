package reteller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type HttpLogFields struct {
	Url     string              // both
	Method  string              // both
	IP      string              `json:"ip,omitempty"`     // both
	UserID  string              `json:"userID,omitempty"` // both
	Headers map[string][]string // both
	Body    string              // both

	Status  int           `json:"status,omitempty"`  // response
	Elapsed time.Duration `json:"elapsed,omitempty"` // response
}

// NewRequestFields creates a new HttpLogField struct and fills it with request data.
func NewRequestFields(request *http.Request) HttpLogFields {
	b, _ := ExportRequestBody(request).(string)

	return HttpLogFields{
		Url:     request.URL.String(),
		Method:  request.Method,
		IP:      request.RemoteAddr,
		Headers: ExportHeaders(request.Header),
		Body:    b,
		Status:  0,
		Elapsed: 0,
	}
}

// NewResponseFields creates a new HttpLogField struct. Url, Method and IP are provided from request.
func NewResponseFields(request *http.Request) HttpLogFields {
	return HttpLogFields{
		Url:    request.URL.String(),
		Method: request.Method,
		IP:     request.RemoteAddr,
	}
}

func ExportHeaders(h http.Header) map[string][]string {
	headers := make(map[string][]string)

	for name, values := range h.Clone() {
		headers[name] = values
	}

	return headers
}

func ExportRequestBody(r *http.Request) interface{} {
	var body interface{}
	var bodyBytes []byte

	if r == nil || r.Body == nil {
		return []byte(`{"error":"body is nil"}`)
	}

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		return map[string]string{"error": "failed to read response body"}
	}

	defer r.Body.Close()

	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if json.Valid(bodyBytes) {
		err := json.Unmarshal(bodyBytes, &body)
		if err != nil {
			return map[string]string{"error": "failed to unmarshal json"}
		}
		return body
	}

	return string(bodyBytes)
}
