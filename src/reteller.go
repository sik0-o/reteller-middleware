package reteller

import (
	"context"
	"io"
	"net/http"
)

func InjectCtx(r *http.Request) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), ctxkey, NewEntry(r)))

	return r
}

func CtxEntry(r *http.Request) (*ctxEntry, bool) {
	e, ok := r.Context().Value(ctxkey).(*ctxEntry)

	return e, ok
}

func CtxBuffer(r *http.Request) io.Writer {
	entry, ok := r.Context().Value(ctxkey).(*ctxEntry)
	if !ok {
		return nil
	}

	return entry.Buffer()
}

func RequestData(r *http.Request) *HttpLogFields {
	e, ok := CtxEntry(r)
	if !ok {
		return nil
	}

	return &e.request
}

func ResponseData(r *http.Request) *HttpLogFields {
	e, ok := CtxEntry(r)
	if !ok {
		return nil
	}

	return &e.response
}
