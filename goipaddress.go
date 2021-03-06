package goipaddress

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var errInvalidIP = errors.New("invalid IP address given")

/*IPv4Range contains the range of IP addresses*/
type IPv4Range []string

/*IPv4Addr its a full IPv4 addres representation*/
type IPv4Addr struct {
	AddrIP net.IP
	IntIP  int64
}

/*IPv4Network represents an IPs network*/
type IPv4Network struct {
	AddrIP  string
	IPrange IPv4Range
}

// ToInt convert IPv4 address to its integer representation
func ToInt(ipAddr string) int64 {
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

// IPv4Create fills fields of IPv4Addr struct with appropriate values
func IPv4Create(ipAddr string) (IPv4Addr, error) {
	var ip IPv4Addr
	if isValid(ipAddr) {
		return IPv4Addr{net.ParseIP(ipAddr), ToInt(ipAddr)}, nil
	}
	return ip, errInvalidIP
}

/*
IPv4NetworkCreate creates new IPv4Network instance
This function can handle IPs like:
192.*.1.*, 192.1-23.1.8-10, 192.1-23.*.8 and CIDR notation,
*/
func IPv4NetworkCreate(ipAddr string) (IPv4Network, error) {
	if isValid(ipAddr) {
		return IPv4Network{ipAddr, parseAddr(ipAddr)}, nil
	}
	return IPv4Network{}, errInvalidIP
}

//Check ipAddr is a standard IP addr like 192.168.1.1
func isValid(ipAddr string) bool {
	if !strings.Contains(ipAddr, "-") && !strings.Contains(ipAddr, "*") && !strings.Contains(ipAddr, "/") {
		re, _ := regexp.Compile(`^(\d){1,3}\.(\d){1,3}\.(\d){1,3}\.(\d){1,3}$`) //default ip regexp
		if !re.Match([]byte(ipAddr)) {
			return false
		}
	}
	if strings.Contains(ipAddr, "/") {
		cidr := strings.Split(ipAddr, "/")[1]
		cval, err := strconv.Atoi(cidr)
		if err != nil {
			return false
		}
		if cval < 1 || cval > 32 {
			return false
		}
	}
	for _, delim := range []string{"-", "*"} {
		ipAddr = strings.Replace(ipAddr, delim, ".", -1)
	}
	splited := strings.Split(ipAddr, ".")
	for _, val := range splited {
		ival, err := strconv.Atoi(val)
		if err != nil {
			continue
		}
		if ival > 255 || ival < 0 {
			return false
		}
	}
	return true
}

func parseAster(ipAddr string, storage *[]string) {
	if strings.Count(ipAddr, "*") > 1 {
		for i := 0; i <= 255; i++ {
			repl := strings.Replace(ipAddr, "*", strconv.Itoa(i), 1)
			parseAster(repl, storage)
		}
	} else {
		for i := 0; i <= 255; i++ {
			*storage = append(*storage, strings.Replace(ipAddr, "*", strconv.Itoa(i), 1))
		}
	}
}

func parseHyphen(ipAddr string, storage *[]string) {
	reg, _ := regexp.Compile(`(\d*)-(\d*)`)
	IPrange := fmt.Sprintf("%s", reg.Find([]byte(ipAddr)))
	repl := strings.Replace(ipAddr, IPrange, "@", -1)
	s := strings.Split(IPrange, "-")
	beg, _ := strconv.Atoi(s[0])
	end, _ := strconv.Atoi(s[1])
	if strings.Count(ipAddr, "-") > 1 {
		for ; beg <= end; beg++ {
			repl2 := strings.Replace(repl, "@", strconv.Itoa(beg), 1)
			parseHyphen(repl2, storage)
		}
	} else {
		for ; beg <= end; beg++ {
			*storage = append(*storage, strings.Replace(repl, "@", strconv.Itoa(beg), 1))
		}
	}
}

func parseCIDR(ipAddr string, storage *[]string) {
	inc := func(s net.IP) {
		for j := len(s) - 1; j >= 0; j-- {
			s[j]++
			if s[j] > 0 {
				break
			}
		}
	}
	ip, network, _ := net.ParseCIDR(ipAddr)
	for e := ip.Mask(network.Mask); network.Contains(e); inc(e) {
		*storage = append(*storage, e.String())
	}
}

func parseAddr(ipAddr string) []string {
	var aStore []string
	var bStore []string
	ast := strings.Contains(ipAddr, "*")
	hyp := strings.Contains(ipAddr, "-")
	cidr := strings.Contains(ipAddr, "/")
	switch {
	case ast && hyp:
		parseHyphen(ipAddr, &aStore)
		for _, elem := range aStore {
			parseAster(elem, &bStore)
		}
		return bStore
	case ast:
		parseAster(ipAddr, &aStore)
		return aStore
	case hyp:
		parseHyphen(ipAddr, &aStore)
		return aStore
	case cidr:
		parseCIDR(ipAddr, &aStore)
	}
	return aStore
}
