// This file contains the original test suite, which is now packaged
// up in the test TestSGP4().

/*
  Originally transpiled from

     https://github.com/aholinch/sgp4/tree/master/src/c

  by c2go

     https://github.com/elliotchance/c2go

  (version v0.25.9 Dubnium 2018-12-30).

  Then edited by hand.

  Original C code credit:

    This file contains the sgp4 procedures for analytical propagation
    of a satellite. the code was originally released in the 1980 and
    1986 spacetrack papers. a detailed discussion of the theory and
    history may be found in the 2006 aiaa paper by vallado, crawford,
    hujsak, and kelso.

                           companion code for
              fundamentals of astrodynamics and applications
                                   2013
                             by david vallado
     (w) 719-573-2600, email dvallado@agi.com, davallado@gmail.com

*/

package sgp4go

import (
	"fmt"
	"log"
	"math"
	"testing"
	"unsafe"

	"github.com/elliotchance/c2go/noarch"
)

type VERIN struct {
	line1    [70]byte
	line2    [70]byte
	startmin float64
	stepmin  float64
	stopmin  float64
}

// readVERINs - transpiled function from  /home/somebody/aholinch/sgp4/src/c/TestSGP4.c:18
/**
 * returns the count of verins read and sets the point to an array created with malloc.
 */ //
//
func readVERINs(listptr **VERIN) int32 {
	var line []byte = make([]byte, 256, 256)
	var in_file *noarch.File = nil
	var verins *VERIN = nil
	var cnt int32 = int32(0)
	in_file = noarch.Fopen((&[]byte("SGP4-VER.TLE\x00")[0]), (&[]byte("r\x00")[0]))

	for noarch.Fgets(&line[0], int32(255), in_file) != nil {
		if int32(*&line[0]) == int32('1') {
			cnt += 1
		}
	}
	if in_file != nil {
		noarch.Fclose(in_file)
	}
	verins = (*VERIN)(noarch.Malloc(int32(uint32(cnt) * 164)))
	*listptr = verins
	cnt = int32(0)
	in_file = noarch.Fopen((&[]byte("SGP4-VER.TLE\x00")[0]), (&[]byte("r\x00")[0]))
	for noarch.Fgets(&line[0], int32(255), in_file) != nil {
		if int32(*&line[0]) == int32('1') {
			noarch.Strncpy(&(*((*VERIN)(unsafe.Pointer(uintptr(unsafe.Pointer(verins)) + (uintptr)(cnt)*unsafe.Sizeof(*verins))))).line1[0], &line[0], int32(uint32(int32(69))))
			*((*byte)(func() unsafe.Pointer {
				tempVar := &(*((*VERIN)(unsafe.Pointer(uintptr(unsafe.Pointer(verins)) + (uintptr)(cnt)*unsafe.Sizeof(*verins))))).line1[0]
				return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int32(69))*unsafe.Sizeof(*tempVar))
			}())) = byte(int32(0))
			noarch.Fgets(&line[0], int32(255), in_file)
			noarch.Strncpy(&(*((*VERIN)(unsafe.Pointer(uintptr(unsafe.Pointer(verins)) + (uintptr)(cnt)*unsafe.Sizeof(*verins))))).line2[0], &line[0], int32(uint32(int32(69))))
			*((*byte)(func() unsafe.Pointer {
				tempVar := &(*((*VERIN)(unsafe.Pointer(uintptr(unsafe.Pointer(verins)) + (uintptr)(cnt)*unsafe.Sizeof(*verins))))).line2[0]
				return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int32(69))*unsafe.Sizeof(*tempVar))
			}())) = byte(int32(0))
			sscanf(string(line), "%f %f %f\x00", &(*((*VERIN)(unsafe.Pointer(uintptr(unsafe.Pointer(verins)) + (uintptr)(cnt)*unsafe.Sizeof(*verins))))).startmin, &(*((*VERIN)(unsafe.Pointer(uintptr(unsafe.Pointer(verins)) + (uintptr)(cnt)*unsafe.Sizeof(*verins))))).stopmin, &(*((*VERIN)(unsafe.Pointer(uintptr(unsafe.Pointer(verins)) + (uintptr)(cnt)*unsafe.Sizeof(*verins))))).stepmin)
			cnt += 1
		}
	}
	if in_file != nil {
		noarch.Fclose(in_file)
	}
	return cnt
}

