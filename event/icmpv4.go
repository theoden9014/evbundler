package event

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// not work on
type ICMPv4Event struct {
	seq        int
	listenAddr string
	addr       string
}

func NewICMPv4Event(seq int, listenAddr, addr string) *ICMPv4Event {
	return &ICMPv4Event{seq: seq, listenAddr: listenAddr, addr: addr}
}

func (e ICMPv4Event) Name() string {
	return "ICMPv4"
}

func (e ICMPv4Event) Fire(ctx context.Context) error {
	// start listening for icmp replies
	c, err := icmp.ListenPacket("ip4:icmp", e.listenAddr)
	if err != nil {
		return err
	}
	defer c.Close()

	// resolve any DNS
	dst, err := net.ResolveIPAddr("ip4", e.addr)
	if err != nil {
		return err
	}

	// make icmp packet
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: e.seq,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	b, err := m.Marshal(nil)
	if err != nil {
		return err
	}

	// send a icmp packet
	if _, err := c.WriteTo(b, dst); err != nil {
		return err
	}

	// wait for a reply icmp packet
	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		return err
	}
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return err
	}

	// check a replay icmp packet
	rm, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return err
	}

	switch rm.Type {
	case ipv6.ICMPTypeEchoReply:
		return nil
	default:
		return fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}
}
