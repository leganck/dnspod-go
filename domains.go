package dnspod

import (
	"encoding/json"
	"fmt"
)

const (
	methodDomainList   = "Domain.List"
	methodDomainCreate = "Domain.Create"
	methodDomainInfo   = "Domain.Info"
	methodDomainRemove = "Domain.Remove"
)

// DomainInfo handles domain information.
type DomainInfo struct {
	DomainTotal   json.Number `json:"domain_total,omitempty"`
	AllTotal      json.Number `json:"all_total,omitempty"`
	MineTotal     json.Number `json:"mine_total,omitempty"`
	ShareTotal    json.Number `json:"share_total,omitempty"`
	VipTotal      json.Number `json:"vip_total,omitempty"`
	IsMarkTotal   json.Number `json:"ismark_total,omitempty"`
	PauseTotal    json.Number `json:"pause_total,omitempty"`
	ErrorTotal    json.Number `json:"error_total,omitempty"`
	LockTotal     json.Number `json:"lock_total,omitempty"`
	SpamTotal     json.Number `json:"spam_total,omitempty"`
	VipExpire     json.Number `json:"vip_expire,omitempty"`
	ShareOutTotal json.Number `json:"share_out_total,omitempty"`
}

// Domain handles domain.
type Domain struct {
	ID               json.Number `json:"id,omitempty"`
	Name             string      `json:"name,omitempty"`
	PunyCode         string      `json:"punycode,omitempty"`
	Grade            string      `json:"grade,omitempty"`
	GradeTitle       string      `json:"grade_title,omitempty"`
	Status           string      `json:"status,omitempty"`
	ExtStatus        string      `json:"ext_status,omitempty"`
	Records          string      `json:"records,omitempty"`
	GroupID          json.Number `json:"group_id,omitempty"`
	IsMark           string      `json:"is_mark,omitempty"`
	Remark           string      `json:"remark,omitempty"`
	IsVIP            string      `json:"is_vip,omitempty"`
	SearchenginePush string      `json:"searchengine_push,omitempty"`
	UserID           string      `json:"user_id,omitempty"`
	CreatedOn        string      `json:"created_on,omitempty"`
	UpdatedOn        string      `json:"updated_on,omitempty"`
	TTL              string      `json:"ttl,omitempty"`
	CNameSpeedUp     string      `json:"cname_speedup,omitempty"`
	Owner            string      `json:"owner,omitempty"`
	AuthToAnquanBao  bool        `json:"auth_to_anquanbao,omitempty"`
}

type domainListWrapper struct {
	Status  Status     `json:"status"`
	Info    DomainInfo `json:"info"`
	Domains []Domain   `json:"domains"`
}

type domainWrapper struct {
	Status Status     `json:"status"`
	Info   DomainInfo `json:"info"`
	Domain Domain     `json:"domain"`
}

// DomainsService handles communication with the domain related methods of the DNSPod API.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/domains.html
// - https://docs.dnspod.com/api/
type DomainsService struct {
	client *Client
}

// List the domains.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/domains.html#domain-list
// - https://docs.dnspod.com/api/5fe1b40a6e336701a2111f5b/
func (s *DomainsService) List() ([]Domain, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()

	returnedDomains := domainListWrapper{}

	res, err := s.client.post(methodDomainList, payload, &returnedDomains)
	if err != nil {
		return nil, res, err
	}

	if returnedDomains.Status.Code != "1" {
		return nil, nil, fmt.Errorf("could not get domains: %s", returnedDomains.Status.Message)
	}

	return returnedDomains.Domains, res, nil
}

// Create a new domain.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/domains.html#domain-create
// - https://docs.dnspod.com/api/5fe1a9e36e336701a2111d3d/
func (s *DomainsService) Create(domainAttributes Domain) (Domain, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Set("domain", domainAttributes.Name)
	payload.Set("group_id", domainAttributes.GroupID.String())
	payload.Set("is_mark", domainAttributes.IsMark)

	returnedDomain := domainWrapper{}

	res, err := s.client.post(methodDomainCreate, payload, &returnedDomain)
	if err != nil {
		return Domain{}, res, err
	}

	return returnedDomain.Domain, res, nil
}

// Get fetches a domain.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/domains.html#domain-info
// - https://docs.dnspod.com/api/5fe1b37d6e336701a2111f2b/
func (s *DomainsService) Get(id, domain string) (Domain, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Set("domain_id", id)
	payload.Set("domain", domain)

	returnedDomain := domainWrapper{}

	res, err := s.client.post(methodDomainInfo, payload, &returnedDomain)
	if err != nil {
		return Domain{}, res, err
	}

	return returnedDomain.Domain, res, nil
}

// Delete a domain.
//
// DNSPod API docs:
// - https://dnsapi.cn/Domain.Remove
// - https://docs.dnspod.com/api/5fe1ac446e336701a2111dd1/
func (s *DomainsService) Delete(id string, domain string) (*Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Set("domain_id", id)
	payload.Set("domain", domain)
	returnedDomain := domainWrapper{}

	return s.client.post(methodDomainRemove, payload, &returnedDomain)
}
