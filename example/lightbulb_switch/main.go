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
	accessorySn   string = "ExmplLBS"
	accessoryPin  string = "11112222"
)

func main() {
	// runtime.GOMAXPROCS(4)
	// log.Debug.Enable()
	acc := homekit.NewAccessoryLightbulbSwitch(accessory.Info{Name: accessoryName, SerialNumber: accessorySn, Manufacturer: "alpr777", Model: "ACC-TEST", FirmwareRevision: "1.2"})
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
				acc.LightbulbSwitch.On.SetValue(!acc.LightbulbSwitch.On.GetValue())
				fmt.Printf("acc lightbulb update on: %T - %v \n", acc.LightbulbSwitch.On.GetValue(), acc.LightbulbSwitch.On.GetValue())
				continue
			}
		}
	}()
	go acc.LightbulbSwitch.On.OnValueRemoteUpdate(func(v bool) {
		fmt.Printf("acc lightbulb remote update on: %T - %v \n", v, v)
	})
	fmt.Println("homekit accessory transport start [", acc.Info.SerialNumber.GetValue(), "/", acc.Info.Name.GetValue(), "]")
	hc.OnTermination(func() { <-transp.Stop() })
	transp.Start()
}
