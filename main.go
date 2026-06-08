package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Iface struct {
	Name        string
	Path        string
	Options     map[string]string
	IPv4Address []string
}

func main() {
	etcnetPath := "testdata/ifaces"

	entries, err := os.ReadDir(etcnetPath)
	if err != nil {
		fmt.Println("error reading ifaces:", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		ifaceName := entry.Name()
		ifacePath := filepath.Join(etcnetPath, ifaceName)
		optionsPath := filepath.Join(ifacePath, "options")

		data, err := os.ReadFile(optionsPath)
		if err != nil {
			fmt.Println("iface:", ifaceName, "options: not found")
			continue
		}

		iface := parseOptions(ifaceName, ifacePath, string(data))

		fmt.Println("iface:", iface.Name)
		fmt.Println("  path:", iface.Path)
		fmt.Println("  type:", iface.Options["TYPE"])
		fmt.Println("  bootproto:", iface.Options["BOOTPROTO"])
		fmt.Println("  ipv4:", iface.IPv4Address)
	}
}

func parseOptions(ifaceName string, ifacePath string, data string) Iface {
	// Функция разбивки  ключ значения, тоесть представим BOOTPROTO=static
	// key = BOOTPROTO value = static
	// Что бы по умному понимать значение и не создавать их а изменять
	options := make(map[string]string)

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		options[key] = value
	}

	ipv4Address := []string{}
	ipv4Path := filepath.Join(ifacePath, "ipv4address")
	ipv4Data, err := os.ReadFile(ipv4Path)
	if err == nil {
		ipv4Lines := strings.Split(string(ipv4Data), "\n")
		for _, line := range ipv4Lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			ipv4Address = append(ipv4Address, line)
		}

	}

	return Iface{
		Name:        ifaceName,
		Path:        ifacePath,
		Options:     options,
		IPv4Address: ipv4Address,
	}
}
