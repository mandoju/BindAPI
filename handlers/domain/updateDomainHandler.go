package domain

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mandoju/BindAPI/utils"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type RequestDataExtractor struct {
	Address func(request *http.Request) string
	Secret  func(request *http.Request) string
	Domain  func(request *http.Request) string
}
type WebserviceResponse struct {
	Success  bool
	Message  string
	Domain   string
	Domains  []string
	Address  string
	AddrType string
}

var secret, _ = utils.GetDnsSecretKey()

func BuildWebserviceResponseFromRequest(r *http.Request, extractors RequestDataExtractor) WebserviceResponse {
	response := WebserviceResponse{}

	sharedSecret := extractors.Secret(r)
	fmt.Println(sharedSecret)
	fmt.Println(secret)
	response.Domains = strings.Split(extractors.Domain(r), ",")
	response.Address = extractors.Address(r)

	if sharedSecret != secret {
		log.Println(fmt.Sprintf("Invalid shared secret: %s", sharedSecret))
		response.Success = false
		response.Message = "Invalid Credentials"
		return response
	}

	for _, domain := range response.Domains {
		if domain == "" {
			response.Success = false
			response.Message = fmt.Sprintf("Domain not set")
			log.Println("Domain not set")
			return response
		}
	}

	// kept in the response for compatibility reasons
	response.Domain = strings.Join(response.Domains, ",")

	if ValidIP4(response.Address) {
		response.AddrType = "A"
	} else if ValidIP6(response.Address) {
		response.AddrType = "AAAA"
	} else {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)

		if err != nil {
			response.Success = false
			response.Message = fmt.Sprintf("%q is neither a valid IPv4 nor IPv6 address", r.RemoteAddr)
			log.Println(fmt.Sprintf("Invalid address: %q", r.RemoteAddr))
			return response
		}

		if ValidIP4(ip) {
			response.AddrType = "A"
			response.Address = ip
		} else if ValidIP6(ip) {
			response.AddrType = "AAAA"
			response.Address = ip
		} else {
			response.Success = false
			response.Message = fmt.Sprintf("%s is neither a valid IPv4 nor IPv6 address", response.Address)
			log.Println(fmt.Sprintf("Invalid address: %s", response.Address))
			return response
		}
	}

	response.Success = true

	return response
}
func UpdateDomainHandler(w http.ResponseWriter, r *http.Request) {
		extractor := RequestDataExtractor{
			Address: func(r *http.Request) string { return r.URL.Query().Get("addr") },
			Secret:  func(r *http.Request) string { return r.URL.Query().Get("secret") },
			Domain:  func(r *http.Request) string { return r.URL.Query().Get("domain") },
		}
	fmt.Println(extractor.Address(r))

	fmt.Println(extractor.Secret(r))
		response := BuildWebserviceResponseFromRequest(r, extractor)

		if response.Success == false {
			json.NewEncoder(w).Encode(response)
			return
		}

		for _, domain := range response.Domains {
			result := UpdateRecord(domain, response.Address, response.AddrType)

			if result != "" {
				response.Success = false
				response.Message = result

				json.NewEncoder(w).Encode(response)
				return
			}
		}

		response.Success = true
		response.Message = fmt.Sprintf("Updated %s record for %s to IP address %s", response.AddrType, response.Domain, response.Address)

		json.NewEncoder(w).Encode(response)
	}


func UpdateRecord(domain string, ipaddr string, addrType string) string {
	log.Println(fmt.Sprintf("%s record update request: %s -> %s", addrType, domain, ipaddr))

	f, err := ioutil.TempFile(os.TempDir(), "dyndns")
	if err != nil {
		return err.Error()
	}

	defer os.Remove(f.Name())
	w := bufio.NewWriter(f)

	w.WriteString(fmt.Sprintf("server %s\n", "18.220.253.183"))
	w.WriteString(fmt.Sprintf("zone %s\n", "example.com"))
	w.WriteString(fmt.Sprintf("update delete %s.%s %s\n", domain, "Domain", addrType))
	w.WriteString(fmt.Sprintf("update add %s.%s %v %s %s\n", domain, "Domain", "84600", addrType, ipaddr))
	w.WriteString("send\n")

	w.Flush()
	f.Close()

	cmd := exec.Command("C:\\Users\\jorge\\Desktop\\utilitários\\BIND9.14.10.x64\\nsupdate.exe", f.Name())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return err.Error() + ": " + stderr.String()
	}
	fmt.Println(out.String())
	return out.String()
}
func ValidIP4(ipAddress string) bool {
	testInput := net.ParseIP(ipAddress)
	if testInput == nil {
		return false
	}

	return (testInput.To4() != nil)
}

func ValidIP6(ip6Address string) bool {
	testInputIP6 := net.ParseIP(ip6Address)
	if testInputIP6 == nil {
		return false
	}

	return (testInputIP6.To16() != nil)
}