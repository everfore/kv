package kv

type KV interface {
	Get(k string) ([]byte, error)
	Set(k string, v string) error
}

var (
	Kv KV
)

func SetKV(kv KV) {
	Kv = kv
}
