package main

import (
	"fmt"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func readPcap(filename string, filter string) (*pcap.Handle, error) {
	fmt.Printf("\nOpening file \"%s\" for parse", filename)
	file, err := pcap.OpenOffline(filename)
	if err != nil {
		err = fmt.Errorf("error opening PCAP file %v", err)
		return nil, err
	}
	if len(filter) > 0 {
		err = file.SetBPFFilter(filter)
		if err != nil {
			err = fmt.Errorf("error setting filter %v", err)
			return nil, err
		}
	}
	return file, nil
}
func openTsWrite(filename string) (*os.File, error) {
	fmt.Printf("\nCreating file \"%s\" ", filename)
	f, err := os.Create(filename)
	if err != nil {
		err = fmt.Errorf("error opening tsfile %v", err)
		return nil, err
	}
	return f, err
}

func main() {
	var pcapfile string
	var filter string
	var tsfile string
	argsWithoutProg := os.Args[1:]
	for i, arg := range argsWithoutProg {
		fmt.Printf("\n %d:\"%v\"", i, arg)
	}
	if len(argsWithoutProg) >= 3 {
		pcapfile = argsWithoutProg[0]
		filter = argsWithoutProg[1]
		tsfile = argsWithoutProg[2]
	} else {
		fmt.Printf("\nUsage : pcapfile filter tsfile")
		return
	}
	prdr, err := readPcap(pcapfile, filter)
	if err != nil {
		fmt.Printf("\n%v", err)
	}
	tswtr, err := openTsWrite(tsfile)
	if err != nil {
		fmt.Printf("\n%v", err)
	}
	packetSource := gopacket.NewPacketSource(prdr, prdr.LinkType())
	pkts := 0
	for {
		packet, err := packetSource.NextPacket()
		if err != nil {
			break
		}
		applicationLayer := packet.ApplicationLayer()
		if applicationLayer != nil {
			payload := applicationLayer.Payload()
			if len(payload) == (188*7)+12 {
				//Trim them off
				payload = payload[12:]
			}
			_, err := tswtr.Write(payload)
			if err != nil {
				fmt.Printf("\nError Writing %v", err)
			}
			pkts += int(len(payload) / 188)
		}
	}
	fmt.Printf("\nCompleted : %d pkts written", pkts)
}
