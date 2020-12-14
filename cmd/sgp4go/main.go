// package main is a command-line program that reads TLEs on stdin and
// writes propagation data to stdout.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/morphism/sgp4go"

	sat "github.com/jsmorph/go-satellite"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var (
		ts = func(t time.Time) string {
			return t.Format(time.RFC3339Nano)
		}

		from     = flag.String("from", ts(time.Now()), "Propagation start time")
		to       = flag.String("to", ts(time.Now().Add(time.Minute)), "Propagation end time")
		interval = flag.Duration("interval", 6*time.Second, "Propagation end time")
	)

	flag.Parse()

	t0, err := time.Parse(time.RFC3339Nano, *from)
	if err != nil {
		return err
	}
	t1, err := time.Parse(time.RFC3339Nano, *to)
	if err != nil {
		return err
	}

	in := bufio.NewReader(os.Stdin)
	err = DoTLEs(in, 3, func(lines []string) error {
		t, err := sgp4go.NewTLE(lines[1], lines[2])
		if err != nil {
			return err
		}

		return Prop(t, t0, t1, *interval)
	})
	if err != nil {
		return err
	}

	return nil
}

func TimeToGST(t time.Time) (float64, float64) {
	var (
		y   = t.Year()
		m   = int(t.Month())
		d   = t.Day()
		h   = t.Hour()
		min = t.Minute()
		sec = t.Second()
		ns  = t.Nanosecond()
	)

	return sat.GSTimeFromDateNano(y, m, d, h, min, sec, ns)
}

type LatLonAlt struct {
	Lat, Lon, Alt float64
}

// ECIToLLA converts ECI coordinates to latitude, longitude, and
// altitude (km).
func ECIToLLA(t time.Time, p sgp4go.Vect) (*LatLonAlt, error) {

	gmst, _ := TimeToGST(t)

	x := sat.Vector3{
		X: float64(p.X),
		Y: float64(p.Y),
		Z: float64(p.Z),
	}

	// sat.ECIToLLA is very slow.
	alt, _, ll := sat.ECIToLLA(x, gmst)

	d, err := sat.LatLongDeg(ll)
	if err != nil {
		return nil, err
	}

	return &LatLonAlt{
		Lat: d.Latitude,
		Lon: d.Longitude,
		Alt: alt,
	}, nil
}

// Prop propagates over the given time range.
func Prop(o *sgp4go.TLE, from, to time.Time, interval time.Duration) error {

	for t := from; t.Before(to); t = t.Add(interval) {
		s, err := o.Prop(t)
		if err != nil {
			return err
		}

		lla, err := ECIToLLA(t, s.ECI)
		if err != nil {
			return err
		}

		m := map[string]interface{}{
			"Norad": o.NoradCatNum(),
			"At":    t,
			"State": s,
			"LLA":   lla,
		}
		js, err := json.Marshal(&m)
		if err != nil {
			log.Fatalf("prop json.Marshal error %s on %#v", err, m)
		}
		fmt.Printf("%s\n", js)
	}

	return nil
}

// DoTLEs iterates over TLE line groups.
func DoTLEs(r *bufio.Reader, group int, f func(lines []string) error) error {

	var (
		tle []string
		i   int
	)
	for {
		line, err := r.ReadString('\n')
		line = strings.TrimRight(line, "\n\r")
		if 0 < len(line) {
			if tle == nil {
				tle = make([]string, group)
			}
			tle[i%group] = line
			if (i+1)%group == 0 {
				if err := f(tle); err != nil {
					return err
				}
				tle = nil
			}
			i++
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}
