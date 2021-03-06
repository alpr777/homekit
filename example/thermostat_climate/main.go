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
	accessoryName string = "thermostat"
	accessorySn   string = "ExmplTRMCLM"
	accessoryPin  string = "11112222"
)

func main() {
	// runtime.GOMAXPROCS(4)
	// log.Debug.Enable()
	acc := homekit.NewAccessoryThermostat(accessory.Info{Name: accessoryName, SerialNumber: accessorySn, Manufacturer: "alpr777", Model: "ACC-TEST", FirmwareRevision: "1.2"})
	transp, err := hc.NewIPTransport(hc.Config{StoragePath: "./" + acc.Info.SerialNumber.GetValue(), Pin: accessoryPin}, acc.Accessory)
	if err != nil {
		fmt.Println("accessory [", acc.Info.SerialNumber.GetValue(), "/", acc.Info.Name.GetValue(), "]", "error create transport:", err)
		os.Exit(1)
	}
	go func() {
		tickerUpdateState := time.NewTicker(10 * time.Second)
		tickerUpdateTemp := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-tickerUpdateState.C:
				state := acc.Thermostat.CurrentHeatingCoolingState.GetValue()
				if state >= 2 {
					acc.Thermostat.CurrentHeatingCoolingState.SetValue(0)
				} else {
					acc.Thermostat.CurrentHeatingCoolingState.SetValue(state + 1)
				}
				fmt.Printf("acc thermostat update current state: %T - %v \n", acc.Thermostat.CurrentHeatingCoolingState.GetValue(), acc.Thermostat.CurrentHeatingCoolingState.GetValue())
				continue
			case <-tickerUpdateTemp.C:
				acc.Thermostat.CurrentTemperature.SetValue(float64(time.Now().Second()-30) + float64(time.Now().Second()+40)/100)
				fmt.Printf("acc thermostat update current temp: %T - %v \n", acc.Thermostat.CurrentTemperature.GetValue(), acc.Thermostat.CurrentTemperature.GetValue())
				continue
			}
		}
	}()
	go acc.Thermostat.TargetHeatingCoolingState.OnValueRemoteUpdate(func(v int) {
		fmt.Printf("acc thermostat remote update target state: %T - %v \n", v, v)
	})
	go acc.Thermostat.TargetTemperature.OnValueRemoteUpdate(func(v float64) {
		fmt.Printf("acc thermostat remote update target temp: %T - %v \n", v, v)
	})
	fmt.Println("homekit accessory transport start [", acc.Info.SerialNumber.GetValue(), "/", acc.Info.Name.GetValue(), "]")
	hc.OnTermination(func() { <-transp.Stop() })
	transp.Start()
}
