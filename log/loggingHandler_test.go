package log

import (
	"net/http/httptest"
	"testing"
)

func TestStatusResponseWriter(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := &StatusResponseWriter{0, recorder}
	rw.WriteHeader(300)
	if res := rw.Status(); res != 300 {
		t.Errorf("Expected status to be 300, but was %d", res)
	}
}
