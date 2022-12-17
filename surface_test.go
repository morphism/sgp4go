package sgp4go

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func BenchmarkProp(b *testing.B) {
	var (
		ms       = time.Now().UTC().UnixNano() / 1000 / 1000 // ToDo: Use constant.
		line1    = "1 39132U PLANET   20016.08334491  .00000000  00000+0 -47542-3 0    07"
		line2    = "2 39132 064.8760 163.6520 0036285 284.0373 175.5769 15.07452065    00"
		tle, err = NewTLE(line1, line2)
	)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r, v, err := tle.PropUnixMillis(ms)
		if err != nil {
			b.Fatal(err)
		}
		if r[0] == 0 {
			b.Fatal(r)
		}
		if v[0] == 0 {
			b.Fatal(v)
		}
	}
}

var ISS = `ISS (ZARYA)             
1 25544U 98067A   20349.28181795  .00001103  00000-0  27992-4 0  9997
2 25544  51.6443 177.3570 0001731 128.2351  43.6939 15.49184106259930`

func Example() {
	var (
		tle = ISS

		lines = strings.Split(tle, "\n")
		o, _  = NewTLE(lines[1], lines[2])
		then  = "2020-12-12T12:00:00.000Z"
		t0, _ = time.Parse(time.RFC3339Nano, then)

		e, _ = o.Prop(t0)
	)

	fmt.Printf("%#v", e)

	// Output:
	// sgp4go.Ephemeris{V:sgp4go.Vect{X:-5.677866567405847, Y:-3.554868601056742, Z:3.696844342787944}, ECI:sgp4go.Vect{X:-4522.507182662008, Y:2857.518282965092, Z:-4201.949004177315}}
}

func getExample(t *testing.T) *TLE {
	var (
		tle = `ISS (ZARYA)             
1 25544U 98067A   20349.28181795  .00001103  00000-0  27992-4 0  9997
2 25544  51.6443 177.3570 0001731 128.2351  43.6939 15.49184106259930`

		lines = strings.Split(tle, "\n")
		o, err  = NewTLE(lines[1], lines[2])
	)

	if err != nil {
		t.Fatal(err)
	}

	return o
}

func TestSemiMajorAxis(t *testing.T) {
	var (
		alt = 408.0
		earthRadius = 6371.0
		want = alt + earthRadius
		o = getExample(t)
		m = o.SemiMajorAxisMeters()
		got = m / 1000
		epsilon = 20.0 // km
	)

	if  got < want - epsilon || want + epsilon < got {
		t.Fatal(got -want)
	}
}

func TestLines(t *testing.T) {
	var (
		tle = ISS
		lines = strings.Split(tle, "\n")
		o, err  = NewTLE(lines[1], lines[2])
	)

	if err != nil {
		t.Fatal(err)
	}
	line1, line2 := o.Lines()

	if line1 != lines[1] {
		t.Fatal(line1)
	}
	
	if line2 != lines[2] {
		t.Fatal(line2)
	}
}

func TestEqual(t *testing.T) {
	var (
		x = getExample(t)
		y = getExample(t)
	)

	t.Run("same", func(t *testing.T) {
		if !x.EqualValues(y) {
			t.Fatal(false)
		}
	})
	
	t.Run("different", func(t *testing.T) {
		y.epoch++
		if x.EqualValues(y) {
			t.Fatal(true)
		}
	})
}
