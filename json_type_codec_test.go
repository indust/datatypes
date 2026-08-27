package datatypes_test

import (
	"reflect"
	"testing"

	"github.com/ugorji/go/codec"
	"gorm.io/datatypes"
)

var _ codec.Selfer = (*datatypes.JSONType[any])(nil)

func TestJSONTypeCodec(t *testing.T) {
	type metadata struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	type payload struct {
		ID       int                          `json:"id"`
		Metadata datatypes.JSONType[metadata] `json:"metadata"`
		Detail   datatypes.JSONType[any]      `json:"detail"`
	}

	want := payload{
		ID:       1,
		Metadata: datatypes.NewJSONType(metadata{Name: "test", Count: 2}),
		Detail:   datatypes.NewJSONType[any](map[string]any{"enabled": true}),
	}
	var handle codec.MsgpackHandle
	handle.MapType = reflect.TypeOf(map[string]any{})

	var data []byte
	if err := codec.NewEncoderBytes(&data, &handle).Encode(want); err != nil {
		t.Fatal(err)
	}

	var wire struct {
		Metadata metadata       `json:"metadata"`
		Detail   map[string]any `json:"detail"`
	}
	if err := codec.NewDecoderBytes(data, &handle).Decode(&wire); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(wire.Metadata, want.Metadata.Data()) || wire.Detail["enabled"] != true {
		t.Fatalf("encoded as %#v", wire)
	}

	var got payload
	if err := codec.NewDecoderBytes(data, &handle).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip got %#v, want %#v", got, want)
	}
}
