package storage

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
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
)

// AccountsClient is the the Storage Management Client.
type AccountsClient struct {
	ManagementClient
}

// NewAccountsClient creates an instance of the AccountsClient client.
func NewAccountsClient(subscriptionID string) AccountsClient {
	return NewAccountsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewAccountsClientWithBaseURI creates an instance of the AccountsClient
// client.
func NewAccountsClientWithBaseURI(baseURI string, subscriptionID string) AccountsClient {
	return AccountsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CheckNameAvailability checks that the storage account name is valid and is
// not already in use.
//
// accountName is the name of the storage account within the specified
// resource group. Storage account names must be between 3 and 24 characters
// in length and use numbers and lower-case letters only.
func (client AccountsClient) CheckNameAvailability(accountName AccountCheckNameAvailabilityParameters) (result CheckNameAvailabilityResult, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName.Name", Name: validation.Null, Rule: true, Chain: nil},
				{Target: "accountName.Type", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "storage.AccountsClient", "CheckNameAvailability")
	}

	req, err := client.CheckNameAvailabilityPreparer(accountName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "CheckNameAvailability", nil, "Failure preparing request")
	}

	resp, err := client.CheckNameAvailabilitySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "CheckNameAvailability", resp, "Failure sending request")
	}

	result, err = client.CheckNameAvailabilityResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.AccountsClient", "CheckNameAvailability", resp, "Failure responding to request")
	}

	return
}

// CheckNameAvailabilityPreparer prepares the CheckNameAvailability request.
func (client AccountsClient) CheckNameAvailabilityPreparer(accountName AccountCheckNameAvailabilityParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Storage/checkNameAvailability", pathParameters),
		autorest.WithJSON(accountName),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// CheckNameAvailabilitySender sends the CheckNameAvailability request. The method will close the
// http.Response Body if it receives an error.
func (client AccountsClient) CheckNameAvailabilitySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CheckNameAvailabilityResponder handles the response to the CheckNameAvailability request. The method always
// closes the http.Response Body.
func (client AccountsClient) CheckNameAvailabilityResponder(resp *http.Response) (result CheckNameAvailabilityResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Create asynchronously creates a new storage account with the specified
// parameters. If an account is already created and a subsequent create
// request is issued with different properties, the account properties will
// be updated. If an account is already created and a subsequent create or
// update request is issued with the exact same set of properties, the
// request will succeed. This method may poll for completion. Polling can be
// canceled by passing the cancel channel argument. The channel will be used
// to cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group within the user's
// subscription. accountName is the name of the storage account within the
// specified resource group. Storage account names must be between 3 and 24
// characters in length and use numbers and lower-case letters only.
// parameters is the parameters to provide for the created account.
func (client AccountsClient) Create(resourceGroupName string, accountName string, parameters AccountCreateParameters, cancel <-chan struct{}) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.Sku", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "parameters.Sku.Tier", Name: validation.ReadOnly, Rule: true, Chain: nil}}},
				{Target: "parameters.Location", Name: validation.Null, Rule: true, Chain: nil},
				{Target: "parameters.AccountPropertiesCreateParameters", Name: validation.Null, Rule: false,
					Chain: []validation.Constraint{{Target: "parameters.AccountPropertiesCreateParameters.CustomDomain", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "parameters.AccountPropertiesCreateParameters.CustomDomain.Name", Name: validation.Null, Rule: true, Chain: nil}}},
						{Target: "parameters.AccountPropertiesCreateParameters.Encryption", Name: validation.Null, Rule: false,
							Chain: []validation.Constraint{{Target: "parameters.AccountPropertiesCreateParameters.Encryption.Services", Name: validation.Null, Rule: false,
								Chain: []validation.Constraint{{Target: "parameters.AccountPropertiesCreateParameters.Encryption.Services.Blob", Name: validation.Null, Rule: false,
									Chain: []validation.Constraint{{Target: "parameters.AccountPropertiesCreateParameters.Encryption.Services.Blob.LastEnabledTime", Name: validation.ReadOnly, Rule: true, Chain: nil}}},
								}},
								{Target: "parameters.AccountPropertiesCreateParameters.Encryption.KeySource", Name: validation.Null, Rule: true, Chain: nil},
							}},
					}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "storage.AccountsClient", "Create")
	}

	req, err := client.CreatePreparer(resourceGroupName, accountName, parameters, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "Create", nil, "Failure preparing request")
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "Create", resp, "Failure sending request")
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Create", resp, "Failure responding to request")
	}

	return
}

