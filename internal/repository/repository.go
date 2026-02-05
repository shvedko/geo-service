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
