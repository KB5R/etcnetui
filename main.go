package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Iface struct {
	Name        string
	Path        string
	Options     map[string]string
	IPv4Address []string
	IPv4Route   []string
}

type AppState struct {
	Ifaces        []Iface
	SelectedIndex int
}

func main() {
	etcnetPath := "testdata/ifaces"

	selectedIndex := 0

	if len(os.Args) > 1 {
		value, err := strconv.Atoi(os.Args[1])
		if err == nil {
			selectedIndex = value
		}
	}

	ifaces, err := loadIfaces(etcnetPath)
	if err != nil {
		fmt.Println("error loading ifaces:", err)
		return
	}

	state := AppState{
		Ifaces:        ifaces,
		SelectedIndex: selectedIndex,
	}

	printState(state)
	printSelectedIface(state)
}

func loadIfaces(etcnetPath string) ([]Iface, error) {
	var ifaces []Iface

	entries, err := os.ReadDir(etcnetPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		ifaceName := entry.Name()
		ifacePath := filepath.Join(etcnetPath, ifaceName)

		iface, err := parseOptions(ifaceName, ifacePath)
		if err != nil {
			continue
		}

		ifaces = append(ifaces, iface)
	}
	return ifaces, nil
}

func printState(state AppState) {
	for index, iface := range state.Ifaces {
		marker := " "
		if index == state.SelectedIndex {
			marker = ">"
		}
		ipv4 := "none"
		if len(iface.IPv4Address) > 0 {
			ipv4 = iface.IPv4Address[0]
		}
		route := "none"
		if len(iface.IPv4Route) > 0 {
			route = iface.IPv4Route[0]
		}
		fmt.Println(marker, index, "iface:", iface.Name, "type:", iface.Options["TYPE"], iface.Options["BOOTPROTO"], ipv4, route)
	}
}

func printSelectedIface(state AppState) {
	if len(state.Ifaces) == 0 {
		fmt.Println("no ifaces found")
		return
	}
	if state.SelectedIndex < 0 || state.SelectedIndex >= len(state.Ifaces) {
		fmt.Println("selected iface index is out of range")
		return
	}
	iface := state.Ifaces[state.SelectedIndex]
	fmt.Println()
	fmt.Println("selected iface:", iface.Name)
	fmt.Println("path:", iface.Path)
	fmt.Println("type:", iface.Options["TYPE"])
	fmt.Println("bootproto:", iface.Options["BOOTPROTO"])
	fmt.Println("ipv4:", iface.IPv4Address)
	fmt.Println("routes:", iface.IPv4Route)
}

func parseOptions(ifaceName string, ifacePath string) (Iface, error) {
	// Функция разбивки  ключ значения, тоесть представим BOOTPROTO=static
	// key = BOOTPROTO value = static
	// Что бы по умному понимать значение и не создавать их а изменять
	// Так же добавлен парс ip и route
	optionsPath := filepath.Join(ifacePath, "options")

	data, err := os.ReadFile(optionsPath)
	if err != nil {
		return Iface{}, err
	}

	options := make(map[string]string)

	lines := strings.Split(string(data), "\n")
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
	ipv4PathAddress := filepath.Join(ifacePath, "ipv4address")
	ipv4DataAddress, err := os.ReadFile(ipv4PathAddress)
	if err == nil {
		ipv4LinesAddress := strings.Split(string(ipv4DataAddress), "\n")
		for _, line := range ipv4LinesAddress {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			ipv4Address = append(ipv4Address, line)
		}
	}

	ipv4Route := []string{}
	ipv4PathRoute := filepath.Join(ifacePath, "ipv4route")
	ipv4DataRoute, err := os.ReadFile(ipv4PathRoute)
	if err == nil {
		ipv4LinesRoute := strings.Split(string(ipv4DataRoute), "\n")
		for _, line := range ipv4LinesRoute {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			ipv4Route = append(ipv4Route, line)
		}
	}

	return Iface{
		Name:        ifaceName,
		Path:        ifacePath,
		Options:     options,
		IPv4Address: ipv4Address,
		IPv4Route:   ipv4Route,
	}, nil
}