// CreatePreparer prepares the Create request.
func (client AccountsClient) CreatePreparer(resourceGroupName string, accountName string, parameters AccountCreateParameters, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client AccountsClient) CreateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client AccountsClient) CreateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Delete deletes a storage account in Microsoft Azure.
//
// resourceGroupName is the name of the resource group within the user's
// subscription. accountName is the name of the storage account within the
// specified resource group. Storage account names must be between 3 and 24
// characters in length and use numbers and lower-case letters only.
func (client AccountsClient) Delete(resourceGroupName string, accountName string) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "storage.AccountsClient", "Delete")
	}

	req, err := client.DeletePreparer(resourceGroupName, accountName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "Delete", nil, "Failure preparing request")
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "Delete", resp, "Failure sending request")
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client AccountsClient) DeletePreparer(resourceGroupName string, accountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client AccountsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client AccountsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// GetProperties returns the properties for the specified storage account
// including but not limited to name, SKU name, location, and account status.
// The ListKeys operation should be used to retrieve storage keys.
//
// resourceGroupName is the name of the resource group within the user's
// subscription. accountName is the name of the storage account within the
// specified resource group. Storage account names must be between 3 and 24
// characters in length and use numbers and lower-case letters only.
func (client AccountsClient) GetProperties(resourceGroupName string, accountName string) (result Account, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "storage.AccountsClient", "GetProperties")
	}

	req, err := client.GetPropertiesPreparer(resourceGroupName, accountName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "GetProperties", nil, "Failure preparing request")
	}

	resp, err := client.GetPropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "GetProperties", resp, "Failure sending request")
	}

	result, err = client.GetPropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.AccountsClient", "GetProperties", resp, "Failure responding to request")
	}

	return
}

// GetPropertiesPreparer prepares the GetProperties request.
func (client AccountsClient) GetPropertiesPreparer(resourceGroupName string, accountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetPropertiesSender sends the GetProperties request. The method will close the
// http.Response Body if it receives an error.
func (client AccountsClient) GetPropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetPropertiesResponder handles the response to the GetProperties request. The method always
// closes the http.Response Body.
func (client AccountsClient) GetPropertiesResponder(resp *http.Response) (result Account, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List lists all the storage accounts available under the subscription. Note
// that storage keys are not returned; use the ListKeys operation for this.
func (client AccountsClient) List() (result AccountListResult, err error) {
	req, err := client.ListPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "List", nil, "Failure preparing request")
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "List", resp, "Failure sending request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.AccountsClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client AccountsClient) ListPreparer() (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Storage/storageAccounts", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client AccountsClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client AccountsClient) ListResponder(resp *http.Response) (result AccountListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByResourceGroup lists all the storage accounts available under the
// given resource group. Note that storage keys are not returned; use the
// ListKeys operation for this.
//
// resourceGroupName is the name of the resource group within the user's
// subscription.
func (client AccountsClient) ListByResourceGroup(resourceGroupName string) (result AccountListResult, err error) {
	req, err := client.ListByResourceGroupPreparer(resourceGroupName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "ListByResourceGroup", nil, "Failure preparing request")
	}

	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "ListByResourceGroup", resp, "Failure sending request")
	}

	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.AccountsClient", "ListByResourceGroup", resp, "Failure responding to request")
	}

	return
}

