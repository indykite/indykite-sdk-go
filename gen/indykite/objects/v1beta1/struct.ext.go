// Copyright (c) 2022 IndyKite
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objectsv1beta1

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	nullValue                   = &Value{Value: &Value_NullValue{NullValue: structpb.NullValue_NULL_VALUE}}
	_         protoreflect.List = &ArrayValue{}
)

func Null() *Value {
	return &Value{Value: &Value_NullValue{NullValue: structpb.NullValue_NULL_VALUE}}
}

func Bool(v bool) *Value {
	return &Value{Value: &Value_BoolValue{BoolValue: v}}
}

func Int64(v int64) *Value {
	return &Value{Value: &Value_IntegerValue{IntegerValue: v}}
}

func Float64(v float64) *Value {
	return &Value{Value: &Value_DoubleValue{DoubleValue: v}}
}

func String(v string) *Value {
	return &Value{
		Value: &Value_StringValue{
			StringValue: v,
		},
	}
}

func Any(v proto.Message) (*Value, error) {
	if a, ok := v.(*anypb.Any); ok {
		return &Value{Value: &Value_AnyValue{AnyValue: a}}, nil
	}
	a, err := anypb.New(v)
	if err != nil {
		return nil, err
	}
	return &Value{Value: &Value_AnyValue{AnyValue: a}}, nil
}

func ToValue(v interface{}) (*Value, error) {
	val, _, err := toProtoValue(reflect.ValueOf(v))
	return val, err
}

func FromValue(v *Value) (interface{}, error) {
	return createFromProtoValue(v)
}

func ToMapValue(in map[string]interface{}) (map[string]*Value, error) {
	if in == nil {
		return nil, nil
	}
	out := make(map[string]*Value, len(in))
	for k, v := range in {
		val, _, err := toProtoValue(reflect.ValueOf(v))
		if err != nil {
			return nil, err
		}
		if val == nil {
			continue
		}
		out[k] = val
	}
	return out, nil
}

func ToMap(fields map[string]*Value) (map[string]interface{}, error) {
	if fields == nil {
		return nil, nil
	}
	ret := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		r, err := createFromProtoValue(v)
		if err != nil {
			return nil, err
		}
		ret[k] = r
	}
	return ret, nil
}

func createFromProtoValue(v *Value) (interface{}, error) {
	if v == nil || v.Value == nil {
		return nil, nil
	}
	switch v := v.Value.(type) {
	case *Value_NullValue:
		return nil, nil
	case *Value_BoolValue:
		return v.BoolValue, nil
	case *Value_BytesValue:
		return v.BytesValue, nil
	case *Value_StringValue:
		return v.StringValue, nil
	case *Value_DoubleValue:
		return v.DoubleValue, nil
	case *Value_IntegerValue:
		return v.IntegerValue, nil
	case *Value_UnsignedIntegerValue:
		return v.UnsignedIntegerValue, nil
	case *Value_DurationValue:
		return v.DurationValue.AsDuration().String(), nil
	case *Value_ValueTime:
		return v.ValueTime.AsTime().UTC().Format(time.RFC3339), nil
	case *Value_GeoPointValue:
		return fmt.Sprintf("POINT (%v %v)", v.GeoPointValue.GetLatitude(), v.GeoPointValue.GetLongitude()), nil
	case *Value_ArrayValue:
		if v.ArrayValue == nil {
			return nil, nil
		}
		values := v.ArrayValue.Values
		ret := make([]interface{}, len(values))
		for i, v := range values {
			r, err := createFromProtoValue(v)
			if err != nil {
				return nil, err
			}
			ret[i] = r
		}
		return ret, nil
	case *Value_MapValue:
		if v.MapValue == nil {
			return nil, nil
		}
		fields := v.MapValue.Fields
		ret := make(map[string]interface{}, len(fields))
		for k, v := range fields {
			r, err := createFromProtoValue(v)
			if err != nil {
				return nil, err
			}
			ret[k] = r
		}
		return ret, nil
	case *Value_AnyValue:
		if v.AnyValue == nil {
			return nil, nil
		}
		typeURL := v.AnyValue.GetTypeUrl()
		value := v.AnyValue.GetValue()
		ret := map[string]interface{}{"typeUrl": typeURL, "value": value}
		return ret, nil
	default:
		return nil, fmt.Errorf("unknown value type %T", v)
	}
}

