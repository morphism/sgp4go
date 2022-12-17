package sgp4go

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Error adds error behavior to SGP4 status codes.
type Error int

// Error makes Error an error.
//
// Also see HasDecayed().
func (e Error) Error() string {
	var msg string
	switch int(e) {
	case 1:
		msg = "mean elements, ecc >= 1.0 or ecc < -0.001 or a < 0.95 er"
	case 2:
		msg = "mean motion less than 0.0"
	case 3:
		msg = "pert elements, ecc < 0.0  or  ecc > 1.0"
	case 4:
		msg = "semi-latus rectum < 0.0"
	case 5:
		msg = "epoch elements are sub-orbital"
	case 6:
		// See HasDecayed().
		msg = decayError
	default:
		msg = "NA"
	}
	return fmt.Sprintf("code=%d: %s", e, msg)
}

// decayError is a magic string used to detected that situation.
const decayError = "satellite has decayed"

// HasDecayed determines if the given error indicates the object has
// decayed.
func HasDecayed(e error) bool {
	// Too much trouble to try to check type (in the face of
	// fmt.wrapError, etc).
	if e == nil {
		return false
	}
	return strings.Contains(e.Error(), decayError)
}

// PropUnixMillis attempts to propagate a the given time in Unix
// milliseconds.
//
// Also see Prop().
func (tle *TLE) PropUnixMillis(ms int64) ([]float64, []float64, error) {
	// ToDo: Support higher time resolution.
	var (
		r = make([]float64, 3)
		v = make([]float64, 3)
	)

	tle.Lock()
	tle.Rec.error = 0
	getRVForDate(tle, ms, (*float64)(&r[0]), (*float64)(&v[0]))
	e := tle.sgp4Error
	tle.sgp4Error = 0
	tle.Unlock()

	if e != 0 {
		return nil, nil, fmt.Errorf("SGP4 error at ms=%d: %w", ms, Error(e))
	}
	return r, v, nil
}

// PropForMins propagates to the given minutes.
//
// This method calls simpler SGP4 functions, which is good for verification.
func (tle *TLE) PropForMins(mins float64) ([]float64, []float64, error) {
	var (
		r = make([]float64, 3)
		v = make([]float64, 3)
	)

	tle.Lock()
	tle.Rec.error = 0
	getRV(tle, mins, (*float64)(&r[0]), (*float64)(&v[0]))
	e := tle.sgp4Error
	tle.Rec.error = 0
	tle.Unlock()

	if e != 0 {
		return nil, nil, fmt.Errorf("SGP4 error at mins=%f: %w", mins, Error(e))
	}
	return r, v, nil
}

// NewTLE constructs a new TLE (which can be propagated).
//
// Also see Set().
func NewTLE(line1, line2 string) (*TLE, error) {
	tle := &TLE{}
	bs1 := []byte(line1)
	bs2 := []byte(line2)
	parseLines(tle, (*byte)(&bs1[0]), (*byte)(&bs2[0]))
	// ToDo: Detect and report errors!
	return tle, nil
}

// Set allows the caller to provide high-precision values than what a
// TLE can perhaps provide; however, this code has not (yet) been
// tested with respect to this additional precision.
func (tle *TLE) Set(epoch time.Time, mm1, mm2, bstar, incl, ra, ecc, aop, anom, mm, rev float64) {

	if !epoch.IsZero() {
		tle.epoch = epoch.UnixNano() / 1000_000
	}
	if mm1 != 0 {
		tle.ndot = mm1
	}

	if mm2 != 0 {
		tle.nddot = mm2
	}

	if bstar != 0 {
		tle.bstar = bstar
	}

	if incl != 0 {
		tle.incDeg = incl
	}

	if ra != 0 {
		tle.raanDeg = ra
	}

	if ecc != 0 {
		tle.ecc = ecc
	}

	if aop != 0 {
		tle.argpDeg = aop
	}

	if anom != 0 {
		tle.maDeg = anom
	}

	if mm != 0 {
		tle.n = mm
	}

	if rev != 0 {
		tle.revnum = int64(rev)
	}

	setValsToRec(tle, &(*tle).Rec)
}

// Vect is a 3-vector.
type Vect struct {
	X, Y, Z float64
}

// Ephemeris represents position and velocity.
type Ephemeris struct {
	// V is velocity in m/sec.
	V Vect

	// ECI is position in Earth-Centered Intertial coordinates.
	ECI Vect
}

// Prop propagates the given TLE.
//
// Currently the resolution is only milliseconds.
func (o *TLE) Prop(t time.Time) (Ephemeris, error) {
	// ToDo: Increase resolution.
	p, v, err := o.PropUnixMillis(t.UnixNano() / 1000 / 1000)
	var e Ephemeris
	if err == nil {
		e = Ephemeris{
			ECI: Vect{float64(p[0]), float64(p[1]), float64(p[2])},
			V:   Vect{float64(v[0]), float64(v[1]), float64(v[2])},
		}
	}
	return e, err
}

// NoradCatNum returns the NORAD catalog number of the TLE.
func (o *TLE) NoradCatNum() int {
	return int(o.objectNum)
}

// SemiMajorAxis returns what you would expect (hopefully).
func (tle *TLE) SemiMajorAxisMeters() float64 {
	var (
		u      = 3.986004418e14
		secs   = float64(24 * 60 * 60)
		mm = tle.n * 2*math.Pi/secs
	)

	return math.Pow(u, 1.0/3) / math.Pow(mm, 2.0/3)
}
