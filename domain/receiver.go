package domain

import (
	"github.com/google/uuid"
	"time"
)

type ReceiverType string
type ReceiverFormat string
type ReceiverSettings map[string]string

const (
	ReceiverTypeMessage ReceiverType = "send_team_message"
)

const (
	DataReceiverFormatAuto ReceiverFormat = "auto"
	DataReceiverFormatJSON ReceiverFormat = "json"
	DataReceiverFormatXML  ReceiverFormat = "xml"
	DataReceiverFormatTEXT ReceiverFormat = "text"
)

type Receiver struct {
	Id         string
	Active     bool
	Title      string
	Format     ReceiverFormat
	Type       ReceiverType
	Settings   ReceiverSettings
	DateCreate time.Time
	DateUpdate time.Time
}

func (r *Receiver) Edit(title string, status bool, dataType ReceiverType, dataSettings ReceiverSettings, format ReceiverFormat) error {
	r.Title = title
	r.Active = status
	r.Type = dataType
	r.Format = format
	r.Settings = dataSettings
	r.DateUpdate = time.Now()

	return nil
}

func NewReceiver(title string, status bool, dataType ReceiverType, dataSettings ReceiverSettings, format ReceiverFormat) (*Receiver, error) {
	return &Receiver{
		Id:         uuid.NewString(),
		Active:     status,
		Title:      title,
		Type:       dataType,
		Settings:   dataSettings,
		Format:     format,
		DateCreate: time.Now(),
		DateUpdate: time.Now(),
	}, nil
}

func GetReceiverTypes() []ReceiverType {
	return []ReceiverType{ReceiverTypeMessage}
}

func GetReceiverFormats() []ReceiverFormat {
	return []ReceiverFormat{DataReceiverFormatAuto, DataReceiverFormatJSON, DataReceiverFormatXML, DataReceiverFormatTEXT}
}

func GetReceiverType(name string) *ReceiverType {
	for _, t := range GetReceiverTypes() {
		if string(t) == name {
			return &t
		}
	}

	return nil
}

func GetReceiverFormat(name string) *ReceiverFormat {
	for _, t := range GetReceiverFormats() {
		if string(t) == name {
			return &t
		}
	}

	return nil
}
