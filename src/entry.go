package reteller

import (
	"bytes"
	"io"
	"net/http"
)

type ctxEntry struct {
	buf      *bytes.Buffer
	request  HttpLogFields
	response HttpLogFields
}

func (ent *ctxEntry) Buffer() io.Writer {
	return ent.buf
}

func NewEntry(r *http.Request) *ctxEntry {
	return &ctxEntry{
		buf:      &bytes.Buffer{},
		request:  NewRequestFields(r),
		response: NewResponseFields(r),
	}
}

func (ent *ctxEntry) Request() HttpLogFields {
	return ent.request
}

func (ent *ctxEntry) Response() HttpLogFields {
	return ent.response
}
