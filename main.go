// Copyright 2023 sjp27 <https://github.com/sjp27>. All rights reserved.
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

// Utility to view ESP32 application tracing.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

const version = "v2.0"

func main() {
	if len(os.Args) < 3 {
		fmt.Println(version + " Usage: esp_app_trace_viewer <interface> <target>")
		fmt.Println("e.g.")
		fmt.Println("esp_app_trace_viewer interface/ftdi/esp32_devkitj_v1.cfg target/esp32s3.cfg")
		os.Exit(0)
	}

	go tcpServer()

	cmd := exec.Command("openocd", "-c gdb_port 50540", "-c telnet_port 50538", "-f", os.Args[1], "-c adapter_khz 10000", "-f", os.Args[2])
	if err := cmd.Start(); err != nil {
		println("OpenOCD start failed:", err.Error())
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", "127.0.0.1:50538")
	if err != nil {
		println("OpenOCD connect failed:", err.Error())
		log.Fatal(err)
	}

	const s = "esp apptrace start tcp://127.0.0.1:50536 1 -1 -1 0 0\n"
	_, err = conn.Write([]byte(s))
	if err != nil {
		println("OpenOCD apptrace start failed:", err.Error())
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		char, _, _ := reader.ReadRune()

		if char == 'x' {
			println("Exit")
			break
		} else if char == 'r' {
			println("Reset target")
			_, err = conn.Write([]byte("reset run\n"))
			if err != nil {
				println("OpenOCD target reset failed:", err.Error())
				log.Fatal(err)
			}
		}
	}

	conn.Close()

	if err := cmd.Process.Kill(); err != nil {
		log.Fatal("Failed to stop OpenOCD:", err)
	}
}

func tcpServer() {
	ln, _ := net.Listen("tcp", "127.0.0.1:50536")

	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(string(message))
	}
}
