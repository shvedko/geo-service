package repository

func (p AddPointParams) IsValid() bool {
	if p.Lon < -180 || p.Lon > 180 {
		return false
	}
	if p.Lat < -90 || p.Lat > 90 {
		return false
	}
	if p.Name == "" {
		return false
	}
	return true
}

func (p GetPointsFromBoxParams) IsValid() bool {
	if p.MinLon < -180 || p.MinLon > 180 {
		return false
	}
	if p.MinLat < -90 || p.MinLat > 90 {
		return false
	}
	if p.MaxLon < -180 || p.MaxLon > 180 {
		return false
	}
	if p.MaxLat < -90 || p.MaxLat > 90 {
		return false
	}
	if p.MinLon > p.MaxLon || p.MinLat > p.MaxLat {
		return false
	}
	return true
}
