# SGP4 C transpiled to Go

_Warning: This code is experimental; please see the
[LICENSE](LICENSE)._

An
[SGP4](http://celestrak.com/publications/AIAA/2006-6753/AIAA-2006-6753-Rev2.pdf)
implementation based on
[`github.com/aholinch/sgp4/tree/master/src/c`](https://github.com/aholinch/sgp4/tree/master/src/c),
which ultimately originated with [David
Vallado](https://celestrak.com/software/vallado-sw.php).

```
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
```

This implementation is a hand-edited transpilation of [C
sources](https://github.com/aholinch/sgp4/tree/master/src/c) to
[Go](https://golang.org/) by
[`c2go`](https://github.com/elliotchance/c2go) (version v0.25.9
Dubnium 2018-12-30), and the emitted code was edited by hand.  The
original C implementation test suite was included in this process.

The substantive edits (of Go sources emitted by the transpiler) were
the use of 64-bit integers to address at least one 32-bit overflow and
using floating point constants instead of naked integer constants when
the entire expression consisted of the latter with some division.
Example: `var x2o3 float64 = 2 / 3` was edited to be `var x2o3 float64
= 2.0 / 3.0`.

With those changes, the original (transpiled and hand-edited) tests
almost all pass.  The one exception is for NORAD ID
[20413](https://www.n2yo.com/satellite/?s=20413) at `mins=1844335`,
where 1e-07 < _rdist_ < 1e-06.  The tests have been edited to tolerate
_rdist_ < 1e-06 rather than demand _rdist_ < 1e-07.

## Usage

The main constructor is `NewTLE()`, and the primary method is
`Prop()`.  See `ExampleProp()` in
[`surface_test.go`](surface_test.go).

See [the documentation](https://godoc.org/github.com/morphism/sgp4go)
for details.

The example command-line program [`sgp4go`](cmd/sgp4go) reads TLEs
from `stdin` and writes propagation data to `stdout`. See
[`test.sh`](test.sh) for an example invocation.

Some `sgp4go` executables are available
[here](https://github.com/morphism/sgp4go/releases).

## References

1. [`space-trace.org`](https://www.space-track.org/)'s [SGP4 binaries
   and example code](https://www.space-track.org/documentation#/sgp4).

1. [SGP4 in
   AIAA-2006-6753-Rev2](http://celestrak.com/publications/AIAA/2006-6753/AIAA-2006-6753-Rev2.pdf).

1. [Fundamentals of Astrodynamics and
   Applications](https://celestrak.com/software/vallado-sw.php) (also
   [at
   Amazon](https://www.amazon.com/Fundamentals-Astrodynamics-Applications-Technology-Library/dp/1881883183/ref=pd_lpo_14_t_0/140-1650425-3257455)).
   
