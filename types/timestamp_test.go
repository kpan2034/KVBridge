package types_test

import (
	. "KVBridge/types"
	"bytes"
	"testing"
	"time"
)

const (
	high48Mask uint64 = 0xFFFFFFFFFFFF0000
	roundMask  uint64 = 0x0000000000008000
	low16Mask  uint64 = 0x000000000000FFFF
)

func TestTimestampFromUnixTime(t *testing.T) {
	tunix := time.Now().UnixNano()
	// t.Logf("tunix:\t %16x", tunix)
	ts := uint64(tunix)
	// t.Logf("uint64:\t %16x", ts)
	round := ts & roundMask
	// t.Logf("round?:\t %16x", round)
	ts = ts & high48Mask
	// t.Logf("high 48:\t %16x", ts)
	ts = ts + (round << 1)
	// t.Logf("final:\t %16x", ts)
}

func TestTimestampTick(t *testing.T) {
	t1 := NewTimestamp()
	t2 := t1.Tick()

	if t2.Uint64()-t1.Uint64() != 1 {
		t.Errorf("tick invalid: expected: %v, got: %v", t2, t1)
	}
	return
}

func TestTimestampL(t *testing.T) {
	t1 := NewTimestamp()
	t1 = t1.Tick().Tick().Tick()
	t2 := t1.L()
	if (t2.Uint64() & 0x000000000000FFFF) != 0 {
		t.Errorf("last 16 bytes should be 0, got: t2")
	}
}

func TestTimestampC(t *testing.T) {
	t1 := NewTimestamp()
	t1 = t1.Tick().Tick().Tick()
	t2 := t1.C()
	if (t2.Uint64() & 0xFFFFFFFFFFFF0000) != 0 {
		t.Errorf("first 48 bytes should be 0, got: %s", t2)
	}
	if t2.Uint64() != 3 {
		t.Errorf("tick value should be 3, got: %s", t2)
	}
}

func TestTimestampSend(t *testing.T) {
	t1 := NewTimestamp()
	time.Sleep(10 * time.Millisecond)
	t2 := t1.Send()
	// last 16 bytes of t2 should be zero
	if t2.C().Uint64() != 0 {
		t.Errorf("send did not update ts: got: %v, old: %v", t2, t1)
	}
	return
}

func TestTimestampReceive(t *testing.T) {
	t1 := NewTimestamp()
	// t.Logf("ts: %s", ts)
	m := NewTimestamp()
	// t.Logf("m: %s", m)

	// lj := ts.L()
	// lm := m.L()
	// latest := lj.Max(lm).Max(NewTimestamp())
	// if lj.Equals(latest) && lj.Equals(lm) {
	// 	cj := ts.C().Max(m.C()).Tick()
	// 	t.Logf("cj: %s", cj)
	// 	t.Logf("ts b1: %16x", FromUnixTime(int64(lj.Uint64())).Uint64()+cj.Uint64())
	// 	return
	// } else if latest.Equals(lj) {
	// 	t.Logf("ts b2: %s", ts.Tick())
	// 	return
	// } else if latest.Equals(lm) {
	// 	t.Logf("ts b3: %s", m.Tick())
	// 	return
	// }
	// t.Logf("ts b4: %s", latest)
	//
	res1 := t1.Receive(m)

	if res1.Equals(t1) || res1.Equals(m) || !res1.Equals(res1.Max(t1)) {
		t.Errorf("receive failed: expected res1 > t1 and res1 > m, got: %v", res1)
	}
}

func TestTimestampEncode(t *testing.T) {
	randomUnixTime := int64(0x0001000200030004)
	ts := FromUnixTime(randomUnixTime).Tick().Tick().Tick().Tick()
	encoded := ts.Encode()
	expected := []byte{0x04, 0x00, 0x03, 0x00, 0x02, 0x00, 0x01, 0x00}

	if !bytes.Equal(encoded, expected) {
		t.Errorf("got: %v, expected: %v", encoded, expected)
	}
}

func TestTimestampDecode(t *testing.T) {
	randomBytes := []byte{0x04, 0x00, 0x03, 0x00, 0x02, 0x00, 0x01, 0x00}
	decoded, err := DecodeToTimestamp(randomBytes)
	if err != nil {
		t.Errorf("error decoding to timestamp: %v", err)
	}
	expectedUnixTime := int64(0x0001000200030004)
	expected := FromUnixTime(expectedUnixTime).Tick().Tick().Tick().Tick()

	if !decoded.Equals(expected) {
		t.Errorf("got: %v, expected: %v", decoded, expected)
	}
}
