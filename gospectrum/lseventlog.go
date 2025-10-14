package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type EventLogInstance struct {
	SequenceNumber string          `json:"sequence_number,omitempty"`
	LastTimestamp  types.Timestamp `json:"last_timestamp,omitempty"`
	ObjectType     string          `json:"object_type,omitempty"`
	ObjectId       string          `json:"object_id,omitempty"`
	ObjectName     string          `json:"object_name,omitempty"`
	CopyId         string          `json:"copy_id,omitempty"`
	Status         string          `json:"status,omitempty"`
	Fixed          string          `json:"fixed,omitempty"`
	EventId        string          `json:"event_id,omitempty"`
	ErrorCode      string          `json:"error_code,omitempty"`
	Description    string          `json:"description,omitempty"`
}

type EventLogOptions struct {
	Filtervalue string `json:"filtervalue,omitempty"`
	Alert       *bool  `json:"alert,omitempty"`
	Message     *bool  `json:"message,omitempty"`
	Monitoring  *bool  `json:"monitoring,omitempty"`
	Expired     *bool  `json:"expired,omitempty"`
	Fixed       *bool  `json:"fixed,omitempty"`
}

func NewLsEventLogOptions() *EventLogOptions {
	return &EventLogOptions{}
}

func (_opts *EventLogOptions) AddFilterValue(f string) {
	if _opts.Filtervalue == "" {
		_opts.Filtervalue = f
	} else {
		_opts.Filtervalue += ":" + f
	}
}
func (_opts *EventLogOptions) SetAlert(b bool) {
	_opts.Alert = new(bool)
	*_opts.Alert = b
}

func (_opts *EventLogOptions) SetMessage(b bool) {
	_opts.Message = new(bool)
	*_opts.Message = b
}

func (_opts *EventLogOptions) SetMonitoring(b bool) {
	_opts.Monitoring = new(bool)
	*_opts.Monitoring = b
}

func (_opts *EventLogOptions) SetExpired(b bool) {
	_opts.Expired = new(bool)
	*_opts.Expired = b
}

func (_opts *EventLogOptions) SetFixed(b bool) {
	_opts.Fixed = new(bool)
	*_opts.Fixed = b
}

func (_c *SpectrumClient) GetEventLog(opts *EventLogOptions) ([]*EventLogInstance, error) {
	// Try Login
	err := _c.login()
	if err != nil {
		return nil, err
	}
	reqBody, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	req, err := api.SpectrumAPILsEventLog.NewRequest(_c.endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", _c.token)
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*EventLogInstance
	err = json.Unmarshal(body, &data)

	return data, err
}
