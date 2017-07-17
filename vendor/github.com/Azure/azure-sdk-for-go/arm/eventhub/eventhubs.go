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
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
)

// EventHubsClient is the azure Event Hubs client
type EventHubsClient struct {
	ManagementClient
}

// NewEventHubsClient creates an instance of the EventHubsClient client.
func NewEventHubsClient(subscriptionID string) EventHubsClient {
	return NewEventHubsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewEventHubsClientWithBaseURI creates an instance of the EventHubsClient
// client.
func NewEventHubsClientWithBaseURI(baseURI string, subscriptionID string) EventHubsClient {
	return EventHubsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates or updates a new Event Hub as a nested resource
// within a namespace.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. eventHubName is the Event Hub name. parameters is
// parameters supplied to create an Event Hub resource.
func (client EventHubsClient) CreateOrUpdate(resourceGroupName string, namespaceName string, eventHubName string, parameters CreateOrUpdateParameters) (result ResourceType, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.Location", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "eventhub.EventHubsClient", "CreateOrUpdate")
	}

	req, err := client.CreateOrUpdatePreparer(resourceGroupName, namespaceName, eventHubName, parameters)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "CreateOrUpdate", nil, "Failure preparing request")
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "CreateOrUpdate", resp, "Failure sending request")
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client EventHubsClient) CreateOrUpdatePreparer(resourceGroupName string, namespaceName string, eventHubName string, parameters CreateOrUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"eventHubName":      autorest.Encode("path", eventHubName),
		"namespaceName":     autorest.Encode("path", namespaceName),
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
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client EventHubsClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client EventHubsClient) CreateOrUpdateResponder(resp *http.Response) (result ResourceType, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// CreateOrUpdateAuthorizationRule creates or updates an authorization rule
// for the specified Event Hub.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. eventHubName is the Event Hub name. authorizationRuleName
// is the authorization rule name. parameters is the shared access
// authorization rule.
func (client EventHubsClient) CreateOrUpdateAuthorizationRule(resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, parameters SharedAccessAuthorizationRuleCreateOrUpdateParameters) (result SharedAccessAuthorizationRuleResource, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.SharedAccessAuthorizationRuleProperties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "parameters.SharedAccessAuthorizationRuleProperties.Rights", Name: validation.Null, Rule: true, Chain: nil}}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "eventhub.EventHubsClient", "CreateOrUpdateAuthorizationRule")
	}

	req, err := client.CreateOrUpdateAuthorizationRulePreparer(resourceGroupName, namespaceName, eventHubName, authorizationRuleName, parameters)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "CreateOrUpdateAuthorizationRule", nil, "Failure preparing request")
	}

	resp, err := client.CreateOrUpdateAuthorizationRuleSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "CreateOrUpdateAuthorizationRule", resp, "Failure sending request")
	}

	result, err = client.CreateOrUpdateAuthorizationRuleResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "CreateOrUpdateAuthorizationRule", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdateAuthorizationRulePreparer prepares the CreateOrUpdateAuthorizationRule request.
func (client EventHubsClient) CreateOrUpdateAuthorizationRulePreparer(resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, parameters SharedAccessAuthorizationRuleCreateOrUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"authorizationRuleName": autorest.Encode("path", authorizationRuleName),
		"eventHubName":          autorest.Encode("path", eventHubName),
		"namespaceName":         autorest.Encode("path", namespaceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules/{authorizationRuleName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// CreateOrUpdateAuthorizationRuleSender sends the CreateOrUpdateAuthorizationRule request. The method will close the
// http.Response Body if it receives an error.
func (client EventHubsClient) CreateOrUpdateAuthorizationRuleSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CreateOrUpdateAuthorizationRuleResponder handles the response to the CreateOrUpdateAuthorizationRule request. The method always
// closes the http.Response Body.
func (client EventHubsClient) CreateOrUpdateAuthorizationRuleResponder(resp *http.Response) (result SharedAccessAuthorizationRuleResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes an Event Hub from the specified namespace and resource group.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. eventHubName is the name of the Event Hub to delete.
func (client EventHubsClient) Delete(resourceGroupName string, namespaceName string, eventHubName string) (result autorest.Response, err error) {
	req, err := client.DeletePreparer(resourceGroupName, namespaceName, eventHubName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "Delete", nil, "Failure preparing request")
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "Delete", resp, "Failure sending request")
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client EventHubsClient) DeletePreparer(resourceGroupName string, namespaceName string, eventHubName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"eventHubName":      autorest.Encode("path", eventHubName),
		"namespaceName":     autorest.Encode("path", namespaceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client EventHubsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client EventHubsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// DeleteAuthorizationRule deletes an Event Hubs authorization rule.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. eventHubName is the Event Hub name. authorizationRuleName
// is the authorization rule name.
func (client EventHubsClient) DeleteAuthorizationRule(resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string) (result autorest.Response, err error) {
	req, err := client.DeleteAuthorizationRulePreparer(resourceGroupName, namespaceName, eventHubName, authorizationRuleName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "DeleteAuthorizationRule", nil, "Failure preparing request")
	}

	resp, err := client.DeleteAuthorizationRuleSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "DeleteAuthorizationRule", resp, "Failure sending request")
	}

	result, err = client.DeleteAuthorizationRuleResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "DeleteAuthorizationRule", resp, "Failure responding to request")
	}

	return
}

