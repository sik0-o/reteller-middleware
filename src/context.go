package reteller

var (
	ctxkey = &contextKey{"mainEntry"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "reteller.middleware:" + k.name
}
