package dnspod

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const (
	methodRecordList   = "Record.List"
	methodRecordCreate = "Record.Create"
	methodRecordInfo   = "Record.Info"
	methodRecordRemove = "Record.Remove"
	methodRecordModify = "Record.Modify"
)

// Record is the DNS record representation.
type Record struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Line          string `json:"line,omitempty"`
	LineID        string `json:"line_id,omitempty"`
	Type          string `json:"type,omitempty"`
	TTL           string `json:"ttl,omitempty"`
	Value         string `json:"value,omitempty"`
	MX            string `json:"mx,omitempty"`
	Enabled       string `json:"enabled,omitempty"`
	Status        string `json:"status,omitempty"`
	MonitorStatus string `json:"monitor_status,omitempty"`
	Remark        string `json:"remark,omitempty"`
	UpdateOn      string `json:"updated_on,omitempty"`
	UseAQB        string `json:"use_aqb,omitempty"`
	Weight        *int   `json:"weight,omitempty"`
}

// RecordModify is the DNS record modify representation.
type RecordModify struct {
	ID     json.Number `json:"id,omitempty"`
	Name   string      `json:"name,omitempty"`
	Value  string      `json:"value,omitempty"`
	Status string      `json:"status,omitempty"`
}

type recordsWrapper struct {
	Status  Status     `json:"status"`
	Info    DomainInfo `json:"info"`
	Records []Record   `json:"records"`
}

type recordWrapper struct {
	Status Status     `json:"status"`
	Info   DomainInfo `json:"info"`
	Record Record     `json:"record"`
}

type recordModifyWrapper struct {
	Status Status       `json:"status"`
	Record RecordModify `json:"record"`
}

// RecordsService handles communication with the DNS records related methods of the dnspod API.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/records.html
// - https://docs.dnspod.com/api/
type RecordsService struct {
	client *Client
}

type RecordParam struct {
	DomainID string `json:"domain_id"`
	Domain   string `json:"domain"`
}

func (p RecordParam) ToURLValues(values url.Values) url.Values {

	return values
}

// List List the domain records.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/records.html#record-list
// - https://docs.dnspod.com/api/5fe19a7a6e336701a2111bb9/

type ListParams struct {
	*RecordParam
	Offset       string `json:"offset"`
	Length       string `json:"length"`
	SubDomain    string `json:"sub_domain"`
	RecordType   string `json:"record_type"`
	RecordLine   string `json:"record_line"`
	RecordLineID string `json:"record_line_id"`
	keyword      string `json:"keyword"`
}

func (p ListParams) ToURLValues(values url.Values) url.Values {
	p.RecordParam.ToURLValues(values)
	if p.Offset != "" {
		values.Add("offset", p.Offset)
	}
	if p.Length != "" {
		values.Add("length", p.Length)
	}
	if p.SubDomain != "" {
		values.Add("sub_domain", p.SubDomain)
	}
	if p.RecordType != "" {
		values.Add("record_type", p.RecordType)
	}
	if p.RecordLine != "" {
		values.Add("record_line", p.RecordLine)
	}
	if p.RecordLineID != "" {
		values.Add("record_line_id", p.RecordLineID)
	}
	if p.keyword != "" {
		values.Add("keyword", p.keyword)
	}

	return values
}

func (s *RecordsService) List(p ListParams) ([]Record, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload = p.ToURLValues(payload)
	if p.DomainID != "" {
		payload.Add("domain_id", p.DomainID)
	}
	if p.Domain != "" {
		payload.Add("domain", p.Domain)
	}
	if p.RecordType != "" {
		payload.Add("record_type", p.RecordType)
	}

	wrappedRecords := recordsWrapper{}

	res, err := s.client.post(methodRecordList, payload, &wrappedRecords)
	if err != nil {
		return nil, res, err
	}

	if wrappedRecords.Status.Code != "1" && wrappedRecords.Status.Code != "10" {
		return nil, nil, fmt.Errorf("could not get domains: %s", wrappedRecords.Status.Message)
	}

	return wrappedRecords.Records, res, nil
}

