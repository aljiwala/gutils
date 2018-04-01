package netx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/context/ctxhttp"
)

const gmapsURL = `https://maps.googleapis.com/maps/api/geocode/json?sensors=false&address={{.Address}}`

var (
	// ErrNotFound ...
	ErrNotFound = errors.New("not found")
	// ErrTooManyResults ...
	ErrTooManyResults = errors.New("too many results")
)

type (
	// Coordinates ...
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	// Location ...
	Location struct {
		Address string
		Coordinates
	}

	// MapsLocation ...
	MapsLocation struct {
		Coordinates
	}

	// MapsGeometry ...
	MapsGeometry struct {
		Location MapsLocation `json:"location"`
	}

	// MapsResult ...
	MapsResult struct {
		FormattedAddress string       `json:"formatted_address"`
		Geometry         MapsGeometry `json:"geometry"`
	}

	// MapsResponse ...
	MapsResponse struct {
		Status  string       `json:"status"`
		Results []MapsResult `json:"results"`
	}

	// Net ...
	Net struct {
		Address   string
		Bitmask   uint8
		Mask      string
		Hostmask  string
		Broadcast string
		First     string
		Last      string
		Size      uint32
	}
)

// Get ...
func Get(ctx context.Context, address string) (Location, error) {
	var (
		loc  Location
		data MapsResponse
	)

	select {
	case <-ctx.Done():
		return loc, ctx.Err()
	default:
	}

	aURL := strings.Replace(gmapsURL, "{{.Address}}", url.QueryEscape(address), 1)
	resp, err := ctxhttp.Get(ctx, nil, aURL)
	if err != nil {
		return loc, errors.Wrapf(err, aURL)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		return loc, errors.Wrapf(err, aURL)
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return loc, errors.Wrapf(err, "decode")
	}

	switch data.Status {
	case "OK":
	case "ZERO_RESULTS":
		return loc, ErrNotFound
	default:
		return loc, errors.Wrapf(err, "status=%q", data.Status)
	}

	switch len(data.Results) {
	case 0:
		return loc, ErrNotFound
	case 1:
	default:
		return loc, ErrTooManyResults
	}

	result := data.Results[0]
	loc.Address = result.FormattedAddress
	loc.Coordinates.Latitude, loc.Coordinates.Longitude =
		result.Geometry.Location.Coordinates.Latitude,
		result.Geometry.Location.Coordinates.Longitude

	return loc, nil
}

// RemovePort removes the "port" part of an hostname.
func RemovePort(host string) string {
	shost, _, err := net.SplitHostPort(host)
	// Probably doesn't have a port, which is an error.
	if err != nil {
		return host
	}
	return shost
}

// LocalIP should return local IP address.
func LocalIP() (string, error) {
	addr, err := net.ResolveUDPAddr("udp", "1.2.3.4:1")
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	host, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return "", err
	}

	// host = "10.180.2.66"
	return host, nil
}

// LocalDNSName should return host name.
func LocalDNSName() (hostname string, err error) {
	var ip string
	ip, err = LocalIP()
	if err != nil {
		return
	}

	cmd := exec.Command("host", ip)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}

	tmp := out.String()
	arr := strings.Split(tmp, ".\n")

	if len(arr) > 1 {
		content := arr[0]
		arr = strings.Split(content, " ")
		return arr[len(arr)-1], nil
	}

	err = fmt.Errorf("parse host %s fail", ip)
	return
}

// IntranetIP get internal IP addr.
func IntranetIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // Interface down.
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // Loopback interface.
		}

		if strings.HasPrefix(iface.Name, "docker") ||
			strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // Not an IPv4 address.
			}

			return ip.String(), nil
		}
	}

	return "", errors.New("Are you connected to the network?")
}

// ExtranetIP should return external IP address.
func ExtranetIP() (ip string, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("Get external IP error: %v", p)
		} else if err != nil {
			err = errors.New("Get external IP error: " + err.Error())
		}
	}()

	resp, err := http.Get("http://pv.sohu.com/cityjson?ie=utf-8")
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	idx := bytes.Index(b, []byte(`"cip": "`))
	b = b[idx+len(`"cip": "`):]
	idx = bytes.Index(b, []byte(`"`))
	b = b[:idx]
	ip = string(b)

	return
}