// ListByResourceGroupPreparer prepares the ListByResourceGroup request.
func (client AccountsClient) ListByResourceGroupPreparer(resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
// http.Response Body if it receives an error.
func (client AccountsClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client AccountsClient) ListByResourceGroupResponder(resp *http.Response) (result AccountListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListKeys lists the access keys for the specified storage account.
//
// resourceGroupName is the name of the resource group within the user's
// subscription. accountName is the name of the storage account within the
// specified resource group. Storage account names must be between 3 and 24
// characters in length and use numbers and lower-case letters only.
func (client AccountsClient) ListKeys(resourceGroupName string, accountName string) (result AccountListKeysResult, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "storage.AccountsClient", "ListKeys")
	}

	req, err := client.ListKeysPreparer(resourceGroupName, accountName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "ListKeys", nil, "Failure preparing request")
	}

	resp, err := client.ListKeysSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "ListKeys", resp, "Failure sending request")
	}

	result, err = client.ListKeysResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.AccountsClient", "ListKeys", resp, "Failure responding to request")
	}

	return
}

// ListKeysPreparer prepares the ListKeys request.
func (client AccountsClient) ListKeysPreparer(resourceGroupName string, accountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/listKeys", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListKeysSender sends the ListKeys request. The method will close the
// http.Response Body if it receives an error.
func (client AccountsClient) ListKeysSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListKeysResponder handles the response to the ListKeys request. The method always
// closes the http.Response Body.
func (client AccountsClient) ListKeysResponder(resp *http.Response) (result AccountListKeysResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// RegenerateKey regenerates one of the access keys for the specified storage
// account.
//
// resourceGroupName is the name of the resource group within the user's
// subscription. accountName is the name of the storage account within the
// specified resource group. Storage account names must be between 3 and 24
// characters in length and use numbers and lower-case letters only.
// regenerateKey is specifies name of the key which should be regenerated --
// key1 or key2.
func (client AccountsClient) RegenerateKey(resourceGroupName string, accountName string, regenerateKey AccountRegenerateKeyParameters) (result AccountListKeysResult, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: regenerateKey,
			Constraints: []validation.Constraint{{Target: "regenerateKey.KeyName", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "storage.AccountsClient", "RegenerateKey")
	}

	req, err := client.RegenerateKeyPreparer(resourceGroupName, accountName, regenerateKey)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "RegenerateKey", nil, "Failure preparing request")
	}

	resp, err := client.RegenerateKeySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "RegenerateKey", resp, "Failure sending request")
	}

	result, err = client.RegenerateKeyResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.AccountsClient", "RegenerateKey", resp, "Failure responding to request")
	}

	return
}

// RegenerateKeyPreparer prepares the RegenerateKey request.
func (client AccountsClient) RegenerateKeyPreparer(resourceGroupName string, accountName string, regenerateKey AccountRegenerateKeyParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/regenerateKey", pathParameters),
		autorest.WithJSON(regenerateKey),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// RegenerateKeySender sends the RegenerateKey request. The method will close the
// http.Response Body if it receives an error.
func (client AccountsClient) RegenerateKeySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// RegenerateKeyResponder handles the response to the RegenerateKey request. The method always
// closes the http.Response Body.
func (client AccountsClient) RegenerateKeyResponder(resp *http.Response) (result AccountListKeysResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Update the update operation can be used to update the SKU, encryption,
// access tier, or tags for a storage account. It can also be used to map the
// account to a custom domain. Only one custom domain is supported per
// storage account; the replacement/change of custom domain is not supported.
// In order to replace an old custom domain, the old value must be
// cleared/unregistered before a new value can be set. The update of multiple
// properties is supported. This call does not change the storage keys for
// the account. If you want to change the storage account keys, use the
// regenerate keys operation. The location and name of the storage account
// cannot be changed after creation.
//
// resourceGroupName is the name of the resource group within the user's
// subscription. accountName is the name of the storage account within the
// specified resource group. Storage account names must be between 3 and 24
// characters in length and use numbers and lower-case letters only.
// parameters is the parameters to provide for the updated account.
func (client AccountsClient) Update(resourceGroupName string, accountName string, parameters AccountUpdateParameters) (result Account, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "storage.AccountsClient", "Update")
	}

	req, err := client.UpdatePreparer(resourceGroupName, accountName, parameters)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "Update", nil, "Failure preparing request")
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "Update", resp, "Failure sending request")
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Update", resp, "Failure responding to request")
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client AccountsClient) UpdatePreparer(resourceGroupName string, accountName string, parameters AccountUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client AccountsClient) UpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client AccountsClient) UpdateResponder(resp *http.Response) (result Account, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
