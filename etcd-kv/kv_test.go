package etcd

import (
	"testing"
)

func TestKV(t *testing.T) {
	t.Log(Etcd.Get("a"))

	Etcd.Set("a", "value")
	t.Log(Etcd.Get("a"))
}
