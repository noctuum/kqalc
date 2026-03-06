package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/godbus/dbus/v5"
)

const (
	serviceName = "org.kde.krunner1.kqalc"
	objectPath  = "/kqalc"
	ifaceName   = "org.kde.krunner1"
)

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("dbus: %v", err)
	}
	defer conn.Close()

	runner := &Runner{}
	if err := conn.Export(runner, dbus.ObjectPath(objectPath), ifaceName); err != nil {
		log.Fatalf("export: %v", err)
	}

	reply, err := conn.RequestName(serviceName, dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Fatalf("request name: %v", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatal("dbus name already taken")
	}

	log.Printf("kqalc running on %s %s", serviceName, objectPath)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("kqalc stopping")
}
