package engine

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/Microsoft/oe-engine/pkg/api"
)

// setPropertiesDefaults for the container Properties, returns true if certs are generated
func setPropertiesDefaults(oe *api.OpenEnclave, isUpgrade bool) {
	if len(oe.PackageBaseURL) == 0 {
		oe.PackageBaseURL = DefaultPackageBaseURL
	}
	if oe.Properties.MasterProfile == nil {
		oe.Properties.MasterProfile = &api.MasterProfile{}
	}
	if len(oe.Properties.MasterProfile.OSImageName) == 0 {
		oe.Properties.MasterProfile.OSImageName = api.OsImageDefault
	}
	setMasterNetworkDefaults(oe.Properties, isUpgrade)
}

// SetMasterNetworkDefaults for masters
func setMasterNetworkDefaults(a *api.Properties, isUpgrade bool) {
	if a.MasterProfile == nil {
		return
	}

	if !a.MasterProfile.IsCustomVNET() {
		a.MasterProfile.Subnet = DefaultSubnet
		// StaticIP is not reset if it is upgrade and some value already exists
		if !isUpgrade || len(a.MasterProfile.StaticIP) == 0 {
			a.MasterProfile.StaticIP = DefaultStaticIP
		}
	}

	if a.MasterProfile.HTTPSourceAddressPrefix == "" {
		a.MasterProfile.HTTPSourceAddressPrefix = "*"
	}
}

func combineValues(inputs ...string) string {
	valueMap := make(map[string]string)
	for _, input := range inputs {
		applyValueStringToMap(valueMap, input)
	}
	return mapToString(valueMap)
}

func applyValueStringToMap(valueMap map[string]string, input string) {
	values := strings.Split(input, ",")
	for index := 0; index < len(values); index++ {
		// trim spaces (e.g. if the input was "foo=true, bar=true" - we want to drop the space after the comma)
		value := strings.Trim(values[index], " ")
		valueParts := strings.Split(value, "=")
		if len(valueParts) == 2 {
			valueMap[valueParts[0]] = valueParts[1]
		}
	}
}

func mapToString(valueMap map[string]string) string {
	// Order by key for consistency
	keys := []string{}
	for key := range valueMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	for _, key := range keys {
		buf.WriteString(fmt.Sprintf("%s=%s,", key, valueMap[key]))
	}
	return strings.TrimSuffix(buf.String(), ",")
}
