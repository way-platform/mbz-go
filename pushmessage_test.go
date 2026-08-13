package mbz

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/way-platform/mbz-go/api/servicesv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPushMessage_AsProto_mixedSnapshot(t *testing.T) {
	const ts = int64(1_700_000_000_000)
	msg := PushMessage{
		MessageID:        "msg-1",
		VIN:              "WDD123",
		Timestamp:        ts,
		MessageType:      "vehiclesignal",
		Version:          "1.0",
		ServiceID:        "svc-1",
		SendingBehaviour: servicesv1.SignalSendingBehaviourONCHANGE,
		Data: json.RawMessage(`{
			"signals": [
				{
					"name": "Service.Content",
					"timestamp": 1700000000000,
					"value": "SERVICE_B",
					"type": "ENUM"
				},
				{
					"name": "InternalCombustionEngine.LimpMode.IsActive",
					"timestamp": 1700000000000,
					"value": "false",
					"type": "BOOLEAN"
				},
				{
					"name": "SupplyBattery.Status",
					"timestamp": 1700000000000,
					"value": "GREEN",
					"type": "ENUM"
				}
			]
		}`),
	}

	got, err := msg.AsProto()
	if err != nil {
		t.Fatalf("AsProto() error = %v, want nil", err)
	}

	want := &mbzv1.PushMessage{}
	want.SetMessageId("msg-1")
	want.SetVin("WDD123")
	want.SetTime(timestamppb.New(time.UnixMilli(ts).UTC()))
	want.SetMessageType(mbzv1.MessageType_SIGNALS)
	want.SetVersion("1.0")
	want.SetServiceId("svc-1")
	want.SetSendingBehavior(mbzv1.SendingBehavior_ON_CHANGE)

	wantSignals := []*mbzv1.Signal{
		signalProto(
			mbzv1.SignalIdentifier_INTERNALCOMBUSTIONENGINE_LIMPMODE_ISACTIVE,
			mbzv1.SignalType_BOOLEAN,
			ts,
			func(s *mbzv1.Signal) { s.SetBoolValue(false) },
		),
		signalProto(
			mbzv1.SignalIdentifier_SUPPLYBATTERY_STATUS,
			mbzv1.SignalType_ENUM,
			ts,
			func(s *mbzv1.Signal) { s.SetEnumValue("GREEN") },
		),
	}
	want.SetSignals(wantSignals)

	wantUnmapped := []*mbzv1.UnmappedSignal{unmappedProto(
		"Service.Content",
		ts,
		"SERVICE_B",
		"ENUM",
		mbzv1.UnmappedSignal_UNKNOWN_NAME,
	)}
	want.SetUnmappedSignals(wantUnmapped)

	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Fatalf("AsProto() mismatch (-want +got):\n%s", diff)
	}
}

func TestPushMessage_AsProto_invalidEnum(t *testing.T) {
	const ts = int64(1_700_000_000_000)
	msg := PushMessage{
		MessageID:        "msg-2",
		VIN:              "WDD123",
		Timestamp:        ts,
		MessageType:      "vehiclesignal",
		Version:          "1.0",
		ServiceID:        "svc-1",
		SendingBehaviour: servicesv1.SignalSendingBehaviourONCHANGE,
		Data: json.RawMessage(`{
			"signals": [
				{
					"name": "SupplyBattery.Status",
					"timestamp": 1700000000000,
					"value": "PURPLE",
					"type": "ENUM"
				},
				{
					"name": "InternalCombustionEngine.LimpMode.IsActive",
					"timestamp": 1700000000000,
					"value": "false",
					"type": "BOOLEAN"
				}
			]
		}`),
	}

	got, err := msg.AsProto()
	if err != nil {
		t.Fatalf("AsProto() error = %v, want nil", err)
	}

	want := &mbzv1.PushMessage{}
	want.SetMessageId("msg-2")
	want.SetVin("WDD123")
	want.SetTime(timestamppb.New(time.UnixMilli(ts).UTC()))
	want.SetMessageType(mbzv1.MessageType_SIGNALS)
	want.SetVersion("1.0")
	want.SetServiceId("svc-1")
	want.SetSendingBehavior(mbzv1.SendingBehavior_ON_CHANGE)
	want.SetSignals([]*mbzv1.Signal{
		signalProto(
			mbzv1.SignalIdentifier_INTERNALCOMBUSTIONENGINE_LIMPMODE_ISACTIVE,
			mbzv1.SignalType_BOOLEAN,
			ts,
			func(s *mbzv1.Signal) { s.SetBoolValue(false) },
		),
	})
	want.SetUnmappedSignals([]*mbzv1.UnmappedSignal{unmappedProto(
		"SupplyBattery.Status",
		ts,
		"PURPLE",
		"ENUM",
		mbzv1.UnmappedSignal_INVALID_ENUM,
	)})

	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Fatalf("AsProto() mismatch (-want +got):\n%s", diff)
	}
}

