package host_test

import (
	"testing"

	"github.com/yottachain/P2PHost"
)

var localMa = []string{"/ip4/0.0.0.0/tcp/9000"}
var localMa2 = []string{"/ip4/0.0.0.0/tcp/9001"}
var localMa3 = []string{"/ip4/0.0.0.0/tcp/9002"}

func TestNewHost(t *testing.T) {
	host, err := host.NewHost(localMa2, nil)
	if err != nil {
		t.Fatalf("new host error: %s", err)
	} else {
		t.Log(host)
	}
}

// TestSendMessage 测试发送接受消息
func TestSendMessage(t *testing.T) {
	// 创建host1模拟接受消息
	h1, err := host.NewHost(localMa, nil)
	h1.RegisterHandler("ping", func(msg host.Msg) []byte {
		if string(msg.MsgType) == "ping" {
			return []byte("pong")
		} else {
			return []byte("error")
		}
	})
	// 创建host2模拟发送消息
	h2, err := host.NewHost(localMa2, nil)
	if err != nil {
		t.Fatalf("new host error: %s", err)
	} else {
		t.Log("new host success")
	}
	// 连接节点1
	err = h2.Connect(h1.ID(), h1.Addrs())
	if err != nil {
		t.Fatalf("connect err :%s", err)
	} else {
		t.Log("connect success")
	}
	// 发送ping
	res, err := h2.SendMsg(h1.ID(), "ping", []byte("ping"))
	if err != nil {
		t.Fatalf("sendMsg err :%s", err)

	} else {
		t.Log("sendMsg success")
		if string(res) == "pong" {
			t.Log("res success")
		} else {
			t.Fatal(string(res))
		}
	}
}

// TestRealy 测试中继节点连接
func TestRealy(t *testing.T) {
	h1, err := host.NewHost(localMa, nil)
	h2, err := host.NewHost(localMa2, nil)
	h3, err := host.NewHost(localMa3, nil)

	h1.Connect(h2.ID(), h2.Addrs())
	h3.Connect(h2.ID(), h2.Addrs())

	h1.Connect(h3.ID(), []string{"/p2p-circuit/p2p/" + h3.ID().Pretty()})

	h3.RegisterHandler("ping", func(msg host.Msg) []byte {
		if string(msg.MsgType) == "ping" {
			return []byte("pong")
		} else {
			return nil
		}
	})
	res, err := h1.SendMsg(h3.ID(), "ping", []byte("ping"))
	if err != nil {
		t.Fatalf("err:%s", err)
	} else {
		if string(res) != "pong" {
			t.Fatalf("err:%s", res)
		}
	}
	t.Log(string(res))
}