// DeleteAuthorizationRulePreparer prepares the DeleteAuthorizationRule request.
func (client EventHubsClient) DeleteAuthorizationRulePreparer(resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"authorizationRuleName": autorest.Encode("path", authorizationRuleName),
		"eventHubName":          autorest.Encode("path", eventHubName),
		"namespaceName":         autorest.Encode("path", namespaceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules/{authorizationRuleName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// DeleteAuthorizationRuleSender sends the DeleteAuthorizationRule request. The method will close the
// http.Response Body if it receives an error.
func (client EventHubsClient) DeleteAuthorizationRuleSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// DeleteAuthorizationRuleResponder handles the response to the DeleteAuthorizationRule request. The method always
// closes the http.Response Body.
func (client EventHubsClient) DeleteAuthorizationRuleResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets an Event Hubs description for the specified Event Hub.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. eventHubName is the Event Hub name.
func (client EventHubsClient) Get(resourceGroupName string, namespaceName string, eventHubName string) (result ResourceType, err error) {
	req, err := client.GetPreparer(resourceGroupName, namespaceName, eventHubName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "Get", nil, "Failure preparing request")
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "Get", resp, "Failure sending request")
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client EventHubsClient) GetPreparer(resourceGroupName string, namespaceName string, eventHubName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"eventHubName":      autorest.Encode("path", eventHubName),
		"namespaceName":     autorest.Encode("path", namespaceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client EventHubsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client EventHubsClient) GetResponder(resp *http.Response) (result ResourceType, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetAuthorizationRule gets an authorization rule for an Event Hub by rule
// name.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. eventHubName is the Event Hub name. authorizationRuleName
// is the authorization rule name.
func (client EventHubsClient) GetAuthorizationRule(resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string) (result SharedAccessAuthorizationRuleResource, err error) {
	req, err := client.GetAuthorizationRulePreparer(resourceGroupName, namespaceName, eventHubName, authorizationRuleName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "GetAuthorizationRule", nil, "Failure preparing request")
	}

	resp, err := client.GetAuthorizationRuleSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "GetAuthorizationRule", resp, "Failure sending request")
	}

	result, err = client.GetAuthorizationRuleResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "GetAuthorizationRule", resp, "Failure responding to request")
	}

	return
}

// GetAuthorizationRulePreparer prepares the GetAuthorizationRule request.
func (client EventHubsClient) GetAuthorizationRulePreparer(resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"authorizationRuleName": autorest.Encode("path", authorizationRuleName),
		"eventHubName":          autorest.Encode("path", eventHubName),
		"namespaceName":         autorest.Encode("path", namespaceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules/{authorizationRuleName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetAuthorizationRuleSender sends the GetAuthorizationRule request. The method will close the
// http.Response Body if it receives an error.
func (client EventHubsClient) GetAuthorizationRuleSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetAuthorizationRuleResponder handles the response to the GetAuthorizationRule request. The method always
// closes the http.Response Body.
func (client EventHubsClient) GetAuthorizationRuleResponder(resp *http.Response) (result SharedAccessAuthorizationRuleResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListAll gets all the Event Hubs in a namespace.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name.
func (client EventHubsClient) ListAll(resourceGroupName string, namespaceName string) (result ListResult, err error) {
	req, err := client.ListAllPreparer(resourceGroupName, namespaceName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAll", nil, "Failure preparing request")
	}

	resp, err := client.ListAllSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAll", resp, "Failure sending request")
	}

	result, err = client.ListAllResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAll", resp, "Failure responding to request")
	}

	return
}

// ListAllPreparer prepares the ListAll request.
func (client EventHubsClient) ListAllPreparer(resourceGroupName string, namespaceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"namespaceName":     autorest.Encode("path", namespaceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListAllSender sends the ListAll request. The method will close the
// http.Response Body if it receives an error.
func (client EventHubsClient) ListAllSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListAllResponder handles the response to the ListAll request. The method always
// closes the http.Response Body.
func (client EventHubsClient) ListAllResponder(resp *http.Response) (result ListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListAllNextResults retrieves the next set of results, if any.
func (client EventHubsClient) ListAllNextResults(lastResults ListResult) (result ListResult, err error) {
	req, err := lastResults.ListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAll", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListAllSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAll", resp, "Failure sending next results request")
	}

	result, err = client.ListAllResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAll", resp, "Failure responding to next results request")
	}

	return
}

// ListAuthorizationRules gets the authorization rules for an Event Hub.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name eventHubName is the Event Hub name.
func (client EventHubsClient) ListAuthorizationRules(resourceGroupName string, namespaceName string, eventHubName string) (result SharedAccessAuthorizationRuleListResult, err error) {
	req, err := client.ListAuthorizationRulesPreparer(resourceGroupName, namespaceName, eventHubName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAuthorizationRules", nil, "Failure preparing request")
	}

	resp, err := client.ListAuthorizationRulesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAuthorizationRules", resp, "Failure sending request")
	}

	result, err = client.ListAuthorizationRulesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAuthorizationRules", resp, "Failure responding to request")
	}

	return
}