// Atoi returns the uint32 representation of an ipv4 addr string value.
//
// Example:
//
//	Atoi("192.168.0.1")   // 3232235521
//
func Atoi(addr string) (sum uint32, err error) {
	if len(addr) > 15 {
		return sum, errors.New("addr too long")
	}

	octs := strings.Split(addr, ".")
	if len(octs) != 4 {
		return sum, errors.New("requires 4 octects")
	}

	for i := 0; i < 4; i++ {
		oct, err := strconv.ParseUint(octs[i], 10, 0)
		if err != nil {
			return sum, errors.New("bad octect " + octs[i])
		}
		sum += uint32(oct) << uint32((4-1-i)*8)
	}
	return sum, nil
}

// Itoa returns the string representation of an ipv4 addr uint32 value.
//
// Example:
//
//	Itoa(3232235521)  // "192.168.0.1"
//
func Itoa(addr uint32) string {
	var buf bytes.Buffer

	for i := 0; i < 4; i++ {
		oct := (addr >> uint32((4-1-i)*8)) & 0xff
		buf.WriteString(strconv.FormatUint(uint64(oct), 10))
		if i < 3 {
			buf.WriteByte('.')
		}
	}
	return buf.String()
}

// Not ...
// Example:
//
//	Not("0.0.255.255")  // "255.255.0.0"
//
func Not(addr string) (string, error) {
	i, err := Atoi(addr)
	return Itoa(i ^ 0xffffffff), err
}

// Or ...
// Example:
//
//	Or("0.0.1.1", "1.1.0.0")  // "1.1.1.1"
//
func Or(addra string, addrb string) (addr string, err error) {
	ia, err := Atoi(addra)
	if err != nil {
		return addr, err
	}

	ib, err := Atoi(addrb)
	if err != nil {
		return addr, err
	}

	return Itoa(ia | ib), err
}

// Xor ...
// Example:
//
//	Xor("0.255.255.255", "192.255.255.255")  // "192.0.0.0"
//
func Xor(addra string, addrb string) (addr string, err error) {
	ia, err := Atoi(addra)
	if err != nil {
		return addr, err
	}

	ib, err := Atoi(addrb)
	if err != nil {
		return addr, err
	}

	return Itoa(ia ^ ib), err
}

// Next ...
// Example:
//
//	Next("192.168.0.1")  // "192.168.0.2"
//
func Next(addr string) (string, error) {
	i, err := Atoi(addr)
	return Itoa(i + 1), err
}

// Prev ...
// Example:
//
//	Prev("192.168.0.1")  // "192.168.0.0"
//
func Prev(addr string) (string, error) {
	i, err := Atoi(addr)
	return Itoa(i - 1), err
}

// Network returns information for a netblock.
//
// Example:
//
//	Network("192.168.0.0/24")
//	// {
//	//	Address: "192.168.0.0",
//	//	Bitmask: 24,
//	//	Mask: "255.255.255.0",
//	//	Hostmask: "0.0.0.255",
//	//	Broadcast: "192.168.0.255",
//	//	First: "192.168.0.1",
//	//	Last: "192.168.0.254",
//	//	Size: 254,
//	// }
func Network(block string) (net Net, err error) {
	if len(block) > 18 {
		return net, errors.New("block too long")
	}

	list := strings.Split(block, "/")
	if len(list) != 2 {
		return net, errors.New("invalid block")
	}

	// address
	net.Address = list[0]

	// bitmask
	bitmask, err := strconv.ParseUint(list[1], 10, 0)
	if err != nil {
		return net, err
	}
	if bitmask&31 != bitmask {
		return net, errors.New("invalid bitmask")
	}
	net.Bitmask = uint8(bitmask)

	// mask
	net.Mask = Itoa(0xffffffff >> (32 - net.Bitmask) << (32 - net.Bitmask))
	net.Hostmask, err = Not(net.Mask)
	if err != nil {
		return net, err
	}

	// broadcast
	net.Broadcast, err = Or(net.Address, net.Hostmask)
	if err != nil {
		return net, err
	}

	// first
	addr, err := Xor(net.Hostmask, net.Broadcast)
	if err != nil {
		return net, err
	}

	net.First, err = Next(addr)
	if err != nil {
		return net, err
	}

	// last
	net.Last, err = Prev(net.Broadcast)
	if err != nil {
		return net, err
	}

	// size
	i, err := Atoi(net.Last)
	if err != nil {
		return net, err
	}

	j, err := Atoi(net.First)
	if err != nil {
		return net, err
	}

	net.Size = i - j + 1
	return net, nil
}
