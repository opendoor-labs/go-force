package force

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/opendoor-labs/go-force/sobjects"
)

// Interface all standard and custom objects must implement. Needed for uri generation.
// Note: There was once an ExternalAPIName() interface method which
//   was intended to return the name of the key corresponding to an
//   external id that could be used to operate on sobjects.  Since a
//   SObject can have more than one external id property, this didn't make
//   sense so it was removed.
type SObject interface {
	APIName() string
}

// Response received from force.com API after insert of an sobject.
type SObjectResponse struct {
	Id     string    `force:"id,omitempty"`
	Errors ApiErrors `force:"error,omitempty"` //TODO: Not sure if ApiErrors is the right object
}

func (forceAPI *ForceApi) DescribeSObjects() (map[string]*SObjectMetaData, error) {
	if err := forceAPI.getApiSObjects(); err != nil {
		return nil, err
	}

	return forceAPI.apiSObjects, nil
}

func (forceApi *ForceApi) DescribeSObject(in SObject) (resp *SObjectDescription, err error) {
	// Check cache
	resp, ok := forceApi.apiSObjectDescriptions[in.APIName()]
	if !ok {
		// Attempt retrieval from api
		sObjectMetaData, ok := forceApi.apiSObjects[in.APIName()]
		if !ok {
			err = fmt.Errorf("Unable to find metadata for object: %v", in.APIName())
			return
		}

		uri := sObjectMetaData.URLs[sObjectDescribeKey]

		resp = &SObjectDescription{}
		_, err = forceApi.Get(uri, nil, resp)
		if err != nil {
			return
		}

		// Create Comma Separated String of All Field Names.
		// Used for SELECT * Queries.
		length := len(resp.Fields)
		if length > 0 {
			var allFields bytes.Buffer
			for index, field := range resp.Fields {
				// Field type location cannot be directly retrieved from SQL Query.
				if field.Type != "location" {
					if index > 0 && index < length {
						allFields.WriteString(", ")
					}
					allFields.WriteString(field.Name)
				}
			}

			resp.AllFields = allFields.String()
		}

		forceApi.apiSObjectDescriptions[in.APIName()] = resp
	}

	return
}

func (forceApi *ForceApi) GetSObject(id string, fields []string, out SObject) (err error) {
	uri := strings.Replace(forceApi.apiSObjects[out.APIName()].URLs[rowTemplateKey], idKey, id, 1)

	params := url.Values{}
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	_, err = forceApi.Get(uri, params, out.(interface{}))

	return
}

func (forceApi *ForceApi) InsertSObject(in SObject) (resp *SObjectResponse, err error) {
	uri := forceApi.apiSObjects[in.APIName()].URLs[sObjectKey]

	resp = &SObjectResponse{}
	err = forceApi.Post(uri, nil, in.(interface{}), resp)

	return
}

func (forceApi *ForceApi) UpdateSObject(id string, in SObject) (err error) {
	uri := strings.Replace(forceApi.apiSObjects[in.APIName()].URLs[rowTemplateKey], idKey, id, 1)

	_, err = forceApi.Patch(uri, nil, in.(interface{}), nil)

	return
}

func (forceApi *ForceApi) DeleteSObject(id string, in SObject) (err error) {
	uri := strings.Replace(forceApi.apiSObjects[in.APIName()].URLs[rowTemplateKey], idKey, id, 1)

	err = forceApi.Delete(uri, nil)

	return
}

func idsFromURIs(uris []string) []string {
	ids := make([]string, len(uris))
	for i, uri := range uris {
		parts := strings.Split(uri, "/")
		ids[i] = parts[len(parts)-1]
	}
	return ids
}

func (forceApi *ForceApi) getSingleSFID(uri string, params url.Values) (string, int, error) {
	sobj := sobjects.BaseSObject{}
	statusCode, err := forceApi.Get(uri, params, &sobj)
	return sobj.Id, statusCode, err
}

func (forceApi *ForceApi) getMultipleSFIDs(uri string, params url.Values) ([]string, int, error) {
	// We don't have access to the json that was returned, so make the
	// same call as was made in getSingleSFID() passing in a slice to
	// unmarshal into.
	uris := []string{}
	statusCode, err := forceApi.Get(uri, params, &uris)
	if err != nil {
		return nil, statusCode, err
	}

	return idsFromURIs(uris), statusCode, nil
}

func (forceApi *ForceApi) GetSFIDsByExternalId(apiName, externalKey, externalId string) ([]string, int, error) {
	uri := fmt.Sprintf("%v/%v/%v", forceApi.apiSObjects[apiName].URLs[sObjectKey], externalKey, externalId)
	params := url.Values{"fields": []string{"Id"}}

	sfid, statusCode, err := forceApi.getSingleSFID(uri, params)
	if err == nil {
		return []string{sfid}, statusCode, nil
	}

	// https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/errorcodes.htm
	// Status code 300 (StatusMultipleChoices) is returned when an external
	// ID exists in more than one record. The response body contains the
	// list of matching records.
	if statusCode == http.StatusMultipleChoices {
		return forceApi.getMultipleSFIDs(uri, params)
	}

	return nil, statusCode, err
}

func (forceApi *ForceApi) GetSObjectByExternalId(externalKey, externalId string, fields []string, out SObject) (statusCode int, err error) {
	uri := fmt.Sprintf("%v/%v/%v", forceApi.apiSObjects[out.APIName()].URLs[sObjectKey],
		externalKey, externalId)

	params := url.Values{}
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	return forceApi.Get(uri, params, out.(interface{}))
}

func (forceApi *ForceApi) UpsertSObjectByExternalId(
	externalKey string, externalId string, in SObject) (responseCode int, resp *SObjectResponse, err error) {

	uri := fmt.Sprintf("%v/%v/%v", forceApi.apiSObjects[in.APIName()].URLs[sObjectKey], externalKey, externalId)

	resp = &SObjectResponse{}
	responseCode, err = forceApi.Patch(uri, nil, in.(interface{}), resp)

	return
}

func (forceApi *ForceApi) DeleteSObjectByExternalId(externalKey, externalId string, in SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", forceApi.apiSObjects[in.APIName()].URLs[sObjectKey],
		externalKey, externalId)

	err = forceApi.Delete(uri, nil)

	return
}
