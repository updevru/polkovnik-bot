package domain

import (
	"reflect"
	"testing"
)

func TestNewTime(t *testing.T) {
	type args struct {
		seconds int
	}
	type testCase struct {
		name string
		args args
		want Time
	}
	tests := []testCase{
		{
			name: "Create time object",
			args: args{seconds: 1612415944},
			want: Time{seconds: 1612415944},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTime(tt.args.seconds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_Seconds(t1 *testing.T) {
	type fields struct {
		seconds int
	}
	type testCase struct {
		name   string
		fields fields
		want   int
	}
	tests := []testCase{
		{
			name:   "Get seconds",
			fields: fields{seconds: 1612415944},
			want:   1612415944,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := Time{
				seconds: tt.fields.seconds,
			}
			if got := t.Seconds(); got != tt.want {
				t1.Errorf("Seconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_ToHumanFormat(t1 *testing.T) {
	type fields struct {
		seconds int
	}
	type testCase struct {
		name   string
		fields fields
		want   string
	}
	tests := []testCase{
		{
			name:   "Hours 1h",
			fields: fields{seconds: 3600},
			want:   "1h",
		},
		{
			name:   "Zero hours",
			fields: fields{seconds: 0},
			want:   "0h",
		},
		{
			name:   "Hours 1h 30m",
			fields: fields{seconds: 3600 + 1800},
			want:   "1h 30m",
		},
		{
			name:   "Days 1d 7h",
			fields: fields{seconds: 86400 + 25200 + 1800},
			want:   "1d 7h 30m",
		},
		{
			name:   "Days 2w 5d 7h 8m",
			fields: fields{seconds: 14*86400 + 5*86400 + 7*3600 + 8*60},
			want:   "2w 5d 7h 8m",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := Time{
				seconds: tt.fields.seconds,
			}
			if got := t.ToHumanFormat(); got != tt.want {
				t1.Errorf("ToHumanFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
