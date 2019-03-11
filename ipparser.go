package ipparser //ipv4parser

import (
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

/*IPv4Addr its a full IPv4 addres representation*/
type IPv4Addr struct {
	AddrIP net.IP
	IntIP  int64
}

// IPtoInt convert IPv4 addres to its integer representation
func IPtoInt(ipAddr string) int64 {
	strIP := ipAddr
	binIP := ""
	chunks := strings.Split(strIP, ".")

	for _, elem := range chunks {
		intElem, _ := strconv.Atoi(elem)
		binIPOctet := strconv.FormatInt(int64(intElem), 2)
		paddedOctet := strings.Repeat("0", 8-len(binIPOctet)) + binIPOctet
		binIP += paddedOctet
	}

	convToInt, _ := strconv.ParseInt(binIP, 2, 64)
	return convToInt
}

// FromInt convert int64 representation of IPv4 address to its string view
func FromInt(ipInt int64) string {
	binIP := strconv.FormatInt(ipInt, 2)
	binIP = strings.Repeat("0", 32-len(binIP)) + binIP
	convertedIP := ""
	for i, j := 0, 8; i <= 24; i += 8 {
		oct, _ := strconv.ParseInt(binIP[i:j], 2, 64)
		convertedIP += strconv.Itoa(int(oct)) + "."
		if j < 32 {
			j += 8
		}
	}
	return convertedIP[0 : len(convertedIP)-1]
}

// IPv4create fills fields of IPv4Addr struct with appropriate values
func (ipAddr *IPv4Addr) IPv4create(ip string) {
	splited := strings.Split(ip, ".")
	if func() bool {
		re, _ := regexp.Compile(`(\d){1,3}\.(\d){1,3}\.(\d){1,3}\.(\d){1,3}`)
		if re.Match([]byte(ip)) {
			for _, val := range splited {
				ival, err := strconv.Atoi(val)
				if err != nil {
					log.Fatal("Error parsing IP")
				}
				if ival > 255 || ival < 0 {
					return false
				}
			}
		} else {
			return false
		}
		return true
	}() {
		x := make([]byte, 4)
		for pos := range splited {
			val, _ := strconv.Atoi(splited[pos])
			x[pos] = byte(val)
		}
		ipAddr.AddrIP = net.IPv4(x[0], x[1], x[2], x[3])
		ipAddr.IntIP = IPtoInt(ip)
	}
}
