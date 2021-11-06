package deviceid

import (
	"encoding/hex"
	"encoding/json"
	"net"
	"os"
	"polysdk/consts"
	"polysdk/internal/hash"
	"sort"
	"strings"
)

// EnvPolyDeveceID id environment var of poly device id
const EnvPolyDeveceID = consts.EnvPolyDeveceID

// MyHexDeviceID get my poly device id as hex string.
// It was hashed by hostName+macAddrs.
// It use env hex value "ENV_POLY_DEVICE_ID" firstly.
func MyHexDeviceID() string {
	return hash.HexString(MyDeviceID())
}

// MyDeviceID get my poly device id as hash bytes.
// It was hashed by hostName+macAddrs.
// It use env hex value "ENV_POLY_DEVICE_ID" firstly.
func MyDeviceID() []byte {
	return myDeviceId()
}

//------------------------------------------------------------------------------

func myDeviceId() []byte {
	deviceID := os.Getenv(EnvPolyDeveceID)
	if deviceID != "" && len(deviceID) == hash.DefaultSize()*2 {
		if b, err := hex.DecodeString(deviceID); err != nil {
			return b
		}
	}

	var strList []string
	if hostName := getHostName(); hostName != "" {
		strList = append(strList, hostName)
	}
	if macAddrs := getMacAddrs(); macAddrs != "" {
		strList = append(strList, macAddrs)
	}

	device := strings.Join(strList, ";")
	//println(device)
	hexID := hash.Default(nil, device)
	return hexID
}

func getMacAddrs() string {
	var macAddrs []string
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("<unknown mac address>: " + err.Error())
		return ""
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr.String()
		// ignore virtual net addr like "00:50:56:c0:00:01"
		if mac != "" && !strings.HasPrefix(mac, "00:") {
			macAddrs = append(macAddrs, mac)
		}
	}
	sort.Strings(macAddrs) //mac addrs
	debugShowStrings(macAddrs)
	return strings.Join(macAddrs, ",")
}

func debugShowStrings(strs []string) {
	if debug := false; !debug {
		return // do nothing
	}
	if b, err := json.MarshalIndent(strs, "", "  "); err == nil {
		println(string(b))
	}
}

func getHostName() string {
	if hostName, err := os.Hostname(); err == nil {
		return hostName
	}
	return ""
}
