package util

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Inner struct {
	Field1 string `customtag:"update1"`
	Field2 string `customtag:"update2"`
}

type Outer struct {
	Inner // Embedded struct

	Field3 string `customtag:"update3"`

	Field4 Inner
	Field5 *Inner

	Field6 Inner
	Field7 *Inner
}

func TestProcess_1(t *testing.T) {
	o := &Outer{
		Inner: Inner{
			Field1: "value1",
			Field2: "42",
		},
		Field3: "value3",
		Field4: Inner{
			Field1: "value4",
			Field2: "43",
		},
		Field5: &Inner{
			Field1: "value5",
			Field2: "44",
		},
	}

	err := Process("customtag", o, func(v reflect.Value, field reflect.StructField, tag string) error {
		if v.Kind() == reflect.String {
			v.SetString(tag)
		}

		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &Outer{
		Inner: Inner{
			Field1: "update1",
			Field2: "update2",
		},
		Field3: "update3",
		Field4: Inner{
			Field1: "update1",
			Field2: "update2",
		},
		Field5: &Inner{
			Field1: "update1",
			Field2: "update2",
		},
		Field6: Inner{
			Field1: "update1",
			Field2: "update2",
		},
	}, o)

	a := &Outer{}

	err1 := Process("customtag", a, func(v reflect.Value, field reflect.StructField, tag string) error {
		if v.Kind() == reflect.String {
			v.SetString(tag)
		}

		return nil
	})

	assert.NoError(t, err1)
	assert.Equal(t, &Outer{
		Inner: Inner{
			Field1: "update1",
			Field2: "update2",
		},
		Field3: "update3",
		Field4: Inner{
			Field1: "update1",
			Field2: "update2",
		},
		Field6: Inner{
			Field1: "update1",
			Field2: "update2",
		},
	}, a)

	b := Outer{}
	err2 := Process("customtag", b, func(v reflect.Value, field reflect.StructField, tag string) error {
		if v.Kind() == reflect.String {
			v.SetString(tag)
		}

		return nil
	})

	assert.Error(t, err2)
}

func TestProcess_2(t *testing.T) {
	t.Run("normal struct with pointers", func(t *testing.T) {
		o := &Outer{
			Inner: Inner{
				Field1: "value1",
				Field2: "42",
			},
			Field3: "value3",
			Field4: Inner{
				Field1: "value4",
				Field2: "43",
			},
			Field5: &Inner{
				Field1: "value5",
				Field2: "44",
			},
			Field6: Inner{
				Field1: "value6",
				Field2: "45",
			},
			Field7: &Inner{
				Field1: "value7",
				Field2: "46",
			},
		}

		err := Process("customtag", o, func(v reflect.Value, field reflect.StructField, tag string) error {
			if v.Kind() == reflect.String {
				v.SetString(tag)
			}

			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, &Outer{
			Inner: Inner{
				Field1: "update1",
				Field2: "update2",
			},
			Field3: "update3",
			Field4: Inner{
				Field1: "update1",
				Field2: "update2",
			},
			Field5: &Inner{
				Field1: "update1",
				Field2: "update2",
			},
			Field6: Inner{
				Field1: "update1",
				Field2: "update2",
			},
			Field7: &Inner{
				Field1: "update1",
				Field2: "update2",
			},
		}, o)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var o *Outer
		err := Process("customtag", &o, func(v reflect.Value, field reflect.StructField, tag string) error {
			if v.Kind() == reflect.String {
				v.SetString(tag)
			}

			return nil
		})

		assert.Error(t, err)
		assert.Nil(t, o)
	})

	t.Run("nil pointer in struct", func(t *testing.T) {
		o := &Outer{}

		err := Process("customtag", o, func(v reflect.Value, field reflect.StructField, tag string) error {
			if v.Kind() == reflect.String {
				v.SetString(tag)
			}

			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, &Outer{
			Inner: Inner{
				Field1: "update1",
				Field2: "update2",
			},
			Field3: "update3",
			Field4: Inner{
				Field1: "update1",
				Field2: "update2",
			},
			Field5: (*Inner)(nil),
			Field6: Inner{
				Field1: "update1",
				Field2: "update2",
			},
			Field7: (*Inner)(nil),
		}, o)
	})

	t.Run("non-pointer struct", func(t *testing.T) {
		o := Outer{
			Inner: Inner{
				Field1: "value1",
				Field2: "42",
			},
			Field3: "value3",
		}

		err := Process("customtag", o, func(v reflect.Value, field reflect.StructField, tag string) error {
			if v.Kind() == reflect.String {
				v.SetString(tag)
			}

			return nil
		})

		assert.Error(t, err)
		assert.Equal(t, Outer{
			Inner: Inner{
				Field1: "value1",
				Field2: "42",
			},
			Field3: "value3",
		}, o)
	})

	t.Run("empty struct", func(t *testing.T) {
		o := &Outer{}

		err := Process("customtag", o, func(v reflect.Value, field reflect.StructField, tag string) error {
			if v.Kind() == reflect.String {
				v.SetString(tag)
			}

			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, &Outer{
			Inner: Inner{
				Field1: "update1",
				Field2: "update2",
			},
			Field3: "update3",
			Field4: Inner{
				Field1: "update1",
				Field2: "update2",
			},
			Field5: (*Inner)(nil),
			Field6: Inner{
				Field1: "update1",
				Field2: "update2",
			},
			Field7: (*Inner)(nil),
		}, o)
	})
}

func TestProcess_3(t *testing.T) {
	o := &Outer{}

	err := Process("customtag", o, func(v reflect.Value, field reflect.StructField, tag string) error {
		if err := setValueFromTag(v, field, tag, tag); err != nil {
			return err
		}

		return nil
	})

	assert.NoError(t, err)
}

type timeDurationStruct struct {
	TimeField     time.Time     `customtag:"2022-01-01"`
	DurationField time.Duration `customtag:"1h"`
}

func TestProcess_TimeDurationStruct(t *testing.T) {
	tds := &timeDurationStruct{}
	err := Process("customtag", tds, func(v reflect.Value, field reflect.StructField, tag string) error {
		if err := setValueFromTag(v, field, tag, tag); err != nil {
			return err
		}

		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), tds.TimeField)
	assert.Equal(t, time.Hour, tds.DurationField)
}

func TestProcess_PrimitivesStruct(t *testing.T) {
	type TestStruct struct {
		Field1 string        `customtag:"field1"`
		Field2 int           `customtag:"42"`
		Field3 uint          `customtag:"43"`
		Field4 float64       `customtag:"3.14"`
		Field5 bool          `customtag:"true"`
		Field6 time.Time     `customtag:"2023-03-18T00:00:00Z"`
		Field7 time.Duration `customtag:"1h"`
	}

	ts := &TestStruct{}

	err := Process("customtag", ts, func(v reflect.Value, field reflect.StructField, tag string) error {
		if err := setValueFromTag(v, field, tag, tag); err != nil {
			return err
		}

		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &TestStruct{
		Field1: "field1",
		Field2: 42,
		Field3: 43,
		Field4: 3.14,
		Field5: true,
		Field6: time.Date(2023, 3, 18, 0, 0, 0, 0, time.UTC),
		Field7: time.Hour,
	}, ts)
}

func TestProcess_SliceStruct(t *testing.T) {
	type TestStruct struct {
		Field1 []string        `customtag:"a,b,c"`
		Field2 []int           `customtag:"1,2,3"`
		Field3 []uint          `customtag:"4,5,6"`
		Field4 []float64       `customtag:"1.1,2.2,3.3"`
		Field5 []bool          `customtag:"true,false,true"`
		Field6 []time.Time     `customtag:"2023-03-18T00:00:00Z,2023-03-19T00:00:00Z"`
		Field7 []time.Duration `customtag:"1h,2m,3ms"`
	}

	ts := &TestStruct{}

	err := Process("customtag", ts, func(v reflect.Value, field reflect.StructField, tag string) error {
		if err := setValueFromTag(v, field, tag, tag); err != nil {
			return err
		}

		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &TestStruct{
		Field1: []string{"a", "b", "c"},
		Field2: []int{1, 2, 3},
		Field3: []uint{4, 5, 6},
		Field4: []float64{1.1, 2.2, 3.3},
		Field5: []bool{true, false, true},
		Field6: []time.Time{
			time.Date(2023, 3, 18, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 3, 19, 0, 0, 0, 0, time.UTC),
		},
		Field7: []time.Duration{time.Hour, 2 * time.Minute, 3 * time.Millisecond},
	}, ts)
}

func TestProcess_MapStruct(t *testing.T) {
	type TestStruct struct {
		Field1  map[string]string        `customtag:"a:a,b:b,c:c"`
		Field2  map[string]int           `customtag:"a:1,b:2,c:3"`
		Field3  map[string]uint          `customtag:"a:4,b:5,c:6"`
		Field4  map[string]float64       `customtag:"a:1.1,b:2.2,c:3.3"`
		Field5  map[string]bool          `customtag:"a:true,b:false"`
		Field6  map[string]time.Time     `customtag:"a:2021-03-18,b:2022-03-18"`
		Field7  map[string]time.Duration `customtag:"a:1h,b:2s,c:3ms"`
		Field8  map[string]interface{}   `customtag:"asd:qwe,dfg:1"`
		Field9  map[string]interface{}   `customtag:"asd:qwe,dfg:text1"`
		Field10 map[string]interface{}   `customtag:"asd:qwe,dfg:true"`
		Field11 map[string]interface{}   `customtag:"asd:qwe,dfg:false"`
		Field12 map[string]interface{}   `customtag:"asd:qwe,dfg:0.65"`
		Field13 map[string]interface{}   `customtag:"asd:qwe,dfg:0"`
		Field14 map[string]interface{}   `customtag:"asd:qwe,dfg:123"`
	}

	ts := &TestStruct{}

	err := Process("customtag", ts, func(v reflect.Value, field reflect.StructField, tag string) error {
		if err := setValueFromTag(v, field, tag, tag); err != nil {
			return err
		}

		return nil
	})

	assert.NoError(t, err)

	assert.EqualValues(t, &TestStruct{
		Field1: map[string]string{"a": "a", "b": "b", "c": "c"},
		Field2: map[string]int{"a": 1, "b": 2, "c": 3},
		Field3: map[string]uint{"a": 4, "b": 5, "c": 6},
		Field4: map[string]float64{"a": 1.1, "b": 2.2, "c": 3.3},
		Field5: map[string]bool{"a": true, "b": false},
		Field6: map[string]time.Time{
			"a": time.Date(2021, 3, 18, 0, 0, 0, 0, time.UTC),
			"b": time.Date(2022, 3, 18, 0, 0, 0, 0, time.UTC),
		},
		Field7:  map[string]time.Duration{"a": time.Hour, "b": 2 * time.Second, "c": 3 * time.Millisecond},
		Field8:  map[string]interface{}{"asd": "qwe", "dfg": int64(1)},
		Field9:  map[string]interface{}{"asd": "qwe", "dfg": "text1"},
		Field10: map[string]interface{}{"asd": "qwe", "dfg": true},
		Field11: map[string]interface{}{"asd": "qwe", "dfg": false},
		Field12: map[string]interface{}{"asd": "qwe", "dfg": 0.65},
		Field13: map[string]interface{}{"asd": "qwe", "dfg": int64(0)},
		Field14: map[string]interface{}{"asd": "qwe", "dfg": int64(123)},
	}, ts)
}