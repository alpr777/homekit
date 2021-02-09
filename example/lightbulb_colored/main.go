package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alpr777/homekit"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

const (
	accessoryName string = "light"
	accessorySn   string = "ExmplLB"
	accessoryPin  string = "19283746"
)

func main() {
	// runtime.GOMAXPROCS(4)
	// log.Debug.Enable()
	acc := homekit.NewAccessoryLightbulbColored(accessory.Info{Name: accessoryName, SerialNumber: accessorySn, Manufacturer: "EXAMPLE", Model: "ACC-TEST", FirmwareRevision: "1.2"})
	transp, err := hc.NewIPTransport(hc.Config{StoragePath: "./" + acc.Info.SerialNumber.GetValue(), Pin: accessoryPin}, acc.Accessory)
	if err != nil {
		fmt.Println("accessory [", acc.Info.SerialNumber.GetValue(), "/", acc.Info.Name.GetValue(), "]", "error create transport:", err)
		os.Exit(1)
	}
	go func() {
		tickerUpdateState := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-tickerUpdateState.C:
				acc.LightbulbColored.On.SetValue(!acc.LightbulbColored.On.GetValue())
				continue
			}
		}
	}()
	go acc.LightbulbColored.On.OnValueRemoteUpdate(func(state bool) { fmt.Printf("acc remote update on: %T - %v \n", state, state) })
	go acc.LightbulbColored.Brightness.OnValueRemoteUpdate(func(state int) { fmt.Printf("acc remote update brightness: %T - %v \n", state, state) })
	go acc.LightbulbColored.Saturation.OnValueRemoteUpdate(func(state float64) { fmt.Printf("acc remote update saturation: %T - %v \n", state, state) })
	go acc.LightbulbColored.Hue.OnValueRemoteUpdate(func(state float64) { fmt.Printf("acc remote update hue: %T - %v \n", state, state) })
	fmt.Println("homekit accessory transport start [", acc.Info.SerialNumber.GetValue(), "/", acc.Info.Name.GetValue(), "]")
	hc.OnTermination(func() { <-transp.Stop() })
	transp.Start()
}
