package avcodec

import (
	"runtime"
	"testing"
	"time"
)

// fakePool records Return routing and leak-cleanup callbacks. Unlike a real
// pool it does not Unref or re-arm returned packets, which lets the tests
// observe the raw Packet-side semantics.
type fakePool struct {
	returned chan *Packet
	leaks    chan bool
}

func newFakePool() *fakePool {
	return &fakePool{returned: make(chan *Packet, 4), leaks: make(chan bool, 4)}
}

func (f *fakePool) ReturnPacket(p *Packet)   { f.returned <- p }
func (f *fakePool) PacketLeaked(hadBuf bool) { f.leaks <- hadBuf }

// waitForLeak GCs until the leak cleanup reports, or fails the test.
func waitForLeak(t *testing.T, leaks chan bool) bool {
	t.Helper()
	deadline := time.After(5 * time.Second)
	for {
		runtime.GC()
		select {
		case hadBuf := <-leaks:
			return hadBuf
		case <-deadline:
			t.Fatal("leak cleanup did not fire")
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// assertNoLeakCallback asserts the leak cleanup stays silent across GCs.
func assertNoLeakCallback(t *testing.T, leaks chan bool, msg string) {
	t.Helper()
	for i := 0; i < 3; i++ {
		runtime.GC()
		time.Sleep(20 * time.Millisecond)
	}
	select {
	case <-leaks:
		t.Fatal(msg)
	default:
	}
}

func TestPacket_MakeWritable_COWsSharedBuffer(t *testing.T) {
	src := NewPacket()
	defer src.Free()
	if err := src.SetBytes([]byte{1, 2, 3, 4}); err != nil {
		t.Fatalf("SetBytes: %v", err)
	}
	dst := NewPacket()
	defer dst.Free()
	if err := src.Ref(dst); err != nil {
		t.Fatalf("Ref: %v", err)
	}
	if src.Data() != dst.Data() {
		t.Fatal("expected Ref to share the underlying buffer")
	}

	if err := dst.MakeWritable(); err != nil {
		t.Fatalf("MakeWritable: %v", err)
	}
	if src.Data() == dst.Data() {
		t.Fatal("expected MakeWritable to copy the shared buffer")
	}

	dst.GetDataUnsafe()[0] = 99
	if got := src.GetDataAt(0); got != 1 {
		t.Fatalf("mutation leaked into peer packet: got %d, want 1", got)
	}
	if got := dst.GetDataAt(0); got != 99 {
		t.Fatalf("mutation lost after COW: got %d, want 99", got)
	}

	// Sole owner: MakeWritable must be a no-op, not a copy.
	before := dst.Data()
	if err := dst.MakeWritable(); err != nil {
		t.Fatalf("MakeWritable (sole owner): %v", err)
	}
	if dst.Data() != before {
		t.Fatal("MakeWritable copied a buffer it exclusively owned")
	}
}

func TestPacket_Return_PooledRoute_StopsCleanup(t *testing.T) {
	pool := newFakePool()
	func() {
		pkt := NewPacketWithReturner(pool)
		if err := pkt.SetBytes([]byte{1, 2, 3}); err != nil {
			t.Fatalf("SetBytes: %v", err)
		}
		pkt.Return()
	}()

	var got *Packet
	select {
	case got = <-pool.returned:
	default:
		t.Fatal("Return on a pool-owned packet did not route to the pool")
	}
	if got.CAVPacket == nil {
		t.Fatal("pool-routed packet must not be freed")
	}

	// Return must have stopped the original cleanup handle. Re-arm as a
	// real pool would before caching, then drop the packet: exactly one
	// callback proves the original handle was dead (a still-armed original
	// would double-fire and double-free the C packet).
	got.ArmLeakCleanup()
	got = nil
	waitForLeak(t, pool.leaks)
	assertNoLeakCallback(t, pool.leaks, "original cleanup fired after Return stopped it")
}

func TestPacket_Return_NilReturner_FreesC(t *testing.T) {
	pkt := NewPacketWithReturner(nil)
	if err := pkt.SetBytes([]byte{1, 2, 3}); err != nil {
		t.Fatalf("SetBytes: %v", err)
	}
	pkt.Return()
	if pkt.CAVPacket != nil {
		t.Fatal("Return on a foreign packet must free the C packet")
	}

	plain := NewPacket()
	plain.Return()
	if plain.CAVPacket != nil {
		t.Fatal("Return on a plain NewPacket packet must free the C packet")
	}
}

func TestPacket_Free_OnPoolOwned_Panics(t *testing.T) {
	pool := newFakePool()
	pkt := NewPacketWithReturner(pool)

	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("Free on a pool-owned packet must panic")
			}
		}()
		pkt.Free()
	}()

	// StopCleanup severs ownership, after which Free is legal and must
	// not trigger the (now stopped) leak cleanup.
	pkt.StopCleanup()
	pkt.Free()
	if pkt.CAVPacket != nil {
		t.Fatal("Free after StopCleanup must free the C packet")
	}
	assertNoLeakCallback(t, pool.leaks, "cleanup fired after StopCleanup")
}

func TestPacket_LeakedInFlight_CleanupFiresWithBuffer(t *testing.T) {
	pool := newFakePool()
	func() {
		pkt := NewPacketWithReturner(pool)
		if err := pkt.SetBytes([]byte{9, 9, 9}); err != nil {
			t.Fatalf("SetBytes: %v", err)
		}
		// Dropped without Return: leaked in flight.
	}()
	if hadBuf := waitForLeak(t, pool.leaks); !hadBuf {
		t.Fatal("in-flight leak must report hadBuffer=true")
	}
}

func TestPacket_LeakedIdle_CleanupFiresWithoutBuffer(t *testing.T) {
	pool := newFakePool()
	func() {
		pkt := NewPacketWithReturner(pool)
		if err := pkt.SetBytes([]byte{9, 9, 9}); err != nil {
			t.Fatalf("SetBytes: %v", err)
		}
		pkt.Unref() // what a real pool does before caching
		// Dropped while idle: models a sync.Pool eviction.
	}()
	if hadBuf := waitForLeak(t, pool.leaks); hadBuf {
		t.Fatal("idle eviction must report hadBuffer=false")
	}
}
