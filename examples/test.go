package main

import (
	"../itch"
	"bufio"
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "%s usage:\n\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func handleSysEvents(msg itch.ITCHMessage) {
	ev := msg.(*itch.SystemEventMessage)
	fmt.Printf("System Event: %c %d\n", ev.EventCode, ev.Timestamp)
}

func handleAddOrder(msg itch.ITCHMessage) {
	ev := msg.(*itch.AddOrderMessage)
	/*    order := binary.ReadUvarint(buf[5:13])
	shares := binary.ReadVarint(buf[14:18])
	price := binary.ReadVarint(buf[26:30])*/
	fmt.Printf("Add Order[%d]: %c %d %s %d\n", ev.Timestamp, ev.Indicator, ev.Shares, ev.Ticker, ev.Price)
}

func main() {
	flag.Usage = usage
	var fname = flag.String("file", "", "ITCH data file")
	flag.Parse()

	if *fname == "" {
		usage()
	}

	inp, err := os.Open(*fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n\n", err)
		os.Exit(1)
	}

	var proc = itch.New(bufio.NewReader(inp))
	proc.AddHandler(handleSysEvents, 'S')
	proc.AddHandler(handleAddOrder, 'A')
	proc.Process()

	//msg := &itch.SystemEventMessage{Type: 'S', Description: "The Descr", Timestamp: "The Time"}
}
