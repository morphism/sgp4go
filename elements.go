package sgp4go

type Elements struct {
	// *    em          - eccentricity

	Eccentricity float64

	// *    argpm       - argument of perigee

	ArgOfPerigree float64

	// *    inclm       - inclination

	Inclination float64

	// *    mm          - mean anomaly

	MeanAnomaly float64

	// ?
	// *    n          - mean motion

	MeanMotion float64

	// *    nodem       - right ascension of ascending node

	RightAscension float64
}

func (tle *TLE) Elements() Elements {
	return Elements{
		Eccentricity:   tle.ecc,
		ArgOfPerigree:  tle.argpDeg,
		Inclination:    tle.incDeg,
		MeanAnomaly:    tle.maDeg,
		MeanMotion:     tle.n,
		RightAscension: tle.raanDeg,
	}
}
