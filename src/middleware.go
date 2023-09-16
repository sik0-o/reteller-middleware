package reteller

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r = InjectCtx(r)

		// Создаем буферезированного писателя
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		ww.Tee(CtxBuffer(r))

		// лог запишим после выполнения миддлвары
		t1 := time.Now()
		defer func() {
			resp := ResponseData(r)

			resp.Headers = ExportHeaders(ww.Header())

			if ww.BytesWritten() > 0 {
				if ent, ok := CtxEntry(r); ok {
					resp.Body = ent.buf.String()
				}
			}

			resp.Status = ww.Status()
			resp.Elapsed = time.Since(t1)
		}()

		// выполняем миддлвару
		next.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}