func TestPushMessage_AsProto_allKnown(t *testing.T) {
	const ts = int64(1_700_000_000_000)
	msg := PushMessage{
		MessageID:        "msg-3",
		VIN:              "WDD123",
		Timestamp:        ts,
		MessageType:      "vehiclesignal",
		Version:          "1.0",
		ServiceID:        "svc-1",
		SendingBehaviour: servicesv1.SignalSendingBehaviourONCHANGE,
		Data: json.RawMessage(`{
			"signals": [
				{
					"name": "InternalCombustionEngine.LimpMode.IsActive",
					"timestamp": 1700000000000,
					"value": "false",
					"type": "BOOLEAN"
				},
				{
					"name": "SupplyBattery.Status",
					"timestamp": 1700000000000,
					"value": "GREEN",
					"type": "ENUM"
				}
			]
		}`),
	}

	got, err := msg.AsProto()
	if err != nil {
		t.Fatalf("AsProto() error = %v, want nil", err)
	}
	if len(got.GetUnmappedSignals()) != 0 {
		t.Fatalf("GetUnmappedSignals() = %v, want empty", got.GetUnmappedSignals())
	}
	if len(got.GetSignals()) != 2 {
		t.Fatalf("GetSignals() len = %d, want 2", len(got.GetSignals()))
	}
}

func TestPushMessage_AsProto_unknownSendingBehavior(t *testing.T) {
	const ts = int64(1_700_000_000_000)
	msg := PushMessage{
		MessageID:        "msg-4",
		VIN:              "WDD123",
		Timestamp:        ts,
		MessageType:      "vehiclesignal",
		Version:          "1.0",
		ServiceID:        "svc-1",
		SendingBehaviour: servicesv1.SignalSendingBehaviour("ON_FUTURE_BEHAVIOR"),
		Data: json.RawMessage(`{
			"signals": [
				{
					"name": "InternalCombustionEngine.LimpMode.IsActive",
					"timestamp": 1700000000000,
					"value": "false",
					"type": "BOOLEAN"
				}
			]
		}`),
	}

	got, err := msg.AsProto()
	if err != nil {
		t.Fatalf("AsProto() error = %v, want nil", err)
	}
	if got.GetSendingBehavior() != mbzv1.SendingBehavior_SENDING_BEHAVIOR_UNSPECIFIED {
		t.Fatalf(
			"GetSendingBehavior() = %v, want %v",
			got.GetSendingBehavior(),
			mbzv1.SendingBehavior_SENDING_BEHAVIOR_UNSPECIFIED,
		)
	}
	if len(got.GetSignals()) != 1 {
		t.Fatalf("GetSignals() len = %d, want 1", len(got.GetSignals()))
	}
}

func TestPushMessage_AsProto_garbageDataJSON(t *testing.T) {
	msg := PushMessage{
		MessageID:        "msg-5",
		VIN:              "WDD123",
		Timestamp:        1_700_000_000_000,
		MessageType:      "vehiclesignal",
		Version:          "1.0",
		ServiceID:        "svc-1",
		SendingBehaviour: servicesv1.SignalSendingBehaviourONCHANGE,
		Data:             json.RawMessage(`{"signals": [`),
	}

	_, err := msg.AsProto()
	if err == nil {
		t.Fatal("AsProto() error = nil, want unmarshal error")
	}
}

func signalProto(
	id mbzv1.SignalIdentifier,
	signalType mbzv1.SignalType,
	ts int64,
	setValue func(*mbzv1.Signal),
) *mbzv1.Signal {
	var s mbzv1.Signal
	s.SetId(id)
	s.SetType(signalType)
	s.SetTime(timestamppb.New(time.UnixMilli(ts).UTC()))
	setValue(&s)
	return &s
}

func unmappedProto(
	name string,
	ts int64,
	value string,
	signalType string,
	reason mbzv1.UnmappedSignal_Reason,
) *mbzv1.UnmappedSignal {
	var u mbzv1.UnmappedSignal
	u.SetName(name)
	u.SetTime(timestamppb.New(time.UnixMilli(ts).UTC()))
	u.SetValue(value)
	u.SetType(signalType)
	u.SetReason(reason)
	return &u
}
