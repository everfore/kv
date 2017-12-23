package etcd

import (
	"errors"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	endpoints      = []string{"127.0.0.1:2379"}

	Etcd *EtcdKV
)

func init() {
	Etcd = new(EtcdKV)
	var err error
	Etcd.client, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		println(err)
	}
}

type EtcdKV struct {
	client *clientv3.Client
}

func (e *EtcdKV) Get(key string) ([]byte, error) {
	resp, err := e.client.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	size := len(resp.Kvs)
	if size <= 0 {
		return nil, errors.New("nil keys")
	}
	return resp.Kvs[0].Value, nil
}

func (e *EtcdKV) Set(key string, value string) error {
	_, err := e.client.Put(context.Background(), key, value)
	return err
}

func kk() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		println(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Put(ctx, "/test/hello", "world")
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	// resp, err := cli.Get(ctx, "a/b/k1")
	resp, err := cli.Get(ctx, "/test/hello")
	cancel()

	for _, ev := range resp.Kvs {
		fmt.Printf("get %s : %s\n", ev.Key, ev.Value)
	}

	dlt, err := cli.Delete(context.TODO(), "a/b/k1")
	fmt.Println("err:", err)
	for _, ev := range dlt.PrevKvs {
		fmt.Printf("delete %s : %s\n", ev.Key, ev.Value)
	}

	go func() {
		time.Sleep(3e9)
		_, err = cli.Put(context.TODO(), "key", "xyz")
		ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
		_, err = cli.Txn(ctx).
			If(clientv3.Compare(clientv3.Value("key"), ">", "abc")).
			Then(clientv3.OpPut("key", "XYZ")).
			Else(clientv3.OpPut("key", "ABC")).
			Commit()
		cancel()

	}()

	rch := cli.Watch(context.Background(), "key", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("watch %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	<-rch

	if err != nil {
		println(err)
	}
}