// Create Creates a domain record.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/records.html#record-create
// - https://docs.dnspod.com/api/5fe19a3f6e336701a2111bb0/
func (s *RecordsService) Create(domain, domainId string, recordAttributes Record) (Record, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Add("domain", domain)
	payload.Add("domain_id", domainId)

	if recordAttributes.Name != "" {
		payload.Add("sub_domain", recordAttributes.Name)
	}

	if recordAttributes.Type != "" {
		payload.Add("record_type", recordAttributes.Type)
	}

	if recordAttributes.Line != "" {
		payload.Add("record_line", recordAttributes.Line)
	}

	if recordAttributes.LineID != "" {
		payload.Add("record_line_id", recordAttributes.LineID)
	}

	if recordAttributes.Value != "" {
		payload.Add("value", recordAttributes.Value)
	}

	if recordAttributes.MX != "" {
		payload.Add("mx", recordAttributes.MX)
	}

	if recordAttributes.TTL != "" {
		payload.Add("ttl", recordAttributes.TTL)
	}

	if recordAttributes.Status != "" {
		payload.Add("status", recordAttributes.Status)
	}

	if recordAttributes.Weight != nil {
		payload.Add("weight", strconv.Itoa(*recordAttributes.Weight))
	}

	returnedRecord := recordWrapper{}

	res, err := s.client.post(methodRecordCreate, payload, &returnedRecord)
	if err != nil {
		return Record{}, res, err
	}

	if returnedRecord.Status.Code != "1" {
		return returnedRecord.Record, nil, fmt.Errorf("could not create domains: %s", returnedRecord.Status.Message)
	}

	return returnedRecord.Record, res, nil
}

// Get Fetches the domain record.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/records.html#record-info
// - https://docs.dnspod.com/api/5fe1a2a06e336701a2111bcd/
func (s *RecordsService) Get(domain, domainId string, recordID int) (Record, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Add("domain", domain)
	payload.Add("domain_id", domainId)
	payload.Add("record_id", strconv.Itoa(recordID))

	returnedRecord := recordWrapper{}

	res, err := s.client.post(methodRecordInfo, payload, &returnedRecord)
	if err != nil {
		return Record{}, res, err
	}

	if returnedRecord.Status.Code != "1" {
		return returnedRecord.Record, nil, fmt.Errorf("could not get domains: %s", returnedRecord.Status.Message)
	}

	return returnedRecord.Record, res, nil
}

// Update Updates a domain record.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/records.html#record-modify
// - https://docs.dnspod.com/api/5fe1a5a16e336701a2111c76/
func (s *RecordsService) Update(domainId, domain, recordID string, recordAttributes Record) (RecordModify, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Add("domain_id", domainId)
	payload.Add("domain", domain)
	payload.Add("record_id", recordID)

	if recordAttributes.Name != "" {
		payload.Add("sub_domain", recordAttributes.Name)
	}

	if recordAttributes.Type != "" {
		payload.Add("record_type", recordAttributes.Type)
	}

	if recordAttributes.Line != "" {
		payload.Add("record_line", recordAttributes.Line)
	}

	if recordAttributes.LineID != "" {
		payload.Add("record_line_id", recordAttributes.LineID)
	}

	if recordAttributes.Value != "" {
		payload.Add("value", recordAttributes.Value)
	}

	if recordAttributes.MX != "" {
		payload.Add("mx", recordAttributes.MX)
	}

	if recordAttributes.TTL != "" {
		payload.Add("ttl", recordAttributes.TTL)
	}

	if recordAttributes.Status != "" {
		payload.Add("status", recordAttributes.Status)
	}

	if recordAttributes.Weight != nil {
		payload.Add("weight", strconv.Itoa(*recordAttributes.Weight))
	}

	returnedRecord := recordModifyWrapper{}

	res, err := s.client.post(methodRecordModify, payload, &returnedRecord)
	if err != nil {
		return RecordModify{}, res, err
	}

	if returnedRecord.Status.Code != "1" {
		return returnedRecord.Record, nil, fmt.Errorf("could not get domains: %s", returnedRecord.Status.Message)
	}

	return returnedRecord.Record, res, nil
}

// Delete Deletes a domain record.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/records.html#record-remove
// - https://docs.dnspod.com/api/5fe1a4576e336701a2111c24/
func (s *RecordsService) Delete(domainId int, domain, recordID string) (*Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Add("domain_id", strconv.Itoa(domainId))
	payload.Add("domain", domain)
	payload.Add("record_id", recordID)

	returnedRecord := recordWrapper{}

	res, err := s.client.post(methodRecordRemove, payload, &returnedRecord)
	if err != nil {
		return res, err
	}

	if returnedRecord.Status.Code != "1" {
		return nil, fmt.Errorf("could not get domains: %s", returnedRecord.Status.Message)
	}

	return res, nil
}
