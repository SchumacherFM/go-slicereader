package slicereader

import (
	"fmt"
	"reflect"
	"testing"
)

type testSlice struct {
	name  string
	slice []string
}

var testSlices = []testSlice{
	{
		"empty slice",
		nil,
	},
	{
		"single value slice",
		[]string{"value1"},
	},
	{
		"two valued slice1",
		[]string{"value1", "value2"},
	},
	{
		"two valued slice2",
		[]string{"value1", "value2", "value3", "value4"},
	},
}

func TestNewSliceReader(t *testing.T) {
	type args[T string] struct {
		slice []T
	}
	tests := []struct {
		name string
		args args[string]
		want *SliceReader[string]
	}{
		{
			"Create reader: " + testSlices[0].name,
			args[string]{testSlices[0].slice},
			&SliceReader[string]{s: testSlices[0].slice, i: 0},
		},
		{
			"Create reader: " + testSlices[1].name,
			args[string]{testSlices[1].slice},
			&SliceReader[string]{s: testSlices[1].slice, i: 0},
		},
		{
			"Create reader: " + testSlices[2].name,
			args[string]{testSlices[2].slice},
			&SliceReader[string]{s: testSlices[2].slice, i: 0},
		},
		{
			"Create reader: " + testSlices[3].name,
			args[string]{testSlices[3].slice},
			&SliceReader[string]{s: testSlices[3].slice, i: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSliceReader(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSliceReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReader_Len(t *testing.T) {
	type fields struct {
		s []string
		i int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			"all unreaded: " + testSlices[0].name,
			fields{testSlices[0].slice, 0},
			0,
		},
		{
			"all unreaded: " + testSlices[1].name,
			fields{testSlices[1].slice, 0},
			1,
		},
		{
			"all unreaded: " + testSlices[2].name,
			fields{testSlices[2].slice, 0},
			2,
		},
		{
			"all unreaded: " + testSlices[3].name,
			fields{testSlices[3].slice, 0},
			4,
		},

		{
			"all readed: " + testSlices[1].name,
			fields{testSlices[1].slice, 1},
			0,
		},
		{
			"all readed: " + testSlices[2].name,
			fields{testSlices[2].slice, 2},
			0,
		},
		{
			"partial readed: " + testSlices[3].name,
			fields{testSlices[3].slice, 3},
			1,
		},
		{
			"all readed: " + testSlices[3].name,
			fields{testSlices[3].slice, 4},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &SliceReader[string]{
				s: tt.fields.s,
				i: tt.fields.i,
			}
			if got := sr.Len(); got != tt.want {
				t.Errorf("SliceReader.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReader_Size(t *testing.T) {
	type fields struct {
		s []string
		i int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			"size == input slice size: " + testSlices[0].name,
			fields{testSlices[0].slice, 0},
			int64(len(testSlices[0].slice)),
		},
		{
			"size == input slice size: " + testSlices[1].name,
			fields{testSlices[1].slice, 1},
			int64(len(testSlices[1].slice)),
		},
		{
			"size == input slice size: " + testSlices[2].name,
			fields{testSlices[2].slice, 1},
			int64(len(testSlices[2].slice)),
		},
		{
			"size == input slice size: " + testSlices[3].name,
			fields{testSlices[3].slice, 4},
			int64(len(testSlices[3].slice)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &SliceReader[string]{
				s: tt.fields.s,
				i: tt.fields.i,
			}
			if got := sr.Size(); got != tt.want {
				t.Errorf("SliceReader.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReader_ReadElement(t *testing.T) {
	type fields struct {
		s []string
		i int64
	}
	tests := []struct {
		name     string
		fields   fields
		wantV    string
		wantErr  error
		wantV2   string
		wantErr2 error
	}{
		{
			"Read twice index 0: " + testSlices[0].name,
			fields{s: testSlices[0].slice, i: 0},
			"",
			EOS,
			"",
			EOS,
		},
		{
			"Read twice from index 0: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			"value1",
			nil,
			"",
			EOS,
		},
		{
			"Read twice from index 0: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			"value1",
			nil,
			"value2",
			nil,
		},
		{
			"Read twice from index 3: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 3},
			"value4",
			nil,
			"",
			EOS,
		},
		{
			"Read twice from index 4: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 4},
			"",
			EOS,
			"",
			EOS,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &SliceReader[string]{
				s: tt.fields.s,
				i: tt.fields.i,
			}
			gotE, err := sr.Read()
			if (err != nil || tt.wantErr != nil) && fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.wantErr) {
				t.Errorf("SliceReader.ReadElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantV) {
				t.Errorf("SliceReader.ReadElement() = %v, want %v", gotE, tt.wantV)
			}
			gotE, err = sr.Read()
			if (err != nil || tt.wantErr2 != nil) && fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.wantErr2) {
				t.Errorf("SliceReader.ReadElement() error = %v, wantErr %v", err, tt.wantErr2)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantV2) {
				t.Errorf("SliceReader.ReadElement() = %v, want %v", gotE, tt.wantV2)
			}
		})
	}
}

func TestSliceReader_ReadWhile(t *testing.T) {
	type fields struct {
		s []string
		i int64
	}
	type args struct {
		v func(string) bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantS   []string
		wantErr error
	}{
		{
			"Read till EOS: " + testSlices[0].name,
			fields{s: testSlices[0].slice, i: 0},
			args{func(v string) bool { return true }},
			testSlices[0].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return true }},
			testSlices[1].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return true }},
			testSlices[2].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return true }},
			testSlices[3].slice,
			EOS,
		},
		{
			"Read till first value: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return v != "value1" }},
			nil,
			nil,
		},
		{
			"Read till first value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v != "value1" }},
			nil,
			nil,
		},
		{
			"Read till second value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v != "value2" }},
			testSlices[2].slice[0:1],
			nil,
		},
		{
			"Read till third value: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return v != "value3" }},
			testSlices[3].slice[0:2],
			nil,
		},
		{
			"Read till last value: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return v != "value1" }},
			nil,
			nil,
		},
		{
			"Read till last value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v != "value2" }},
			testSlices[2].slice[0:1],
			nil,
		},
		{
			"2 Read till last value: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return v != "value4" }},
			testSlices[3].slice[0:3],
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &SliceReader[string]{
				s: tt.fields.s,
				i: tt.fields.i,
			}
			gotS, err := sr.ReadWhile(tt.args.v)
			if (err != nil || tt.wantErr != nil) && fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.wantErr) {
				t.Errorf("SliceReader.ReadWhile() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("SliceReader.ReadWhile() = %#v, want %#v", gotS, tt.wantS)
			}
		})
	}
}

func TestSliceReader_ReadUntil(t *testing.T) {
	type fields struct {
		s []string
		i int64
	}
	type args struct {
		v func(string) bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantS   []string
		wantErr error
	}{
		{
			"Read till EOS: " + testSlices[0].name,
			fields{s: testSlices[0].slice, i: 0},
			args{func(v string) bool { return false }},
			testSlices[0].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return false }},
			testSlices[1].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return false }},
			testSlices[2].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return false }},
			testSlices[3].slice,
			EOS,
		},
		{
			"Read till first value: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return v == "value1" }},
			nil,
			nil,
		},
		{
			"Read till first value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v == "value1" }},
			nil,
			nil,
		},
		{
			"Read till second value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v == "value2" }},
			testSlices[2].slice[0:1],
			nil,
		},
		{
			"Read till third value: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return v == "value3" }},
			testSlices[3].slice[0:2],
			nil,
		},
		{
			"Read till last value: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return v == "value1" }},
			nil,
			nil,
		},
		{
			"Read till last value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v == "value2" }},
			testSlices[2].slice[0:1],
			nil,
		},
		{
			"Read till last value: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return v == "value4" }},
			testSlices[3].slice[0:3],
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &SliceReader[string]{
				s: tt.fields.s,
				i: tt.fields.i,
			}
			gotS, err := sr.ReadUntil(tt.args.v)
			if (err != nil || tt.wantErr != nil) && fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.wantErr) {
				t.Errorf("SliceReader.ReadUntil() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("SliceReader.ReadUntil() = %#v, want %#v", gotS, tt.wantS)
			}
		})
	}
}

