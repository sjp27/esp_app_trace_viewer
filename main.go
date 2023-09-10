// Copyright 2023 sjp27 <https://github.com/sjp27>. All rights reserved.
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

// Utility to view ESP32 application tracing.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

const version = "v1.0"

func main() {
	if len(os.Args) < 4 {
		fmt.Println(version + " Usage: esp_app_trace_viewer <interface> <target> <filename>")
		fmt.Println("e.g.")
		fmt.Println("esp_app_trace_viewer interface/ftdi/esp32_devkitj_v1.cfg target/esp32s3.cfg trace.txt")
		os.Exit(0)
	}

	err := os.Remove(os.Args[3])

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

	_, err = conn.Write([]byte("esp apptrace start file://" + os.Args[3] + " 1 -1 -1 0 0\n"))
	if err != nil {
		println("OpenOCD apptrace start failed:", err.Error())
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	go monitorFile("./" + os.Args[3])

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

func monitorFile(filename string) {
	bytesRead := 0

	for {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		fi, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		fileSize := fi.Size()

		if int(fileSize) < bytesRead {
			bytesRead = 0
		}

		if int(fileSize) > bytesRead {
			br := bufio.NewReader(file)
			bytesCount := 0
			for {

				c, err := br.ReadByte()

				if err != nil && !errors.Is(err, io.EOF) {
					fmt.Println(err)
					break
				}

				// Check end of file
				if err != nil {
					break
				}

				if bytesCount > bytesRead {
					fmt.Print(string(c))
				}

				bytesCount++
			}
			bytesRead += (bytesCount - bytesRead)
		}

		file.Close()

		time.Sleep(500 * time.Millisecond)
	}
}
