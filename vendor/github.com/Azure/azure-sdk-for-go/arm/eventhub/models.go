package eventhub

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 0.17.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
)

// AccessRights enumerates the values for access rights.
type AccessRights string

const (
	// Listen specifies the listen state for access rights.
	Listen AccessRights = "Listen"
	// Manage specifies the manage state for access rights.
	Manage AccessRights = "Manage"
	// Send specifies the send state for access rights.
	Send AccessRights = "Send"
)

// EntityStatus enumerates the values for entity status.
type EntityStatus string

const (
	// Active specifies the active state for entity status.
	Active EntityStatus = "Active"
	// Creating specifies the creating state for entity status.
	Creating EntityStatus = "Creating"
	// Deleting specifies the deleting state for entity status.
	Deleting EntityStatus = "Deleting"
	// Disabled specifies the disabled state for entity status.
	Disabled EntityStatus = "Disabled"
	// ReceiveDisabled specifies the receive disabled state for entity status.
	ReceiveDisabled EntityStatus = "ReceiveDisabled"
	// Renaming specifies the renaming state for entity status.
	Renaming EntityStatus = "Renaming"
	// Restoring specifies the restoring state for entity status.
	Restoring EntityStatus = "Restoring"
	// SendDisabled specifies the send disabled state for entity status.
	SendDisabled EntityStatus = "SendDisabled"
	// Unknown specifies the unknown state for entity status.
	Unknown EntityStatus = "Unknown"
)

// NamespaceState enumerates the values for namespace state.
type NamespaceState string

const (
	// NamespaceStateActivating specifies the namespace state activating state
	// for namespace state.
	NamespaceStateActivating NamespaceState = "Activating"
	// NamespaceStateActive specifies the namespace state active state for
	// namespace state.
	NamespaceStateActive NamespaceState = "Active"
	// NamespaceStateCreated specifies the namespace state created state for
	// namespace state.
	NamespaceStateCreated NamespaceState = "Created"
	// NamespaceStateCreating specifies the namespace state creating state for
	// namespace state.
	NamespaceStateCreating NamespaceState = "Creating"
	// NamespaceStateDisabled specifies the namespace state disabled state for
	// namespace state.
	NamespaceStateDisabled NamespaceState = "Disabled"
	// NamespaceStateDisabling specifies the namespace state disabling state
	// for namespace state.
	NamespaceStateDisabling NamespaceState = "Disabling"
	// NamespaceStateEnabling specifies the namespace state enabling state for
	// namespace state.
	NamespaceStateEnabling NamespaceState = "Enabling"
	// NamespaceStateFailed specifies the namespace state failed state for
	// namespace state.
	NamespaceStateFailed NamespaceState = "Failed"
	// NamespaceStateRemoved specifies the namespace state removed state for
	// namespace state.
	NamespaceStateRemoved NamespaceState = "Removed"
	// NamespaceStateRemoving specifies the namespace state removing state for
	// namespace state.
	NamespaceStateRemoving NamespaceState = "Removing"
	// NamespaceStateSoftDeleted specifies the namespace state soft deleted
	// state for namespace state.
	NamespaceStateSoftDeleted NamespaceState = "SoftDeleted"
	// NamespaceStateSoftDeleting specifies the namespace state soft deleting
	// state for namespace state.
	NamespaceStateSoftDeleting NamespaceState = "SoftDeleting"
	// NamespaceStateUnknown specifies the namespace state unknown state for
	// namespace state.
	NamespaceStateUnknown NamespaceState = "Unknown"
)

// Policykey enumerates the values for policykey.
type Policykey string

const (
	// PrimaryKey specifies the primary key state for policykey.
	PrimaryKey Policykey = "PrimaryKey"
	// SecondaryKey specifies the secondary key state for policykey.
	SecondaryKey Policykey = "SecondaryKey"
)

// SkuName enumerates the values for sku name.
type SkuName string

const (
	// Basic specifies the basic state for sku name.
	Basic SkuName = "Basic"
	// Premium specifies the premium state for sku name.
	Premium SkuName = "Premium"
	// Standard specifies the standard state for sku name.
	Standard SkuName = "Standard"
)

// SkuTier enumerates the values for sku tier.
type SkuTier string

