package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestSchedule_isThisDay(t *testing.T) {
	type fields struct {
		WeekDays []string
		Hour     int
		Minute   int
	}
	type args struct {
		date time.Time
	}
	type testCase struct {
		name   string
		fields fields
		args   args
		want   bool
	}

	tests := []testCase{
		{
			name: "Ok",
			fields: fields{
				WeekDays: []string{"Friday"},
			},
			args: args{date: time.Date(2021, 1, 1, 1, 0, 0, 0, time.Local)},
			want: true,
		},
		{
			name: "Fail",
			fields: fields{
				WeekDays: []string{"Thursday"},
			},
			args: args{date: time.Date(2021, 1, 1, 1, 0, 0, 0, time.Local)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Schedule{
				WeekDays: tt.fields.WeekDays,
				Hour:     tt.fields.Hour,
				Minute:   tt.fields.Minute,
			}
			if got := s.isThisDay(tt.args.date); got != tt.want {
				t.Errorf("isThisDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchedule_GetStartTime(t *testing.T) {
	type fields struct {
		WeekDays []string
		Hour     int
		Minute   int
	}
	type args struct {
		date time.Time
	}
	type testCase struct {
		name   string
		fields fields
		args   args
		want   *time.Time
	}
	okDate := time.Date(2021, 1, 1, 1, 25, 0, 0, time.Local)
	tests := []testCase{
		{
			name: "OK",
			fields: fields{
				WeekDays: []string{"Friday"},
				Hour:     1,
				Minute:   25,
			},
			args: args{date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)},
			want: &okDate,
		},
		{
			name: "Nil",
			fields: fields{
				WeekDays: []string{"Thursday"},
				Hour:     1,
				Minute:   25,
			},
			args: args{date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Schedule{
				WeekDays: tt.fields.WeekDays,
				Hour:     tt.fields.Hour,
				Minute:   tt.fields.Minute,
			}
			if got := s.GetStartTime(tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStartTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_IsRun(t1 *testing.T) {
	type fields struct {
		Schedule    Schedule
		LastRunTime time.Time
	}
	type args struct {
		date time.Time
	}
	type testCase struct {
		name   string
		fields fields
		args   args
		want   bool
	}

	tests := []testCase{
		{
			name: "First start",
			fields: fields{
				Schedule: Schedule{
					WeekDays: []string{"Friday"},
					Hour:     1,
					Minute:   25,
				},
			},
			args: args{date: time.Date(2021, 1, 1, 1, 26, 0, 0, time.Local)},
			want: true,
		},
		{
			name: "Another day",
			fields: fields{
				Schedule: Schedule{
					WeekDays: []string{"Thursday"},
					Hour:     1,
					Minute:   25,
				},
			},
			args: args{date: time.Date(2021, 1, 1, 1, 26, 0, 0, time.Local)},
			want: false,
		},
		{
			name: "Already started at today",
			fields: fields{
				Schedule: Schedule{
					WeekDays: []string{"Friday"},
					Hour:     1,
					Minute:   25,
				},
				LastRunTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			},
			args: args{date: time.Date(2021, 1, 1, 1, 26, 0, 0, time.Local)},
			want: false,
		},
		{
			name: "Start next day",
			fields: fields{
				Schedule: Schedule{
					WeekDays: []string{"Friday"},
					Hour:     1,
					Minute:   25,
				},
				LastRunTime: time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
			},
			args: args{date: time.Date(2021, 1, 1, 1, 26, 0, 0, time.Local)},
			want: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := Task{
				Schedule:    tt.fields.Schedule,
				LastRunTime: tt.fields.LastRunTime,
			}
			if got := t.IsRun(tt.args.date); got != tt.want {
				t1.Errorf("IsRun() = %v, want %v", got, tt.want)
			}
		})
	}
}
