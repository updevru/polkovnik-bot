package domain

import (
	"testing"
	"time"
)

func TestWeekendInterval_IsWeekend(t *testing.T) {
	type fields struct {
		Start time.Time
		End   time.Time
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
			name: "more_5_days_start",
			fields: fields{
				Start: time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
				End:   time.Date(2021, 05, 10, 0, 0, 0, 0, time.Local),
			},
			args: args{
				date: time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "more_5_days_first_day",
			fields: fields{
				Start: time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
				End:   time.Date(2021, 05, 10, 0, 0, 0, 0, time.Local),
			},
			args: args{
				date: time.Date(2021, 05, 01, 1, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "more_5_days_last_day",
			fields: fields{
				Start: time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
				End:   time.Date(2021, 05, 10, 0, 0, 0, 0, time.Local),
			},
			args: args{
				date: time.Date(2021, 05, 10, 0, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "more_5_days_end",
			fields: fields{
				Start: time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
				End:   time.Date(2021, 05, 10, 0, 0, 0, 0, time.Local),
			},
			args: args{
				date: time.Date(2021, 05, 11, 0, 0, 0, 0, time.Local),
			},
			want: false,
		},
		{
			name: "one_day_start",
			fields: fields{
				Start: time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
				End:   time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
			},
			args: args{
				date: time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "one_day",
			fields: fields{
				Start: time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
				End:   time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
			},
			args: args{
				date: time.Date(2021, 05, 01, 5, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "one_day_end",
			fields: fields{
				Start: time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
				End:   time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
			},
			args: args{
				date: time.Date(2021, 06, 01, 5, 0, 0, 0, time.Local),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := WeekendInterval{
				Start: tt.fields.Start,
				End:   tt.fields.End,
			}
			if got := i.IsWeekend(tt.args.date); got != tt.want {
				t.Errorf("IsWeekend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeekend_IsWeekend(t *testing.T) {
	type fields struct {
		WeekDays  []string
		Intervals []WeekendInterval
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
			name: "workday",
			fields: fields{
				WeekDays:  []string{"Saturday", "Sunday"},
				Intervals: nil,
			},
			args: args{
				date: time.Date(2021, 05, 12, 5, 0, 0, 0, time.Local),
			},
			want: false,
		},
		{
			name: "weekend_saturday",
			fields: fields{
				WeekDays:  []string{"Saturday", "Sunday"},
				Intervals: nil,
			},
			args: args{
				date: time.Date(2021, 05, 15, 5, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "weekend_sunday",
			fields: fields{
				WeekDays:  []string{"Saturday", "Sunday"},
				Intervals: nil,
			},
			args: args{
				date: time.Date(2021, 05, 16, 5, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "weekend_interval",
			fields: fields{
				WeekDays: nil,
				Intervals: []WeekendInterval{
					{
						Start: time.Date(2021, 05, 01, 0, 0, 0, 0, time.Local),
						End:   time.Date(2021, 05, 10, 0, 0, 0, 0, time.Local),
					},
				},
			},
			args: args{
				date: time.Date(2021, 05, 05, 5, 0, 0, 0, time.Local),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Weekend{
				WeekDays:  tt.fields.WeekDays,
				Intervals: tt.fields.Intervals,
			}
			if got := w.IsWeekend(tt.args.date); got != tt.want {
				t.Errorf("IsWeekend() = %v, want %v", got, tt.want)
			}
		})
	}
}
