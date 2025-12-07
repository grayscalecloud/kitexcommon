package hd

import (
	"reflect"
	"testing"
)

// assertEqual 辅助函数，用于断言两个值是否相等
func assertEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("期望值: %v, 实际值: %v", expected, actual)
	}
}

// assertNil 辅助函数，用于断言值是否为 nil
func assertNil(t *testing.T, actual interface{}) {
	if actual == nil {
		return
	}
	v := reflect.ValueOf(actual)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		if !v.IsNil() {
			t.Errorf("期望值为 nil, 实际值: %v", actual)
		}
	default:
		t.Errorf("期望值为 nil, 实际值: %v (类型: %T)", actual, actual)
	}
}

func Test_Trans(t *testing.T) {
	str := String("tea")
	strVal := StringValue(str)
	assertEqual(t, "tea", strVal)
	assertEqual(t, "", StringValue(nil))

	strSlice := StringSlice([]string{"tea"})
	strSliceVal := StringSliceValue(strSlice)
	assertEqual(t, []string{"tea"}, strSliceVal)
	assertNil(t, StringSlice(nil))
	assertNil(t, StringSliceValue(nil))

	b := Bool(true)
	bVal := BoolValue(b)
	assertEqual(t, true, bVal)
	assertEqual(t, false, BoolValue(nil))

	bSlice := BoolSlice([]bool{false})
	bSliceVal := BoolSliceValue(bSlice)
	assertEqual(t, []bool{false}, bSliceVal)
	assertNil(t, BoolSlice(nil))
	assertNil(t, BoolSliceValue(nil))

	f64 := Float64(2.00)
	f64Val := Float64Value(f64)
	assertEqual(t, float64(2.00), f64Val)
	assertEqual(t, float64(0), Float64Value(nil))

	f32 := Float32(2.00)
	f32Val := Float32Value(f32)
	assertEqual(t, float32(2.00), f32Val)
	assertEqual(t, float32(0), Float32Value(nil))

	f64Slice := Float64Slice([]float64{2.00})
	f64SliceVal := Float64ValueSlice(f64Slice)
	assertEqual(t, []float64{2.00}, f64SliceVal)
	assertNil(t, Float64Slice(nil))
	assertNil(t, Float64ValueSlice(nil))

	f32Slice := Float32Slice([]float32{2.00})
	f32SliceVal := Float32ValueSlice(f32Slice)
	assertEqual(t, []float32{2.00}, f32SliceVal)
	assertNil(t, Float32Slice(nil))
	assertNil(t, Float32ValueSlice(nil))

	i := Int(1)
	iVal := IntValue(i)
	assertEqual(t, 1, iVal)
	assertEqual(t, 0, IntValue(nil))

	i8 := Int8(int8(1))
	i8Val := Int8Value(i8)
	assertEqual(t, int8(1), i8Val)
	assertEqual(t, int8(0), Int8Value(nil))

	i16 := Int16(int16(1))
	i16Val := Int16Value(i16)
	assertEqual(t, int16(1), i16Val)
	assertEqual(t, int16(0), Int16Value(nil))

	i32 := Int32(int32(1))
	i32Val := Int32Value(i32)
	assertEqual(t, int32(1), i32Val)
	assertEqual(t, int32(0), Int32Value(nil))

	i64 := Int64(int64(1))
	i64Val := Int64Value(i64)
	assertEqual(t, int64(1), i64Val)
	assertEqual(t, int64(0), Int64Value(nil))

	iSlice := IntSlice([]int{1})
	iSliceVal := IntValueSlice(iSlice)
	assertEqual(t, []int{1}, iSliceVal)
	assertNil(t, IntSlice(nil))
	assertNil(t, IntValueSlice(nil))

	i8Slice := Int8Slice([]int8{1})
	i8ValSlice := Int8ValueSlice(i8Slice)
	assertEqual(t, []int8{1}, i8ValSlice)
	assertNil(t, Int8Slice(nil))
	assertNil(t, Int8ValueSlice(nil))

	i16Slice := Int16Slice([]int16{1})
	i16ValSlice := Int16ValueSlice(i16Slice)
	assertEqual(t, []int16{1}, i16ValSlice)
	assertNil(t, Int16Slice(nil))
	assertNil(t, Int16ValueSlice(nil))

	i32Slice := Int32Slice([]int32{1})
	i32ValSlice := Int32ValueSlice(i32Slice)
	assertEqual(t, []int32{1}, i32ValSlice)
	assertNil(t, Int32Slice(nil))
	assertNil(t, Int32ValueSlice(nil))

	i64Slice := Int64Slice([]int64{1})
	i64ValSlice := Int64ValueSlice(i64Slice)
	assertEqual(t, []int64{1}, i64ValSlice)
	assertNil(t, Int64Slice(nil))
	assertNil(t, Int64ValueSlice(nil))

	ui := Uint(1)
	uiVal := UintValue(ui)
	assertEqual(t, uint(1), uiVal)
	assertEqual(t, uint(0), UintValue(nil))

	ui8 := Uint8(uint8(1))
	ui8Val := Uint8Value(ui8)
	assertEqual(t, uint8(1), ui8Val)
	assertEqual(t, uint8(0), Uint8Value(nil))

	ui16 := Uint16(uint16(1))
	ui16Val := Uint16Value(ui16)
	assertEqual(t, uint16(1), ui16Val)
	assertEqual(t, uint16(0), Uint16Value(nil))

	ui32 := Uint32(uint32(1))
	ui32Val := Uint32Value(ui32)
	assertEqual(t, uint32(1), ui32Val)
	assertEqual(t, uint32(0), Uint32Value(nil))

	ui64 := Uint64(uint64(1))
	ui64Val := Uint64Value(ui64)
	assertEqual(t, uint64(1), ui64Val)
	assertEqual(t, uint64(0), Uint64Value(nil))

	uiSlice := UintSlice([]uint{1})
	uiValSlice := UintValueSlice(uiSlice)
	assertEqual(t, []uint{1}, uiValSlice)
	assertNil(t, UintSlice(nil))
	assertNil(t, UintValueSlice(nil))

	ui8Slice := Uint8Slice([]uint8{1})
	ui8ValSlice := Uint8ValueSlice(ui8Slice)
	assertEqual(t, []uint8{1}, ui8ValSlice)
	assertNil(t, Uint8Slice(nil))
	assertNil(t, Uint8ValueSlice(nil))

	ui16Slice := Uint16Slice([]uint16{1})
	ui16ValSlice := Uint16ValueSlice(ui16Slice)
	assertEqual(t, []uint16{1}, ui16ValSlice)
	assertNil(t, Uint16Slice(nil))
	assertNil(t, Uint16ValueSlice(nil))

	ui32Slice := Uint32Slice([]uint32{1})
	ui32ValSlice := Uint32ValueSlice(ui32Slice)
	assertEqual(t, []uint32{1}, ui32ValSlice)
	assertNil(t, Uint32Slice(nil))
	assertNil(t, Uint32ValueSlice(nil))

	ui64Slice := Uint64Slice([]uint64{1})
	ui64ValSlice := Uint64ValueSlice(ui64Slice)
	assertEqual(t, []uint64{1}, ui64ValSlice)
	assertNil(t, Uint64Slice(nil))
	assertNil(t, Uint64ValueSlice(nil))
}