func toProtoValue(v reflect.Value) (pbv *Value, sawTransform bool, err error) {
	if !v.IsValid() {
		return nullValue, false, nil
	}
	switch v.Kind() {
	case reflect.Bool:
		return &Value{Value: &Value_BoolValue{BoolValue: v.Bool()}}, false, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &Value{Value: &Value_IntegerValue{IntegerValue: v.Int()}}, false, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return &Value{Value: &Value_UnsignedIntegerValue{UnsignedIntegerValue: v.Uint()}}, false, nil
	case reflect.Float32, reflect.Float64:
		return &Value{Value: &Value_DoubleValue{DoubleValue: v.Float()}}, false, nil
	case reflect.String:
		return &Value{Value: &Value_StringValue{StringValue: v.String()}}, false, nil
	case reflect.Array:
		return arrayToProtoValue(v)
	case reflect.Slice:
		return sliceToProtoValue(v)
	case reflect.Map:
		return mapToProtoValue(v)
	case reflect.Ptr:
		if v.IsNil() {
			return nullValue, false, nil
		}
		return toProtoValue(v.Elem())
	case reflect.Interface:
		if v.NumMethod() == 0 { // empty interface: recurse on its contents
			return toProtoValue(v.Elem())
		}
		fallthrough // any other interface value is an error
	default:
		return nil, false, fmt.Errorf("cannot convert type %s to value", v.Type())
	}
}

// arrayToProtoValue converts a array to a Value protobuf and reports whether a transform was encountered.
func arrayToProtoValue(v reflect.Value) (*Value, bool, error) {
	values := make([]*Value, v.Len())
	for i := 0; i < v.Len(); i++ {
		val, sawTransform, err := toProtoValue(v.Index(i))
		if err != nil {
			return nil, false, err
		}
		if sawTransform {
			return nil, false, fmt.Errorf("transforms cannot occur in an array, but saw some in %v", v.Index(i))
		}
		values[i] = val
	}
	return &Value{Value: &Value_ArrayValue{ArrayValue: &ArrayValue{Values: values}}}, false, nil
}

// sliceToProtoValue converts a slice to a Value protobuf and reports whether a transform was encountered.
func sliceToProtoValue(v reflect.Value) (*Value, bool, error) {
	// A nil slice is converted to a null value.
	if v.IsNil() {
		return nullValue, false, nil
	}
	return arrayToProtoValue(v)
}

// mapToProtoValue converts a map to a Value protobuf and reports whether a transform was encountered.
func mapToProtoValue(v reflect.Value) (*Value, bool, error) {
	if v.Type().Key().Kind() != reflect.String {
		return nil, false, errors.New("map key type must be string")
	}
	// A nil map is converted to a null value.
	if v.IsNil() {
		return nullValue, false, nil
	}
	m := map[string]*Value{}
	sawTransform := false
	for _, k := range v.MapKeys() {
		mi := v.MapIndex(k)
		val, sst, err := toProtoValue(mi)
		if err != nil {
			return nil, false, err
		}
		if sst {
			sawTransform = true
		}
		if val == nil {
			continue
		}
		m[k.String()] = val
	}
	var pv *Value
	if len(m) == 0 && sawTransform {
		// The entire map consisted of transform values.
		pv = nil
	} else {
		pv = &Value{Value: &Value_MapValue{MapValue: &MapValue{Fields: m}}}
	}
	return pv, sawTransform, nil
}

func (x *ArrayValue) Len() int {
	return len(x.Values)
}

func (x *ArrayValue) Get(i int) protoreflect.Value {
	switch v := x.Values[i].Value.(type) {
	case *Value_NullValue:
		return protoreflect.ValueOf(v.NullValue)
	case *Value_BoolValue:
		return protoreflect.ValueOfBool(v.BoolValue)
	case *Value_BytesValue:
		return protoreflect.ValueOfBytes(v.BytesValue)
	case *Value_StringValue:
		return protoreflect.ValueOfString(v.StringValue)
	case *Value_DoubleValue:
		return protoreflect.ValueOfFloat64(v.DoubleValue)
	case *Value_IntegerValue:
		return protoreflect.ValueOfInt64(v.IntegerValue)
	case *Value_UnsignedIntegerValue:
		return protoreflect.ValueOfUint64(v.UnsignedIntegerValue)
	case *Value_ArrayValue:
		return protoreflect.ValueOf(v.ArrayValue)
	case *Value_MapValue:
		return protoreflect.ValueOf(v.MapValue)
	}
	return protoreflect.ValueOf(nil)
}

func (x *ArrayValue) Set(i int, value protoreflect.Value) {
	panic("implement me")
}

func (x *ArrayValue) Append(value protoreflect.Value) {
	panic("implement me")
}

func (x *ArrayValue) AppendMutable() protoreflect.Value {
	panic("implement me")
}

func (x *ArrayValue) Truncate(i int) {
	panic("implement me")
}

func (x *ArrayValue) NewElement() protoreflect.Value {
	panic("implement me")
}

func (x *ArrayValue) IsValid() bool {
	return true
}
