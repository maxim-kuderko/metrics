package metrics

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/metrics/drivers"
	"github.com/maxim-kuderko/metrics/entities"
	"math/rand"
	"net"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func BenchmarkReporter_Send(b *testing.B) {
	b.ReportAllocs()
	r := NewReporter(WithDriver(drivers.NewNoop()), WithConcurrency(2))
	name := `name`
	v := 1.0

	tagsAr := make([][]string, 0, 1000)
	for i := 0; i < 100; i++ {
		tagsAr = append(tagsAr, randArr())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Send(name, v, tagsAr[i%len(tagsAr)]...)
	}
}

func randArr() []string {
	c := int(rand.Int31n(10)) + 1
	output := make([]string, 0, c)
	for c > 0 {
		c--
		output = append(output, RandStringRunes(5))

	}
	return output
}

func BenchmarkReporter_Send_Concurrent(b *testing.B) {
	b.ReportAllocs()
	r := NewReporter(WithDriver(drivers.NewNoop()))
	name := `name`
	v := 0.1
	concurrency := 32
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	tagsAr := make([][]string, 0, 1000)
	for i := 0; i < 100; i++ {
		tagsAr = append(tagsAr, randArr())
	}
	b.ResetTimer()
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/concurrency; i++ {
				r.Send(name, v, tagsAr[i%len(tagsAr)]...)
			}
		}()
	}
	wg.Wait()
}

func TestReporter_Send(t *testing.T) {
	stu := drivers.NewTestStub()
	r := NewReporter(WithDriver(stu))
	count := 100000
	tagsAr := map[string][]string{}
	for i := 0; i < count; i++ {
		tagsA := randArr()
		k := ``
		for _, v := range tagsA {
			k += v
		}
		tagsAr[k] = tagsA
	}
	for _, v := range tagsAr {
		r.Send(`name`, 1.0, v...)
	}

	r.Close()
	c := int64(0)
	i := 0
	if len(tagsAr) != len(stu.Metrics()) {
		t.Fatalf(`bad aggregation expecting %v, got %v`, len(tagsAr), len(stu.Metrics()))
	}
	for _, m := range stu.Metrics() {
		c += m.Values.Count
		if m.Values.Count == 0 {
			fmt.Print(0)
		}
		k := ``
		for _, v := range m.Tags {
			k += v
		}
		if !reflect.DeepEqual(tagsAr[k], m.Tags) {
			t.Fatalf(`expecting %v, got %v`, tagsAr[k], m.Tags)
		}
		i++
	}

	if int(c) != count {
		t.Fatalf(`expecting %v, got %v`, count, c)
	}
}

func TestReporter_Send_Small(t *testing.T) {
	stu := drivers.NewTestStub()
	r := NewReporter(WithDriver(stu))
	count := 2000
	tagsAr := map[string][]string{}
	for i := 0; i < count; i++ {
		tagsA := randArr()
		r.Send(`name`, 1.0)
		k := ``
		for _, v := range tagsA {
			k += v
		}
		tagsAr[k] = tagsA
	}
	r.Close()
	c := int64(0)
	i := 0
	for _, m := range stu.Metrics() {
		c += m.Values.Count
		k := ``
		for _, v := range m.Tags {
			k += v
		}
		if !reflect.DeepEqual(tagsAr[k], m.Tags) {
			t.Fatalf(`expecting %v, got %v`, tagsAr[k], m.Tags)
		}
		i++
	}

	if int(c) != count {
		t.Fatalf(`expecting %v, got %v`, count, c)
	}
}

func TestReporter_SendC(t *testing.T) {
	stu := drivers.NewTestStub()
	r := NewReporter(WithDriver(stu))
	concurrency := 8
	count := 10000000 * concurrency
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	name := `name`
	v := 1.0
	tagsAr := make([][]string, 0, 1000)
	for i := 0; i < 1000; i++ {
		tagsA := randArr()
		k := ``
		for _, v := range tagsA {
			k += v
		}
		tagsAr = append(tagsAr, tagsA)
	}
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < count/concurrency; i++ {
				r.Send(name, v, tagsAr[i%len(tagsAr)]...)
			}
		}()
	}
	wg.Wait()
	r.Close()
	c := int64(0)
	for _, m := range stu.Metrics() {
		c += m.Values.Count
	}

	if int(c) != count {
		t.Fatalf(`expecting %v, got %v`, count, c)
	}
}

