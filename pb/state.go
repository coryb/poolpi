package pb

type IndicatorName string

const (
	Heater1         IndicatorName = "Heater1"
	Valve3          IndicatorName = "Valve3"
	CheckSystem     IndicatorName = "CheckSystem"
	Pool            IndicatorName = "Pool"
	Spa             IndicatorName = "Spa"
	Filter          IndicatorName = "Filter"
	Lights          IndicatorName = "Lights"
	Aux1            IndicatorName = "Aux1"
	Aux2            IndicatorName = "Aux2"
	Service         IndicatorName = "Service"
	Aux3            IndicatorName = "Aux3"
	Aux4            IndicatorName = "Aux4"
	Aux5            IndicatorName = "Aux5"
	Aux6            IndicatorName = "Aux6"
	Valve4          IndicatorName = "Valve4"
	Spillover       IndicatorName = "Spillover"
	SystemOff       IndicatorName = "SystemOff"
	Aux7            IndicatorName = "Aux7"
	Aux8            IndicatorName = "Aux8"
	Aux9            IndicatorName = "Aux9"
	Aux10           IndicatorName = "Aux10"
	Aux11           IndicatorName = "Aux11"
	Aux12           IndicatorName = "Aux12"
	Aux13           IndicatorName = "Aux13"
	Aux14           IndicatorName = "Aux14"
	SuperChlorinate IndicatorName = "SuperChlorinate"
)

type Indicators map[IndicatorName]*Indicator

func (s *StateEvent) Indicators() Indicators {
	return Indicators{
		Heater1:         s.GetHeater1(),
		Valve3:          s.GetValve3(),
		CheckSystem:     s.GetCheckSystem(),
		Pool:            s.GetPool(),
		Spa:             s.GetSpa(),
		Filter:          s.GetFilter(),
		Lights:          s.GetLights(),
		Aux1:            s.GetAux1(),
		Aux2:            s.GetAux2(),
		Service:         s.GetService(),
		Aux3:            s.GetAux3(),
		Aux4:            s.GetAux4(),
		Aux5:            s.GetAux5(),
		Aux6:            s.GetAux6(),
		Valve4:          s.GetValve4(),
		Spillover:       s.GetSpillover(),
		SystemOff:       s.GetSystemOff(),
		Aux7:            s.GetAux7(),
		Aux8:            s.GetAux8(),
		Aux9:            s.GetAux9(),
		Aux10:           s.GetAux10(),
		Aux11:           s.GetAux11(),
		Aux12:           s.GetAux12(),
		Aux13:           s.GetAux13(),
		Aux14:           s.GetAux14(),
		SuperChlorinate: s.GetSuperChlorinate(),
	}
}
