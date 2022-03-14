package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/miekg/dns"
)

var (
	errNoResponse error = errors.New("No response")
	errNoSOA error = errors.New("Did not receive a SOA record")
)

func exitUnknown(message string) {
	fmt.Printf("UNKNOWN - %s\n", message)
	os.Exit(3)
}

func exitCritical(message string) {
	fmt.Printf("CRITICAL - %s\n", message)
	os.Exit(2)
}

func main() {
	zone := flag.String("zone", "", "The zone to check")
	master := flag.String("master", "", "IP of the master server")
	rawslaves := flag.String("slaves", "", "IP of the slaves")
	flag.Parse()

	if len(*zone) == 0 {
		exitUnknown("No zone given")
	} else if len(*master) == 0 {
		exitUnknown("No master given")
	} else if len(*rawslaves) == 0 {
		exitUnknown("No slaves given")
	}

	slaves := strings.Split(*rawslaves, ",")

	masterSerial, masterError := fetchSerialFromServer(*zone, *master)

	if masterError != nil {
		exitCritical("Error querying master: " + masterError.Error())
	}

	exitCode := 0
	messages := make([]string, 0)

	for _, slave := range slaves {
		serial, err := fetchSerialFromServer(*zone, slave)

		if err != nil {
			messages = append(messages, fmt.Sprintf("%s query error: %s", slave, err.Error()))
			exitCode = 2
		} else if masterSerial != serial {
			messages = append(messages, fmt.Sprintf("%s serial mismatch (%d)", slave, serial))
			exitCode = 2
		} else {
			messages = append(messages, fmt.Sprintf("%s is in sync", slave))
		}
	}

	label := "OK"

	if exitCode == 3 {
		label = "CRITICAL"
	}

	fmt.Printf("%s - Zone: %s, Master serial: %d, Slaves: %s\n", label, *zone, masterSerial, strings.Join(messages, ", "))
	os.Exit(exitCode)
}

func fetchSerialFromServer(zone, server string) (uint32, error) {
	for i := 0; i < 3; i++ {
		message := &dns.Msg{}
		message.SetQuestion(dns.Fqdn(zone), dns.TypeSOA)

		result, err := dns.Exchange(message, server+":53")

		if err == nil && len(result.Answer) > 0 {
			if soa, ok := result.Answer[0].(*dns.SOA); ok {
				return soa.Serial, nil
			} else {
				return 0, errNoSOA
			}
		}
	}

	return 0, errNoResponse
}