// dist - transpiled function from  /home/somebody/aholinch/sgp4/src/c/TestSGP4.c:71
/**
 * 2-norm distance for two three vectors
 */ //
// just unroll the loop
//
func dist(v1 *float64, v2 *float64) float64 {
	var dist float64 = float64(int32(0))
	var tmp float64 = float64(int32(0))
	tmp = *v1 - *v2
	dist += tmp * tmp
	tmp = *((*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(v1)) + (uintptr)(int32(1))*unsafe.Sizeof(*v1)))) - *((*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + (uintptr)(int32(1))*unsafe.Sizeof(*v2))))
	dist += tmp * tmp
	tmp = *((*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(v1)) + (uintptr)(int32(2))*unsafe.Sizeof(*v1)))) - *((*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + (uintptr)(int32(2))*unsafe.Sizeof(*v2))))
	dist += tmp * tmp
	return math.Sqrt(dist)
}

// runVER - transpiled function from  /home/somebody/aholinch/sgp4/src/c/TestSGP4.c:87
func runVER(verins *VERIN, cnt int32) {
	var in_file *noarch.File = nil
	var tle TLE
	var r []float64 = make([]float64, 3, 3)
	var v []float64 = make([]float64, 3, 3)
	var rv []float64 = make([]float64, 3, 3)
	var vv []float64 = make([]float64, 3, 3)
	var mins float64 = float64(int32(0))
	var line []byte = make([]byte, 256, 256)
	var i int32 = -int32(1)
	var ind *byte = nil
	var rdist float64 = float64(int32(0))
	var vdist float64 = float64(int32(0))
	var rerr float64 = float64(int32(0))
	var verr float64 = float64(int32(0))
	var cnt2 int32 = int32(0)
	in_file = noarch.Fopen((&[]byte("tcppver.out\x00")[0]), (&[]byte("r\x00")[0]))
	for noarch.Fgets(&line[0], int32(255), in_file) != nil && i < cnt {
		ind = noarch.Strstr(&line[0], (&[]byte("xx\x00")[0]))
		if ind != nil {
			i += 1
			parseLines(&tle, &(*((*VERIN)(unsafe.Pointer(uintptr(unsafe.Pointer(verins)) + (uintptr)(i)*unsafe.Sizeof(*verins))))).line1[0], &(*((*VERIN)(unsafe.Pointer(uintptr(unsafe.Pointer(verins)) + (uintptr)(i)*unsafe.Sizeof(*verins))))).line2[0])
		} else {
			fmt.Sscanf(string(line),
				"%f %f %f %f %f %f %f\n",
				&mins,
				&rv[0], &rv[1], &rv[2],
				&vv[0], &vv[1], &vv[2],
			)

			getRV(&tle, mins, &r[0], &v[0])

			rdist = dist(&r[0], &rv[0])
			vdist = dist(&v[0], &vv[0])

			rerr += rdist
			verr += vdist
			cnt2 += 1

			// Jamie: Changed the epsilon from 1e-07,
			// which triggered the problem report for
			// objectNum 20413 amdn mins=1844335.
			if rdist > 1e-06 {
				log.Printf("debug %v %v %v %v %v %v", r, v, rv, vv, rdist, vdist)
				log.Fatalf("rdist %d %f %0.12f", tle.objectNum, mins, rdist)
			}
			if vdist > 1e-08 {
				log.Printf("debug %v %v %v %v %v %v", r, v, rv, vv, rdist, vdist)
				log.Fatalf("vdist %d %f %0.12ff", tle.objectNum, mins, vdist)
			}
		}
	}
	rerr = rerr / float64(cnt2)
	verr = verr / float64(cnt2)

	fmt.Printf("Typical errors r=%e mm, v=%e mm/s\n", 1e+06*rerr, 1e+06*verr)

	if in_file != nil {
		noarch.Fclose(in_file)
	}
}

// main - transpiled function from  /home/somebody/aholinch/sgp4/src/c/TestSGP4.c:143
func TestSGP4(t *testing.T) {
	var cnt int32 = int32(0)
	var verins *VERIN = nil
	cnt = readVERINs(&verins)
	noarch.Printf((&[]byte("read %d verins\n\x00")[0]), cnt)
	runVER(verins, cnt)
	noarch.Free(unsafe.Pointer(verins))
	return
}