// ListAuthorizationRulesPreparer prepares the ListAuthorizationRules request.
func (client EventHubsClient) ListAuthorizationRulesPreparer(resourceGroupName string, namespaceName string, eventHubName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"eventHubName":      autorest.Encode("path", eventHubName),
		"namespaceName":     autorest.Encode("path", namespaceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListAuthorizationRulesSender sends the ListAuthorizationRules request. The method will close the
// http.Response Body if it receives an error.
func (client EventHubsClient) ListAuthorizationRulesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListAuthorizationRulesResponder handles the response to the ListAuthorizationRules request. The method always
// closes the http.Response Body.
func (client EventHubsClient) ListAuthorizationRulesResponder(resp *http.Response) (result SharedAccessAuthorizationRuleListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListAuthorizationRulesNextResults retrieves the next set of results, if any.
func (client EventHubsClient) ListAuthorizationRulesNextResults(lastResults SharedAccessAuthorizationRuleListResult) (result SharedAccessAuthorizationRuleListResult, err error) {
	req, err := lastResults.SharedAccessAuthorizationRuleListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAuthorizationRules", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListAuthorizationRulesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAuthorizationRules", resp, "Failure sending next results request")
	}

	result, err = client.ListAuthorizationRulesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListAuthorizationRules", resp, "Failure responding to next results request")
	}

	return
}

// ListKeys gets the ACS and SAS connection strings for the Event Hub.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. eventHubName is the Event Hub name. authorizationRuleName
// is the connection string of the namespace for the specified authorization
// rule.
func (client EventHubsClient) ListKeys(resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string) (result ResourceListKeys, err error) {
	req, err := client.ListKeysPreparer(resourceGroupName, namespaceName, eventHubName, authorizationRuleName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListKeys", nil, "Failure preparing request")
	}

	resp, err := client.ListKeysSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListKeys", resp, "Failure sending request")
	}

	result, err = client.ListKeysResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "ListKeys", resp, "Failure responding to request")
	}

	return
}

// ListKeysPreparer prepares the ListKeys request.
func (client EventHubsClient) ListKeysPreparer(resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"authorizationRuleName": autorest.Encode("path", authorizationRuleName),
		"eventHubName":          autorest.Encode("path", eventHubName),
		"namespaceName":         autorest.Encode("path", namespaceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules/{authorizationRuleName}/ListKeys", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListKeysSender sends the ListKeys request. The method will close the
// http.Response Body if it receives an error.
func (client EventHubsClient) ListKeysSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListKeysResponder handles the response to the ListKeys request. The method always
// closes the http.Response Body.
func (client EventHubsClient) ListKeysResponder(resp *http.Response) (result ResourceListKeys, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// RegenerateKeys regenerates the ACS and SAS connection strings for the Event
// Hub.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. eventHubName is the Event Hub name. authorizationRuleName
// is the connection string of the Event Hub for the specified authorization
// rule. parameters is parameters supplied to regenerate the authorization
// rule.
func (client EventHubsClient) RegenerateKeys(resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, parameters RegenerateKeysParameters) (result ResourceListKeys, err error) {
	req, err := client.RegenerateKeysPreparer(resourceGroupName, namespaceName, eventHubName, authorizationRuleName, parameters)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "RegenerateKeys", nil, "Failure preparing request")
	}

	resp, err := client.RegenerateKeysSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "RegenerateKeys", resp, "Failure sending request")
	}

	result, err = client.RegenerateKeysResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhub.EventHubsClient", "RegenerateKeys", resp, "Failure responding to request")
	}

	return
}

// RegenerateKeysPreparer prepares the RegenerateKeys request.
func (client EventHubsClient) RegenerateKeysPreparer(resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, parameters RegenerateKeysParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"authorizationRuleName": autorest.Encode("path", authorizationRuleName),
		"eventHubName":          autorest.Encode("path", eventHubName),
		"namespaceName":         autorest.Encode("path", namespaceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules/{authorizationRuleName}/regenerateKeys", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// RegenerateKeysSender sends the RegenerateKeys request. The method will close the
// http.Response Body if it receives an error.
func (client EventHubsClient) RegenerateKeysSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// RegenerateKeysResponder handles the response to the RegenerateKeys request. The method always
// closes the http.Response Body.
func (client EventHubsClient) RegenerateKeysResponder(resp *http.Response) (result ResourceListKeys, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
