package utils

import "time"

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
func Now() Timestamp {
	ts := uint64(time.Now().Unix()) << 1
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
func (t Timestamp) tick() Timestamp {
	return Timestamp{t.ts + 1}
}

// Exposes the underlying uint64 component of the Timestamp
func (t Timestamp) UInt64() uint64 {
	return t.ts
}

// Returns an updated Timestamp after a send or local event
func (t Timestamp) Send() Timestamp {
	lj := t.L()
	latest := lj.Max(Now())
	if lj.Equals(latest) {
		return t.tick() // this updates cj
	}
	return latest
}

// Returns an updated Timestamp after a receive event
func (t Timestamp) Receive(m Timestamp) Timestamp {
	lj := t.L()
	lm := m.L()
	latest := lj.Max(lm).Max(Now())
	if lj.Equals(latest) && lj.Equals(lm) {
		cj := t.C().Max(m.C()).tick()
		return Timestamp{lj.ts & cj.ts}
	} else if latest.Equals(lj) {
		return t.tick()
	} else if latest.Equals(lm) {
		return m.tick()
	}
	return latest
}
