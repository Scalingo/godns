package main

import "strings"

func IsZonesFilteringEnabled() bool {
	return len(settings.Zones) > 0
}

func ZoneConfigForDomain(name string) (ZoneSettings, bool) {
	for _, zone := range settings.Zones {
		if strings.HasSuffix(name, zone.Name) {
			return zone, true
		}
	}
	return ZoneSettings{}, false
}
