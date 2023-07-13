package resource

import "time"

type ServiceBroker struct {
	GUID          string            `json:"guid"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Name          string            `json:"name"`
	URL           string            `json:"url"`
	Relationships SpaceRelationship `json:"relationships"`
	Links         map[string]Link   `json:"links"`
	Metadata      *Metadata         `json:"metadata"`
}

type ServiceBrokerList struct {
	Pagination Pagination       `json:"pagination"`
	Resources  []*ServiceBroker `json:"resources"`
}

type ServiceBrokerCreate struct {
	Name           string                   `json:"name"`
	URL            string                   `json:"url"`
	Authentication ServiceBrokerCredentials `json:"authentication"`

	Relationships *SpaceRelationship `json:"relationships,omitempty"`
	Metadata      *Metadata          `json:"metadata,omitempty"`
}

type ServiceBrokerUpdate struct {
	Name           *string                   `json:"name,omitempty"`
	URL            *string                   `json:"url,omitempty"`
	Authentication *ServiceBrokerCredentials `json:"authentication,omitempty"`
	Metadata       *Metadata                 `json:"metadata,omitempty"`
}

type ServiceBrokerCredentials struct {
	Type        string                            `json:"type"` // basic
	Credentials ServiceBrokerBasicAuthCredentials `json:"credentials"`
}

type ServiceBrokerBasicAuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewServiceBrokerCreate(name, url, username, password string) *ServiceBrokerCreate {
	return &ServiceBrokerCreate{
		Name: name,
		URL:  url,
		Authentication: ServiceBrokerCredentials{
			Type: "basic",
			Credentials: ServiceBrokerBasicAuthCredentials{
				Username: username,
				Password: password,
			},
		},
	}
}

func (s *ServiceBrokerCreate) WithSpace(guid string) *ServiceBrokerCreate {
	s.Relationships = &SpaceRelationship{
		Space: ToOneRelationship{
			Data: &Relationship{
				GUID: guid,
			},
		},
	}
	return s
}

func NewServiceBrokerUpdate() *ServiceBrokerUpdate {
	return &ServiceBrokerUpdate{}
}

func (s *ServiceBrokerUpdate) WithURL(url string) *ServiceBrokerUpdate {
	s.URL = &url
	return s
}

func (s *ServiceBrokerUpdate) WithName(name string) *ServiceBrokerUpdate {
	s.Name = &name
	return s
}

func (s *ServiceBrokerUpdate) WithCredentials(username, password string) *ServiceBrokerUpdate {
	if s.Authentication == nil {
		s.Authentication = &ServiceBrokerCredentials{}
	}
	s.Authentication.Type = "basic"
	s.Authentication.Credentials.Username = username
	s.Authentication.Credentials.Password = password
	return s
}
