package web

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/brandon1024/OpenEVT/internal/types"
)

var (
	reg = prometheus.NewRegistry()

	connected = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "openevt_connected",
			Help: "Connection status to the inverter (0-disconnected, 1-connected).",
		},
		[]string{"addr", "sn"},
	)
	power = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "openevt_power_ac",
			Help: "Total instantaneous power (AC) of both inverter modules, in W.",
		},
		[]string{"addr", "sn"},
	)
	energy = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "openevt_energy",
			Help: "Total accumulated energy generated by both inverter modules, in kWh.",
		},
		[]string{"addr", "sn"},
	)

	moduleInputVoltageDC = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "openevt_module_input_voltage_dc",
			Help: "Input voltage (DC) for an inverter module, in volts.",
		},
		[]string{"addr", "sn", "module_id", "firmware_version"},
	)
	moduleOutputPowerAC = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "openevt_module_output_power_ac",
			Help: "Instantaneous output power (AC) for an inverter module, in W.",
		},
		[]string{"addr", "sn", "module_id", "firmware_version"},
	)
	moduleTotalEnergy = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "openevt_module_total_energy",
			Help: "Accumulated total energy generated generated for an inverter module, in kWh.",
		},
		[]string{"addr", "sn", "module_id", "firmware_version"},
	)
	moduleTemperature = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "openevt_module_temp",
			Help: "Module temperature, in degrees C.",
		},
		[]string{"addr", "sn", "module_id", "firmware_version"},
	)
	moduleOutputVoltageAC = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "openevt_module_output_voltage_ac",
			Help: "Output voltage (AC) for an inverter module, in volts.",
		},
		[]string{"addr", "sn", "module_id", "firmware_version"},
	)
	moduleOutputFrequencyAC = promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "openevt_module_output_frequency_ac",
			Help: "Output frequency (AC) for an inverter module, in Hz.",
		},
		[]string{"addr", "sn", "module_id", "firmware_version"},
	)
)

var (
	inverter    types.InverterStatus
	inverterMux sync.RWMutex
)

func get() types.InverterStatus {
	inverterMux.RLock()
	defer inverterMux.RUnlock()

	return inverter
}

func set(status types.InverterStatus) {
	inverterMux.Lock()
	defer inverterMux.Unlock()

	inverter = status
}

func UpdateConnectionStatus(addr, sn string, status float64) {
	labels := prometheus.Labels{
		"addr": addr,
		"sn":   sn,
	}

	connected.With(labels).Set(status)
}

func Update(addr string, status *types.InverterStatus) {
	labels := prometheus.Labels{
		"addr": addr,
		"sn":   status.InverterId,
	}

	set(*status)

	power.With(labels).Set(status.Module1.OutputPowerAC + status.Module2.OutputPowerAC)
	energy.With(labels).Set(status.Module1.TotalEnergy + status.Module2.TotalEnergy)

	UpdateModule(addr, status.InverterId, &status.Module1)
	UpdateModule(addr, status.InverterId, &status.Module2)
}

func UpdateModule(addr, sn string, module *types.InverterModuleStatus) {
	labels := prometheus.Labels{
		"addr":             addr,
		"sn":               sn,
		"module_id":        module.ModuleId,
		"firmware_version": module.FirmwareVersion,
	}

	moduleInputVoltageDC.With(labels).Set(module.InputVoltageDC)
	moduleOutputPowerAC.With(labels).Set(module.OutputPowerAC)
	moduleTotalEnergy.With(labels).Set(module.TotalEnergy)
	moduleTemperature.With(labels).Set(module.Temperature)
	moduleOutputVoltageAC.With(labels).Set(module.OutputPowerAC)
	moduleOutputFrequencyAC.With(labels).Set(module.OutputFrequencyAC)
}