func TestReporter_Send_UDP(t *testing.T) {
	count := 10000000
	addr := `127.0.0.1:9999`
	addrS := net.UDPAddr{
		Port: 9999,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ln, _ := net.ListenUDP(`udp`, &addrS)
	wg := sync.WaitGroup{}
	wg.Add(2)
	c := int64(0)

	go func() {
		defer wg.Done()
		buff := make([]byte, 1<<20)
		for {
			ln.SetReadBuffer(1 << 20)
			ln.SetReadDeadline(time.Now().Add(time.Second))
			n, err := ln.Read(buff)
			if err != nil {
				if c != int64(count) {
					t.Fatalf("got %d expexted %d", c, count)
				}
				break
			}
			s := strings.Split(string(buff[:n]), "\n")
			for _, d := range s {
				if len(d) == 0 {
					continue
				}
				tmp := entities.AggregatedMetric{}
				if err := jsoniter.ConfigFastest.Unmarshal([]byte(d), &tmp); err != nil {
					fmt.Println(err)
					continue
				}
				c += tmp.Values.Count
			}
		}
	}()
	go func() {
		defer wg.Done()
		udp, _ := drivers.NewUDP(addr)
		cardinality := 1000
		r := NewReporter(WithDriver(udp), WithConcurrency(8), WithFlushTicker(time.Millisecond*10))
		tagsAr := make([][]string, 0, cardinality)
		for i := 0; i < cardinality; i++ {
			tagsAr = append(tagsAr, randArr())
		}
		for i := 0; i < count; i++ {
			r.Send(`name`, 1, tagsAr[i%cardinality]...)
		}
		fmt.Println(`closing`)
		r.Close()
		fmt.Println(udp.C)
	}()
	wg.Wait()
}

func TestReporter_Send_TCP(t *testing.T) {
	count := 1000000
	addr := `127.0.0.1:9999`
	addrS := net.TCPAddr{
		Port: 9999,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ln, _ := net.ListenTCP(`tcp`, &addrS)
	wg := sync.WaitGroup{}
	wg.Add(2)
	c := int64(0)

	go func() {
		defer wg.Done()
		buff := make([]byte, 1<<20)
		for {
			conn, err := ln.Accept()
			if err != nil {
				t.Error(err)
			}
			for {
				conn.SetReadDeadline(time.Now().Add(time.Second))
				n, err := conn.Read(buff)
				if err != nil {
					if c != int64(count) {
						t.Fatalf("got %d expexted %d", c, count)
					}
					return
				}
				s := strings.Split(string(buff[:n]), "\n")
				for _, d := range s {
					if len(d) == 0 {
						continue
					}
					tmp := entities.AggregatedMetric{}
					if err := jsoniter.ConfigFastest.Unmarshal([]byte(d), &tmp); err != nil {
						fmt.Println(err)
						continue
					}
					c += tmp.Values.Count
				}
			}
		}
	}()
	go func() {
		defer wg.Done()
		udp, _ := drivers.NewTCP(addr)
		cardinality := 1000
		r := NewReporter(WithDriver(udp), WithConcurrency(8), WithFlushTicker(time.Millisecond*10))
		tagsAr := make([][]string, 0, cardinality)
		for i := 0; i < cardinality; i++ {
			tagsAr = append(tagsAr, randArr())
		}
		for i := 0; i < count; i++ {
			r.Send(`name`, 1, tagsAr[i%cardinality]...)
		}
		fmt.Println(`closing`)
		r.Close()
		fmt.Println(udp.C)
	}()
	wg.Wait()
}

func BenchmarkReporter_Send_UDP(b *testing.B) {
	b.ReportAllocs()
	addr := `127.0.0.1:9999`
	udp, _ := drivers.NewUDP(addr)
	r := NewReporter(WithDriver(udp), WithConcurrency(1), WithFlushTicker(time.Second))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Send(`name`, 1)
	}
	r.Close()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