func TestSliceReader_ReadWhileIncl(t *testing.T) {
	type fields struct {
		s []string
		i int64
	}
	type args struct {
		v func(string) bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantS   []string
		wantErr error
	}{
		{
			"Read till EOS: " + testSlices[0].name,
			fields{s: testSlices[0].slice, i: 0},
			args{func(v string) bool { return true }},
			testSlices[0].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return true }},
			testSlices[1].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return true }},
			testSlices[2].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return true }},
			testSlices[3].slice,
			EOS,
		},
		{
			"Read till first value: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return v != "value1" }},
			testSlices[1].slice[0:1],
			nil,
		},
		{
			"Read till first value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v != "value1" }},
			testSlices[2].slice[0:1],
			nil,
		},
		{
			"Read till second value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v != "value2" }},
			testSlices[2].slice[0:2],
			nil,
		},
		{
			"Read till third value: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return v != "value3" }},
			testSlices[3].slice[0:3],
			nil,
		},
		{
			"Read till last value: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return v != "value1" }},
			testSlices[1].slice[0:1],
			nil,
		},
		{
			"Read till last value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v != "value2" }},
			testSlices[2].slice[0:2],
			nil,
		},
		{
			"Read till last value: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return v != "value4" }},
			testSlices[3].slice[0:4],
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &SliceReader[string]{
				s: tt.fields.s,
				i: tt.fields.i,
			}
			gotS, err := sr.ReadWhileIncl(tt.args.v)
			if (err != nil || tt.wantErr != nil) && fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.wantErr) {
				t.Errorf("SliceReader.ReadWhileIncl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("SliceReader.ReadWhileIncl() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestSliceReader_ReadUntilIncl(t *testing.T) {
	type fields struct {
		s []string
		i int64
	}
	type args struct {
		v func(string) bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantS   []string
		wantErr error
	}{
		{
			"Read till EOS: " + testSlices[0].name,
			fields{s: testSlices[0].slice, i: 0},
			args{func(v string) bool { return false }},
			testSlices[0].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return false }},
			testSlices[1].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return false }},
			testSlices[2].slice,
			EOS,
		},
		{
			"Read till EOS: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return false }},
			testSlices[3].slice,
			EOS,
		},
		{
			"Read till first value: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return v == "value1" }},
			testSlices[1].slice[0:1],
			nil,
		},
		{
			"Read till first value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v == "value1" }},
			testSlices[2].slice[0:1],
			nil,
		},
		{
			"Read till second value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v == "value2" }},
			testSlices[2].slice[0:2],
			nil,
		},
		{
			"Read till third value: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return v == "value3" }},
			testSlices[3].slice[0:3],
			nil,
		},
		{
			"Read till last value: " + testSlices[1].name,
			fields{s: testSlices[1].slice, i: 0},
			args{func(v string) bool { return v == "value1" }},
			testSlices[1].slice[0:1],
			nil,
		},
		{
			"Read till last value: " + testSlices[2].name,
			fields{s: testSlices[2].slice, i: 0},
			args{func(v string) bool { return v == "value2" }},
			testSlices[2].slice[0:2],
			nil,
		},
		{
			"Read till last value: " + testSlices[3].name,
			fields{s: testSlices[3].slice, i: 0},
			args{func(v string) bool { return v == "value4" }},
			testSlices[3].slice[0:4],
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &SliceReader[string]{
				s: tt.fields.s,
				i: tt.fields.i,
			}
			gotS, err := sr.ReadUntilIncl(tt.args.v)
			if (err != nil || tt.wantErr != nil) && fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.wantErr) {
				t.Errorf("SliceReader.ReadUntilIncl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("SliceReader.ReadUntilIncl() = %#v, want %#v", gotS, tt.wantS)
			}
		})
	}
}
