package main

import (
	"encoding/json"
	"fmt"
	"github.com/fluepke/esptool/esp32"
	"os"
	"strconv"
	"strings"
)

type DeviceInfo struct {
	ChipType   string
	Revision   string
	Features   []string
	MacAddress string
	Partitions esp32.PartitionList
}

func (d *DeviceInfo) String() string {
	builder := &strings.Builder{}
	fmt.Fprint(builder, underline(bold(("Chip Information"))))
	fmt.Fprint(builder, "\n")
	fmt.Fprintf(builder, "%s: %s\n", bold("Chip Type"), d.ChipType)
	fmt.Fprintf(builder, "%s: %s\n", bold("Revision"), d.Revision)
	fmt.Fprintf(builder, "%s: %s\n", bold("MAC"), d.MacAddress)
	fmt.Fprintf(builder, "%s: %s\n", bold("Features"), strings.Join(d.Features, ", "))
	fmt.Fprintln(builder, bold("Partition Table"))
	if d.Partitions != nil {
		fmt.Fprint(builder, d.Partitions.String())
	} else {
		fmt.Fprint(builder, "** invalid **")
	}
	return builder.String()
}

func infoCommand(jsonOutput bool, esp32 *esp32.ESP32ROM) error {
	macAddress, err := esp32.GetChipMAC()
	if err != nil {
		return fmt.Errorf("Could not retrieve MAC address: %s", err.Error())
	}

	description, err := esp32.GetChipDescription()
	if err != nil {
		return fmt.Errorf("Could not retrieve chip description: %s", err.Error())
	}

	features, err := esp32.GetFeatures()
	if err != nil {
		return fmt.Errorf("Could not retrieve chip features: %s", err.Error())
	}

	featureList := make([]string, 0)
	for feature, status := range features {
		if status {
			featureList = append(featureList, feature.String())
		}
	}

	deviceInfo := &DeviceInfo{
		ChipType:   description.ChipType.String(),
		Revision:   strconv.Itoa(int(description.Revision)),
		Features:   featureList,
		MacAddress: macAddress,
	}

	partitionList, err := esp32.ReadPartitionList()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	if err == nil {
		deviceInfo.Partitions = partitionList
	}

	if jsonOutput {
		prettyJson, err := json.MarshalIndent(deviceInfo, "", "  ")

		if err != nil {
			return fmt.Errorf("Could not generate JSON outputs: %s", err.Error())
		}
		_, err = os.Stdout.Write(prettyJson)
		return err
	}

	_, err = fmt.Println(deviceInfo.String())

	return err
}