const (
	// SkuTierBasic specifies the sku tier basic state for sku tier.
	SkuTierBasic SkuTier = "Basic"
	// SkuTierPremium specifies the sku tier premium state for sku tier.
	SkuTierPremium SkuTier = "Premium"
	// SkuTierStandard specifies the sku tier standard state for sku tier.
	SkuTierStandard SkuTier = "Standard"
)

// ConsumerGroupCreateOrUpdateParameters is parameters supplied to the Create
// Or Update Consumer Group operation.
type ConsumerGroupCreateOrUpdateParameters struct {
	Location                 *string `json:"location,omitempty"`
	Type                     *string `json:"type,omitempty"`
	Name                     *string `json:"name,omitempty"`
	*ConsumerGroupProperties `json:"properties,omitempty"`
}

// ConsumerGroupListResult is the response to the List Consumer Group
// operation.
type ConsumerGroupListResult struct {
	autorest.Response `json:"-"`
	Value             *[]ConsumerGroupResource `json:"value,omitempty"`
	NextLink          *string                  `json:"nextLink,omitempty"`
}

// ConsumerGroupListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client ConsumerGroupListResult) ConsumerGroupListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// ConsumerGroupProperties is
type ConsumerGroupProperties struct {
	CreatedAt    *date.Time `json:"createdAt,omitempty"`
	EventHubPath *string    `json:"eventHubPath,omitempty"`
	UpdatedAt    *date.Time `json:"updatedAt,omitempty"`
	UserMetadata *string    `json:"userMetadata,omitempty"`
}

// ConsumerGroupResource is description of the consumer group resource.
type ConsumerGroupResource struct {
	autorest.Response        `json:"-"`
	ID                       *string             `json:"id,omitempty"`
	Name                     *string             `json:"name,omitempty"`
	Type                     *string             `json:"type,omitempty"`
	Location                 *string             `json:"location,omitempty"`
	Tags                     *map[string]*string `json:"tags,omitempty"`
	*ConsumerGroupProperties `json:"properties,omitempty"`
}

// CreateOrUpdateParameters is parameters supplied to the Create Or Update
// Event Hub operation.
type CreateOrUpdateParameters struct {
	Location   *string     `json:"location,omitempty"`
	Type       *string     `json:"type,omitempty"`
	Name       *string     `json:"name,omitempty"`
	Properties *Properties `json:"properties,omitempty"`
}

// ListResult is the response of the List Event Hubs operation.
type ListResult struct {
	autorest.Response `json:"-"`
	Value             *[]ResourceType `json:"value,omitempty"`
	NextLink          *string         `json:"nextLink,omitempty"`
}

// ListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client ListResult) ListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// NamespaceCreateOrUpdateParameters is parameters supplied to the Create Or
// Update Namespace operation.
type NamespaceCreateOrUpdateParameters struct {
	Location             *string             `json:"location,omitempty"`
	Sku                  *Sku                `json:"sku,omitempty"`
	Tags                 *map[string]*string `json:"tags,omitempty"`
	*NamespaceProperties `json:"properties,omitempty"`
}

// NamespaceListResult is the response of the List Namespace operation.
type NamespaceListResult struct {
	autorest.Response `json:"-"`
	Value             *[]NamespaceResource `json:"value,omitempty"`
	NextLink          *string              `json:"nextLink,omitempty"`
}

// NamespaceListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client NamespaceListResult) NamespaceListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// NamespaceProperties is properties of the namespace.
type NamespaceProperties struct {
	ProvisioningState  *string        `json:"provisioningState,omitempty"`
	Status             NamespaceState `json:"status,omitempty"`
	CreatedAt          *date.Time     `json:"createdAt,omitempty"`
	UpdatedAt          *date.Time     `json:"updatedAt,omitempty"`
	ServiceBusEndpoint *string        `json:"serviceBusEndpoint,omitempty"`
	CreateACSNamespace *bool          `json:"createACSNamespace,omitempty"`
	Enabled            *bool          `json:"enabled,omitempty"`
}

