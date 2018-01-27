package driver

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

var (
	// Drivers Array of available drivers
	Drivers map[string]Driver
)

//CPUFlag CPU flag type
type CPUFlag int

const (
	//CPUOnline CPU is online for domain
	CPUOnline CPUFlag = iota
	//CPURunning Domain is currently running on CPU
	CPURunning
	//CPUHalted CPU is halted/blocked for domain (waiting on IO etc)
	CPUHalted
	//CPUPaused CPU is paused (no further CPU scheduling)
	CPUPaused
)

//DomainFlag Domain Flag
type DomainFlag int

const (
	//DomainOnline Domain is online and running
	DomainOnline DomainFlag = iota
	//DomainShutdown Domain is offline/shutdown
	DomainShutdown
	//DomainCrashed Domain has crashed
	DomainCrashed
	//DomainDying Domain is dying (restart, shutdown)
	DomainDying
	//DomainPaused Domain is waiting for CPU time
	DomainPaused
)

// Driver Driver struct
type Driver interface {
	Name() string
	Detect() bool
	Collect(bool, bool, bool) (map[DomainID]*Domain, error)
	Close()
}

// DomainID Domain #ID
type DomainID uint64

// Timestamp Collection timestamp
type Timestamp int64

// Domain Domain
type Domain struct {
	Name   string
	ID     DomainID
	UUID   string
	OSType string
	Time   Timestamp
	Flags  DomainFlag

	Cpus       []CPU
	Blocks     []BlockDevice
	Interfaces []NetworkInterface

	prv interface{}
}

// BlockIO Block IO
type BlockIO struct {
	Operations uint64
	Bytes      uint64
	Sectors    uint64
	Absolute   bool
}

// BlockDevice Block Device
type BlockDevice struct {
	Name     string
	ReadOnly bool
	IsDisk   bool
	IsCDrom  bool
	Read     BlockIO
	Write    BlockIO
	Flush    BlockIO
}

// CPU CPU
type CPU struct {
	ID      uint64
	Flags   CPUFlag
	Time    float64
	Idle    float64
	IdleSet bool
	Load1   float64
	Load5   float64
	Load15  float64
}

// NetworkIO Network IO
type NetworkIO struct {
	Bytes   uint64
	Packets uint64
	Errors  uint64
	Drops   uint64
}

// NetworkInterface Network Interface
type NetworkInterface struct {
	Name    string
	Mac     net.HardwareAddr
	Bridges []string
	RX      NetworkIO
	TX      NetworkIO
}

func init() {
	Drivers = make(map[string]Driver)
}

//StringToDomainID Convert string to DomainID
func StringToDomainID(id string) DomainID {
	domid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return DomainID(0)
	}
	return DomainID(domid)
}

//AvailableDrivers List of registered driver names
func AvailableDrivers() (drivers []string) {
	for name := range Drivers {
		drivers = append(drivers, name)
	}
	return
}

//IsDriver Test if supplied interface implements the Driver interface
func IsDriver(drv interface{}) bool {
	_, ok := drv.(Driver)
	return ok
}

// RegisterDriver Register Driver
func RegisterDriver(drv interface{}) (err error) {

	//var driver Driver
	//var ok bool

	//log.Printf("Attempting to register driver: %T", drv)

	driver, ok := drv.(Driver)
	if !ok {
		err = fmt.Errorf("RegisterDriver: Invalid type %T", drv)
		log.Fatal(err)
		return
	}

	if driver.Name() == "" {
		err = fmt.Errorf("RegisterDriver: Empty driver name")
		log.Fatal(err)
		return
	}

	Drivers[strings.ToUpper(driver.Name())] = driver

	if !driver.Detect() {
		//log.Printf("RegisterDriver: %s Hypervisor not detected", driver.Name())
		return
	}

	//log.Printf("Driver %s registered successfully.", driver.Name())

	return
}
