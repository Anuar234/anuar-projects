package handlers

func addToRecent(city string) {
	for _, c := range RecentCities {
		if c == city {
			return
		}
	}
	RecentCities = append(RecentCities, city)
}