// NamespaceResource is description of a namespace resource.
type NamespaceResource struct {
	autorest.Response    `json:"-"`
	ID                   *string             `json:"id,omitempty"`
	Name                 *string             `json:"name,omitempty"`
	Type                 *string             `json:"type,omitempty"`
	Location             *string             `json:"location,omitempty"`
	Tags                 *map[string]*string `json:"tags,omitempty"`
	Sku                  *Sku                `json:"sku,omitempty"`
	*NamespaceProperties `json:"properties,omitempty"`
}

// Properties is
type Properties struct {
	CreatedAt              *date.Time   `json:"createdAt,omitempty"`
	MessageRetentionInDays *int64       `json:"messageRetentionInDays,omitempty"`
	PartitionCount         *int64       `json:"partitionCount,omitempty"`
	PartitionIds           *[]string    `json:"partitionIds,omitempty"`
	Status                 EntityStatus `json:"status,omitempty"`
	UpdatedAt              *date.Time   `json:"updatedAt,omitempty"`
}

// RegenerateKeysParameters is parameters supplied to the Regenerate
// Authorization Rule operation.
type RegenerateKeysParameters struct {
	Policykey Policykey `json:"Policykey,omitempty"`
}

// Resource is
type Resource struct {
	ID       *string             `json:"id,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Type     *string             `json:"type,omitempty"`
	Location *string             `json:"location,omitempty"`
	Tags     *map[string]*string `json:"tags,omitempty"`
}

// ResourceListKeys is namespace/EventHub Connection String
type ResourceListKeys struct {
	autorest.Response         `json:"-"`
	PrimaryConnectionString   *string `json:"primaryConnectionString,omitempty"`
	SecondaryConnectionString *string `json:"secondaryConnectionString,omitempty"`
	PrimaryKey                *string `json:"primaryKey,omitempty"`
	SecondaryKey              *string `json:"secondaryKey,omitempty"`
	KeyName                   *string `json:"keyName,omitempty"`
}

// ResourceType is description of the Event Hub resource.
type ResourceType struct {
	autorest.Response `json:"-"`
	ID                *string             `json:"id,omitempty"`
	Name              *string             `json:"name,omitempty"`
	Type              *string             `json:"type,omitempty"`
	Location          *string             `json:"location,omitempty"`
	Tags              *map[string]*string `json:"tags,omitempty"`
	*Properties       `json:"properties,omitempty"`
}

// SharedAccessAuthorizationRuleCreateOrUpdateParameters is parameters
// supplied to the Create Or Update Authorization Rules operation.
type SharedAccessAuthorizationRuleCreateOrUpdateParameters struct {
	Location                                 *string `json:"location,omitempty"`
	Name                                     *string `json:"name,omitempty"`
	*SharedAccessAuthorizationRuleProperties `json:"properties,omitempty"`
}

// SharedAccessAuthorizationRuleListResult is the response from the List
// Namespace operation.
type SharedAccessAuthorizationRuleListResult struct {
	autorest.Response `json:"-"`
	Value             *[]SharedAccessAuthorizationRuleResource `json:"value,omitempty"`
	NextLink          *string                                  `json:"nextLink,omitempty"`
}

// SharedAccessAuthorizationRuleListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client SharedAccessAuthorizationRuleListResult) SharedAccessAuthorizationRuleListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// SharedAccessAuthorizationRuleProperties is sharedAccessAuthorizationRule
// properties.
type SharedAccessAuthorizationRuleProperties struct {
	Rights *[]AccessRights `json:"rights,omitempty"`
}

// SharedAccessAuthorizationRuleResource is description of a namespace
// authorization rule.
type SharedAccessAuthorizationRuleResource struct {
	autorest.Response                        `json:"-"`
	ID                                       *string             `json:"id,omitempty"`
	Name                                     *string             `json:"name,omitempty"`
	Type                                     *string             `json:"type,omitempty"`
	Location                                 *string             `json:"location,omitempty"`
	Tags                                     *map[string]*string `json:"tags,omitempty"`
	*SharedAccessAuthorizationRuleProperties `json:"properties,omitempty"`
}

// Sku is sKU of the namespace.
type Sku struct {
	Name     SkuName `json:"name,omitempty"`
	Tier     SkuTier `json:"tier,omitempty"`
	Capacity *int32  `json:"capacity,omitempty"`
}
