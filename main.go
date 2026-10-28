package main

import (
	"log"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/godbus/dbus/v5"
)

const (
	serviceName = "org.kde.krunner1.kqalc"
	objectPath  = "/kqalc"
	ifaceName   = "org.kde.krunner1"
)

// activatable reports whether the bus can start name on demand. A bus that
// started before ~/.local/share/dbus-1/services/ existed does not watch it,
// so a kqalc started by install.sh or by hand has to stay up.
func activatable(conn *dbus.Conn, name string) bool {
	var names []string
	err := conn.BusObject().Call("org.freedesktop.DBus.ListActivatableNames", 0).Store(&names)
	return err == nil && slices.Contains(names, name)
}

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("dbus: %v", err)
	}
	defer conn.Close()

	runner := newRunner(func() bool { return activatable(conn, serviceName) })
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
	select {
	case <-sig:
	case <-runner.done:
		// The next query makes the bus start kqalc again, from whatever
		// binary is installed by then.
		if _, err := conn.ReleaseName(serviceName); err != nil {
			log.Printf("release name: %v", err)
		}
	}

	log.Println("kqalc stopping")
}
