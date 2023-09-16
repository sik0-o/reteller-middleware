package reteller

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

var (
	ctxkey = &contextKey{"reteller.middleware"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "ctx key:" + k.name
}

type ctxEntry struct {
	buf      *bytes.Buffer
	r        *http.Request
	request  HttpLogFields
	response HttpLogFields
}

func (ent *ctxEntry) Buffer() io.Writer {
	return ent.buf
}

func NewEntry(r *http.Request) *ctxEntry {
	return &ctxEntry{
		buf:      &bytes.Buffer{},
		r:        r,
		request:  NewRequestFields(r),
		response: NewResponseFields(r),
	}
}

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
