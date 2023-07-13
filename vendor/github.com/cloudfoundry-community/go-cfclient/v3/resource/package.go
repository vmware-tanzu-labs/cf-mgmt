package resource

import (
	"encoding/json"
	"fmt"
	"time"
)

type PackageState string

const (
	PackageStateAwaitingUpload   PackageState = "AWAITING_UPLOAD"
	PackageStateProcessingUpload PackageState = "PROCESSING_UPLOAD"
	PackageStateReady            PackageState = "READY"
	PackageStateFailed           PackageState = "FAILED"
	PackageStateCopying          PackageState = "COPYING"
	PackageStateExpired          PackageState = "EXPIRED"
)

type Package struct {
	GUID          string              `json:"guid"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Type          string              `json:"type"` // bits or docker
	DataRaw       json.RawMessage     `json:"data"`
	Data          BitsOrDockerPackage `json:"-"` // depends on value of Type
	State         PackageState        `json:"state"`
	Links         map[string]Link     `json:"links"`
	Relationships AppRelationship     `json:"relationships"`
	Metadata      *Metadata           `json:"metadata"`
}

type PackageCreate struct {
	Type          string          `json:"type"`
	Relationships AppRelationship `json:"relationships"`
	Data          *DockerPackage  `json:"data,omitempty"`
	Metadata      *Metadata       `json:"metadata,omitempty"`
}

type PackageUpdate struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

type PackageList struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []*Package `json:"resources,omitempty"`
}

type PackageCopy struct {
	Relationships AppRelationship `json:"relationships"`
}

type BitsOrDockerPackage struct {
	Bits   *BitsPackage
	Docker *DockerPackage
}

// BitsPackage is the data for Packages of type bits.
// It provides an upload link to which a zip file should be uploaded.
type BitsPackage struct {
	Error    *string             `json:"error"`
	Checksum BitsPackageChecksum `json:"checksum"`
}

type DockerPackage struct {
	Image string `json:"image"`
	*DockerCredentials
}

type BitsPackageChecksum struct {
	Type  string  `json:"type"`  // eg. sha256
	Value *string `json:"value"` // populated after the bits are uploaded
}

type DockerCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func NewPackageCreate(appGUID string) *PackageCreate {
	return &PackageCreate{
		Type: "bits",
		Relationships: AppRelationship{
			App: ToOneRelationship{
				Data: &Relationship{
					GUID: appGUID,
				},
			},
		},
	}
}

func NewDockerPackageCreate(appGUID, image, username, password string) *PackageCreate {
	return &PackageCreate{
		Type: "docker",
		Relationships: AppRelationship{
			App: ToOneRelationship{
				Data: &Relationship{
					GUID: appGUID,
				},
			},
		},
		Data: &DockerPackage{
			Image: image,
			DockerCredentials: &DockerCredentials{
				Username: username,
				Password: password,
			},
		},
	}
}

func NewPackageCopy(appGUID string) *PackageCopy {
	return &PackageCopy{
		Relationships: AppRelationship{
			App: ToOneRelationship{
				Data: &Relationship{
					GUID: appGUID,
				},
			},
		},
	}
}

func (d *Package) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	// new type without any functions to avoid recursive unmarshall call
	type unmarshalPackage Package
	err := json.Unmarshal(data, (*unmarshalPackage)(d))
	if err != nil {
		return err
	}

	// post-processing based on type
	if d.Type == "bits" {
		var p BitsPackage
		err = json.Unmarshal(d.DataRaw, &p)
		if err != nil {
			return err
		}
		d.Data.Bits = &p
		return nil
	} else if d.Type == "docker" {
		var p DockerPackage
		err = json.Unmarshal(d.DataRaw, &p)
		if err != nil {
			return err
		}
		d.Data.Docker = &p
		return nil
	}
	return fmt.Errorf("could not unmarshal data as bits or docker package: %w", err)
}
