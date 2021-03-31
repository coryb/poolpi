package pb

import (
	"fmt"

	"github.com/logrusorgru/aurora/v3"
)

func (e *Event) Summary() string {
	if e == nil {
		return ""
	}
	switch ev := e.Event.(type) {
	case *Event_State:
		return ev.State.Summary()
	case *Event_Message:
		return ev.Message.Summary()
	case *Event_PumpRequest:
		return ev.PumpRequest.Summary()
	case *Event_PumpStatus:
		return ev.PumpStatus.Summary()
	case *Event_MessageUpdate:
		return "Update " + ev.MessageUpdate.Message.Summary()
	case *Event_StateUpdate:
		return "Update " + ev.StateUpdate.State.Summary()
	case *Event_CurrentState:
		return ev.CurrentState.State.Summary() + " " + ev.CurrentState.Message.Summary()
	case *Event_Unknown:
		return ev.Unknown.Summary()
	}
	return ""
}

func (ev *StateEvent) Summary() string {
	if ev == nil {
		return ""
	}
	active := []string{}
	blink := func(s string, should bool) string {
		if !should {
			return s
		}
		return aurora.SlowBlink(s).String()
	}
	if ind := ev.GetHeater1(); ind.GetActive() {
		active = append(active, blink("Heater1", ind.GetCaution()))
	}
	if ind := ev.GetValve3(); ind.GetActive() {
		active = append(active, blink("Valve3", ind.GetCaution()))
	}
	if ind := ev.GetCheckSystem(); ind.GetActive() {
		active = append(active, blink("CheckSystem", ind.GetCaution()))
	}
	if ind := ev.GetPool(); ind.GetActive() {
		active = append(active, blink("Pool", ind.GetCaution()))
	}
	if ind := ev.GetSpa(); ind.GetActive() {
		active = append(active, blink("Spa", ind.GetCaution()))
	}
	if ind := ev.GetFilter(); ind.GetActive() {
		active = append(active, blink("Filter", ind.GetCaution()))
	}
	if ind := ev.GetLights(); ind.GetActive() {
		active = append(active, blink("Lights", ind.GetCaution()))
	}
	if ind := ev.GetAux1(); ind.GetActive() {
		active = append(active, blink("Aux1", ind.GetCaution()))
	}
	if ind := ev.GetAux2(); ind.GetActive() {
		active = append(active, blink("Aux2", ind.GetCaution()))
	}
	if ind := ev.GetService(); ind.GetActive() {
		active = append(active, blink("Service", ind.GetCaution()))
	}
	if ind := ev.GetAux3(); ind.GetActive() {
		active = append(active, blink("Aux3", ind.GetCaution()))
	}
	if ind := ev.GetAux4(); ind.GetActive() {
		active = append(active, blink("Aux4", ind.GetCaution()))
	}
	if ind := ev.GetAux5(); ind.GetActive() {
		active = append(active, blink("Aux5", ind.GetCaution()))
	}
	if ind := ev.GetAux6(); ind.GetActive() {
		active = append(active, blink("Aux6", ind.GetCaution()))
	}
	if ind := ev.GetValve4(); ind.GetActive() {
		active = append(active, blink("Valve4", ind.GetCaution()))
	}
	if ind := ev.GetSpillover(); ind.GetActive() {
		active = append(active, blink("Spillover", ind.GetCaution()))
	}
	if ind := ev.GetSystemOff(); ind.GetActive() {
		active = append(active, blink("SystemOff", ind.GetCaution()))
	}
	if ind := ev.GetAux7(); ind.GetActive() {
		active = append(active, blink("Aux7", ind.GetCaution()))
	}
	if ind := ev.GetAux8(); ind.GetActive() {
		active = append(active, blink("Aux8", ind.GetCaution()))
	}
	if ind := ev.GetAux9(); ind.GetActive() {
		active = append(active, blink("Aux9", ind.GetCaution()))
	}
	if ind := ev.GetAux10(); ind.GetActive() {
		active = append(active, blink("Aux10", ind.GetCaution()))
	}
	if ind := ev.GetAux11(); ind.GetActive() {
		active = append(active, blink("Aux11", ind.GetCaution()))
	}
	if ind := ev.GetAux12(); ind.GetActive() {
		active = append(active, blink("Aux12", ind.GetCaution()))
	}
	if ind := ev.GetAux13(); ind.GetActive() {
		active = append(active, blink("Aux13", ind.GetCaution()))
	}
	if ind := ev.GetAux14(); ind.GetActive() {
		active = append(active, blink("Aux14", ind.GetCaution()))
	}
	if ind := ev.GetSuperChlorinate(); ind.GetActive() {
		active = append(active, blink("SuperChlorinate", ind.GetCaution()))
	}
	return fmt.Sprintf("Active: %v", active)
}

func (ev *MessageEvent) Summary() string {
	if ev == nil {
		return ""
	}
	return fmt.Sprintf("Message: %s", ev.Fancy())
}

func (ev *PumpRequestEvent) Summary() string {
	if ev == nil {
		return ""
	}
	return fmt.Sprintf("Pump speed request: %d%%", ev.SpeedPercent)
}

func (ev *PumpStatusEvent) Summary() string {
	if ev == nil {
		return ""
	}
	return fmt.Sprintf("Pump speed status: %d%% %dW",
		ev.SpeedPercent,
		ev.PowerWatts,
	)
}

func (ev *UnknownEvent) Summary() string {
	if ev == nil {
		return ""
	}
	return fmt.Sprintf("[% x] [% x]", ev.Type, ev.Data)
}
