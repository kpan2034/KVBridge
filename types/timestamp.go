package types

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"time"
)

// HLC Timestamp is a 64-bit value
// Based on Demirbas, Murat et al.
// “Logical Physical Clocks and Consistent Snapshots in Globally Distributed Databases.” (2014).
// Use Timestamp, not *Timestamps! They are cheap to copy over, you know...
type Timestamp struct {
	ts uint64
}

const (
	high48Mask uint64 = 0xFFFFFFFFFFFF0000
	roundMask  uint64 = 0x0000000000008000
	low16Mask  uint64 = 0x000000000000FFFF
)

// Returns physical time, rounded up
// returns a 64-bit integer with the lower 16 bits set to zero
func NewTimestamp() Timestamp {
	return FromUnixTime(time.Now().UnixNano())
}

func FromUnixTime(t int64) Timestamp {
	ts := uint64(t)
	round := ts & roundMask
	ts = ts & high48Mask
	ts = ts + (round << 1)
	return Timestamp{ts}
}

// Gets 'l' component of timestamp
func (t Timestamp) L() Timestamp {
	return Timestamp{t.ts & high48Mask}
}

// Gets 'c' component of timestamp
func (t Timestamp) C() Timestamp {
	return Timestamp{t.ts & low16Mask}
}

// Returns a new Timestamp which is the max of t or r
func (t Timestamp) Max(r Timestamp) Timestamp {
	if t.ts > r.ts {
		return t
	}
	return r
}

// Returns true if both the 'l' and 'c' components of t and r are equal
func (t Timestamp) Equals(r Timestamp) bool {
	return t.ts == r.ts
}

// Returns a new Timestamp with 'c' component incremented by 1
func (t Timestamp) Tick() Timestamp {
	return Timestamp{t.ts + 1}
}

// Exposes the underlying uint64 component of the Timestamp
func (t Timestamp) Uint64() uint64 {
	return t.ts
}

// Returns an updated Timestamp after a send or local event
func (t Timestamp) Send() Timestamp {
	lj := t.L()
	latest := lj.Max(NewTimestamp())
	if lj.Equals(latest) {
		return t.Tick() // this updates cj
	}
	return latest
}

// Returns an updated Timestamp after a receive event
func (t Timestamp) Receive(m Timestamp) Timestamp {
	lj := t.L()
	lm := m.L()
	latest := lj.Max(lm).Max(NewTimestamp())
	if lj.Equals(latest) && lj.Equals(lm) {
		cj := t.C().Max(m.C()).Tick()
		return Timestamp{lj.ts | cj.ts}
	} else if latest.Equals(lj) {
		return t.Tick()
	} else if latest.Equals(lm) {
		return m.Tick()
	}
	return latest
}

// Timestamp implements Stringer interface
func (t Timestamp) String() string {
	return strconv.FormatUint(t.ts, 16)
}

func (t Timestamp) Encode() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, t.ts)
	return b
}

func DecodeToTimestamp(b []byte) (Timestamp, error) {
	buf := bytes.NewReader(b)
	var ts uint64
	err := binary.Read(buf, binary.LittleEndian, &ts)
	if err != nil {
		return Timestamp{}, err
	}
	return Timestamp{ts: ts}, nil
}
