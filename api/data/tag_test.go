package data

import (
	"flag"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

// Main test entrypoint
func TestMain(m *testing.M) {
	log.SetLevel(log.FatalLevel)
	flag.Parse()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestGetTags_successes(t *testing.T) {
	os.Unsetenv("LAYER0_SQLITE_DB_PATH")
	dataStore, err := NewTagSQLiteDataStore()
	if err != nil {
		t.Fatal(err)
	}
	defer dataStore.Close()

	tag1_123 := makeTag("name", "tag1", "123", "service")
	tag0_123 := makeTag("name", "tag0", "123", "service")
	tag2_456 := makeTag("name", "tag2", "456", "service")
	tag2_789 := makeTag("name", "tag2", "789", "service")
	tage_456 := makeTag("environment_name", "envtag", "456", "service")
	tage_789 := makeTag("environment_name", "envother", "789", "service")
	tag_e123 := makeTag("name", "envtag", "e123", "environment")
	tag_e456 := makeTag("name", "envother", "e456", "environment")
	tag_d123 := makeTag("name", "dep", "d123", "deploy")
	tag_c123 := makeTag("name", "cert1", "c123", "certificate")

	allTags := []models.EntityTag{
		tag1_123, tag0_123, tag2_456, tag2_789, tage_456, tage_789,
		tag_e123, tag_e456, tag_d123, tag_c123,
	}

	for _, tag := range allTags {
		dataStore.Insert(tag)
	}

	allExpected := []models.EntityWithTags{
		makeExpected("123", "service", []models.EntityTag{tag1_123, tag0_123}),
		makeExpected("456", "service", []models.EntityTag{tag2_456, tage_456}),
		makeExpected("789", "service", []models.EntityTag{tag2_789, tage_789}),
		makeExpected("e123", "environment", []models.EntityTag{tag_e123}),
		makeExpected("e456", "environment", []models.EntityTag{tag_e456}),
		makeExpected("d123", "deploy", []models.EntityTag{tag_d123}),
		makeExpected("c123", "certificate", []models.EntityTag{tag_c123}),
	}

	getController := func() (*TagLogicLayer, error) {
		return NewTagLogicLayer(dataStore), nil
	}

	testCases := []struct {
		Name             string
		ExpectedResponse []models.EntityWithTags
		Query            map[string]string
		GetTarget        func() (*TagLogicLayer, error)
	}{
		{
			Name:             "Should list all entities without a request filter",
			ExpectedResponse: allExpected,
			Query:            map[string]string{},
			GetTarget:        getController,
		}, {
			Name:             "Should list services with a type filter",
			ExpectedResponse: allExpected[0:3],
			Query:            map[string]string{"type": "service"},
			GetTarget:        getController,
		}, {
			Name:             "Should list environments with another type filter",
			ExpectedResponse: allExpected[5:7],
			Query:            map[string]string{"type": "environment"},
			GetTarget:        getController,
		}, {
			Name:             "Should filter tags by name",
			ExpectedResponse: allExpected[0:1],
			Query:            map[string]string{"type": "service", "name": "tag0"},
			GetTarget:        getController,
		}, {
			Name:             "Should filter tags by name multiple tags",
			ExpectedResponse: allExpected[1:2],
			Query:            map[string]string{"type": "service", "name": "tag2", "environment_name": "envtag"},
			GetTarget:        getController,
		}, {
			Name:             "Should return empty for no match",
			ExpectedResponse: []models.EntityWithTags{},
			Query:            map[string]string{"type": "service", "name": "tag1", "environment_name": "envtag", "bogusTag": "4321"},
			GetTarget:        getController,
		}, {
			Name: "Should return tags without service type",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("456", "service", []models.EntityTag{tag2_456}),
				makeExpected("789", "service", []models.EntityTag{tag2_789}),
			},
			Query:     map[string]string{"name": "tag2"},
			GetTarget: getController,
		}, {
			Name:             "Should list no entities if none exist",
			ExpectedResponse: []models.EntityWithTags{},
			GetTarget:        defaultTagGetTarget,
		},
	}

	for _, testCase := range testCases {
		log.Infof("Test case %s", testCase.Name)
		target, err := testCase.GetTarget()
		if err != nil {
			t.Errorf("Test %s\n failed test setup with error %v", testCase.Name, err)
		}
		outputResponse, err := target.GetTags(testCase.Query)

		if err != nil {
			t.Errorf("Test %s\n Unexpected error: %v", testCase.Name, err)
		}

		compareResponse(t, testCase.Name, sortResponse(outputResponse), sortResponse(testCase.ExpectedResponse))
	}
}

func TestGetTags_errors(t *testing.T) {
	os.Unsetenv("LAYER0_SQLITE_DB_PATH")
	testCases := []struct {
		Name              string
		ExpectedErrorCode errors.ErrorCode
		Query             map[string]string
		GetTarget         func() (*TagLogicLayer, error)
	}{
		{
			Name:              "Should error from invalid service type",
			ExpectedErrorCode: errors.InvalidEntityType,
			Query:             map[string]string{"type": "bogus"},
			GetTarget:         defaultTagGetTarget,
		}, {
			Name:              "Case Sensitive service type",
			ExpectedErrorCode: errors.InvalidEntityType,
			Query:             map[string]string{"type": "Deploy"},
			GetTarget:         defaultTagGetTarget,
		},
	}

	for _, testCase := range testCases {
		target, err := testCase.GetTarget()
		if err != nil {
			t.Errorf("Test %s\n failed test setup with error %v", testCase.Name, err)
		}

		_, err = target.GetTags(testCase.Query)
		if err == nil {
			t.Errorf("Test %s\n should have returned an error", testCase.Name)
		}

		serverError := err.(*errors.ServerError)
		compareResponse(t, testCase.Name, serverError.Code, testCase.ExpectedErrorCode)
	}
}

func TestGetTags_NamePrefix(t *testing.T) {
	os.Unsetenv("LAYER0_SQLITE_DB_PATH")
	dataStore, err := NewTagSQLiteDataStore()
	if err != nil {
		t.Fatal(err)
	}
	defer dataStore.Close()

	tag1_jj := makeTag("name", "jj1", "123", "service")
	tag1_cs := makeTag("name", "cs1", "234", "service")
	tag1_zp := makeTag("name", "zp1", "345", "service")
	tag2_zp := makeTag("name", "zp2", "456", "service")

	allTags := []models.EntityTag{
		tag1_jj, tag1_cs, tag1_zp, tag2_zp,
	}

	for _, tag := range allTags {
		dataStore.Insert(tag)
	}

	getController := func() *TagLogicLayer {
		return NewTagLogicLayer(dataStore)
	}

	testCases := []struct {
		Name             string
		ExpectedResponse []models.EntityWithTags
		Query            map[string]string
		GetTarget        func() *TagLogicLayer
	}{
		{
			Name: "Should return tag with name_prefix",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("123", "service", []models.EntityTag{tag1_jj}),
			},
			Query:     map[string]string{"name_prefix": "jj"},
			GetTarget: getController,
		}, {
			Name: "Should return multiple tags with name_prefix",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("345", "service", []models.EntityTag{tag1_zp}),
				makeExpected("456", "service", []models.EntityTag{tag2_zp}),
			},
			Query:     map[string]string{"name_prefix": "zp"},
			GetTarget: getController,
		}, {
			Name: "Should allow empty name prefix",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("123", "service", []models.EntityTag{tag1_jj}),
				makeExpected("234", "service", []models.EntityTag{tag1_cs}),
				makeExpected("345", "service", []models.EntityTag{tag1_zp}),
				makeExpected("456", "service", []models.EntityTag{tag2_zp}),
			},
			Query:     map[string]string{"name_prefix": ""},
			GetTarget: getController,
		}, {
			Name:             "Should return empty on mismatched prefix",
			ExpectedResponse: []models.EntityWithTags{},
			Query:            map[string]string{"name_prefix": "zz"},
			GetTarget:        getController,
		}, {
			Name: "Should combine with other filters",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("345", "service", []models.EntityTag{tag1_zp}),
			},
			Query:     map[string]string{"name_prefix": "z", "id": "345"},
			GetTarget: getController,
		},
	}

	for _, testCase := range testCases {
		target := testCase.GetTarget()
		outputResponse, err := target.GetTags(testCase.Query)

		if err != nil {
			t.Errorf("Test %s\n Unexpected error: %v\n", testCase.Name, err)
		}

		compareResponse(t, testCase.Name, sortResponse(outputResponse), sortResponse(testCase.ExpectedResponse))
	}
}

func TestGetTags_IdPrefix(t *testing.T) {
	os.Unsetenv("LAYER0_SQLITE_DB_PATH")
	dataStore, err := NewTagSQLiteDataStore()
	if err != nil {
		t.Fatal(err)
	}
	defer dataStore.Close()

	tag1 := makeTag("name", "", "1234567890", "service")
	tag2 := makeTag("name", "", "234567890a", "environment")
	tag3 := makeTag("name", "", "34567890ab", "service")
	tag4 := makeTag("name", "", "1567890abc", "deploy")

	allTags := []models.EntityTag{
		tag1, tag2, tag3, tag4,
	}

	for _, tag := range allTags {
		dataStore.Insert(tag)
	}

	getController := func() *TagLogicLayer {
		return NewTagLogicLayer(dataStore)
	}

	testCases := []struct {
		Name             string
		ExpectedResponse []models.EntityWithTags
		Query            map[string]string
		GetTarget        func() *TagLogicLayer
	}{
		{
			Name: "Should return tag with id_prefix",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("234567890a", "environment", []models.EntityTag{tag2}),
			},
			Query:     map[string]string{"id_prefix": "234"},
			GetTarget: getController,
		}, {
			Name: "Should allow empty id prefix",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("1234567890", "service", []models.EntityTag{tag1}),
				makeExpected("234567890a", "environment", []models.EntityTag{tag2}),
				makeExpected("34567890ab", "service", []models.EntityTag{tag3}),
				makeExpected("1567890abc", "deploy", []models.EntityTag{tag4}),
			},
			Query:     map[string]string{"id_prefix": ""},
			GetTarget: getController,
		}, {
			Name:             "Should return empty on mismatched prefix",
			ExpectedResponse: []models.EntityWithTags{},
			Query:            map[string]string{"id_prefix": "zz"},
			GetTarget:        getController,
		}, {
			Name: "Should combine id_prefix with other filters",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("1567890abc", "deploy", []models.EntityTag{tag4}),
			},
			Query:     map[string]string{"id_prefix": "1", "type": "deploy"},
			GetTarget: getController,
		},
	}

	for _, testCase := range testCases {
		target := testCase.GetTarget()
		outputResponse, err := target.GetTags(testCase.Query)

		if err != nil {
			t.Errorf("Test %s\n Unexpected error: %v\n", testCase.Name, err)
		}

		compareResponse(t, testCase.Name, sortResponse(outputResponse), sortResponse(testCase.ExpectedResponse))
	}
}

func TestGetTags_NamePrefix_errors(t *testing.T) {
	os.Unsetenv("LAYER0_SQLITE_DB_PATH")
	dataStore, err := NewTagSQLiteDataStore()
	if err != nil {
		t.Fatal(err)
	}
	defer dataStore.Close()

	tag1 := makeTag("name", "plant", "apple", "service")
	tag2 := makeTag("name", "pl_nt", "ap_le", "service")

	allTags := []models.EntityTag{
		tag1, tag2,
	}

	for _, tag := range allTags {
		dataStore.Insert(tag)
	}

	getController := func() *TagLogicLayer {
		return NewTagLogicLayer(dataStore)
	}

	testCases := []struct {
		Name             string
		ExpectedResponse []models.EntityWithTags
		Query            map[string]string
		GetTarget        func() *TagLogicLayer
	}{
		{
			Name:             "Should escape wildcards in id search",
			ExpectedResponse: []models.EntityWithTags{},
			Query:            map[string]string{"id_prefix": "%e"},
			GetTarget:        getController,
		}, {
			Name: "Should escape underscore in id search",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("ap_le", "service", []models.EntityTag{tag2}),
			},
			Query:     map[string]string{"id_prefix": "ap_"},
			GetTarget: getController,
		}, {
			Name:             "Should escape wildcards in name search",
			ExpectedResponse: []models.EntityWithTags{},
			Query:            map[string]string{"name_prefix": "%t"},
			GetTarget:        getController,
		}, {
			Name: "Should escape underscore in name search",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("ap_le", "service", []models.EntityTag{tag2}),
			},
			Query:     map[string]string{"name_prefix": "pl_"},
			GetTarget: getController,
		},
	}

	for _, testCase := range testCases {
		target := testCase.GetTarget()
		outputResponse, err := target.GetTags(testCase.Query)

		if err != nil {
			t.Errorf("Test %s\n Unexpected error: %v\n", testCase.Name, err)
		}

		compareResponse(t, testCase.Name, sortResponse(outputResponse), sortResponse(testCase.ExpectedResponse))
	}
}

func TestGetTags_Fuzz(t *testing.T) {
	os.Unsetenv("LAYER0_SQLITE_DB_PATH")
	dataStore, err := NewTagSQLiteDataStore()
	if err != nil {
		t.Fatal(err)
	}
	defer dataStore.Close()

	tag1 := makeTag("name", "chip", "ch123", "service")
	tag2 := makeTag("name", "dale", "ch456", "service")
	tag3 := makeTag("name", "dale-dale", "da789", "service")

	allTags := []models.EntityTag{
		tag1, tag2, tag3,
	}

	for _, tag := range allTags {
		dataStore.Insert(tag)
	}

	getController := func() *TagLogicLayer {
		return NewTagLogicLayer(dataStore)
	}

	testCases := []struct {
		Name             string
		ExpectedResponse []models.EntityWithTags
		Query            map[string]string
		GetTarget        func() *TagLogicLayer
	}{
		{
			Name:             "Should match nothing",
			ExpectedResponse: []models.EntityWithTags{},
			Query:            map[string]string{"fuzz": "da123"},
			GetTarget:        getController,
		}, {
			Name: "Should match multiple",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("ch456", "service", []models.EntityTag{tag2}),
				makeExpected("da789", "service", []models.EntityTag{tag3}),
			},
			Query:     map[string]string{"fuzz": "dal"},
			GetTarget: getController,
		},
	}

	for _, testCase := range testCases {
		target := testCase.GetTarget()
		outputResponse, err := target.GetTags(testCase.Query)

		if err != nil {
			t.Errorf("Test %s\n Unexpected error: %v\n", testCase.Name, err)
		}

		compareResponse(t, testCase.Name, sortResponse(outputResponse), sortResponse(testCase.ExpectedResponse))
	}
}

func TestGetTags_Version(t *testing.T) {
	os.Unsetenv("LAYER0_SQLITE_DB_PATH")
	dataStore, err := NewTagSQLiteDataStore()
	if err != nil {
		t.Fatal(err)
	}
	defer dataStore.Close()

	tag1 := makeTag("version", "1", "123", "deploy")
	tag2 := makeTag("version", "2", "456", "deploy")
	tag3 := makeTag("version", "3", "789", "deploy")
	tag4 := makeTag("version", "4", "1011", "service")
	tag5 := makeTag("enable", "prod", "456", "deploy")

	allTags := []models.EntityTag{
		tag1, tag2, tag3, tag4, tag5,
	}

	for _, tag := range allTags {
		dataStore.Insert(tag)
	}

	getController := func() *TagLogicLayer {
		return NewTagLogicLayer(dataStore)
	}

	testCases := []struct {
		Name             string
		ExpectedResponse []models.EntityWithTags
		Query            map[string]string
		GetTarget        func() *TagLogicLayer
	}{
		{
			Name: "Version latest tag returns across services",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("1011", "service", []models.EntityTag{tag4}),
			},
			Query:     map[string]string{"version": "latest"},
			GetTarget: getController,
		}, {
			Name: "Type Filtered version",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("789", "deploy", []models.EntityTag{tag3}),
			},
			Query:     map[string]string{"version": "latest", "type": "deploy"},
			GetTarget: getController,
		}, {
			Name: "Version filter with other tags",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("456", "deploy", []models.EntityTag{tag2, tag5}),
			},
			Query:     map[string]string{"version": "latest", "enable": "prod"},
			GetTarget: getController,
		}, {
			Name: "Order of version tags",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("789", "deploy", []models.EntityTag{tag3}),
			},
			Query:     map[string]string{"type": "deploy", "version": "latest"},
			GetTarget: getController,
		}, {
			Name: "Specific non-latest version",
			ExpectedResponse: []models.EntityWithTags{
				makeExpected("123", "deploy", []models.EntityTag{tag1}),
			},
			Query:     map[string]string{"type": "deploy", "version": "1"},
			GetTarget: getController,
		},
	}

	for _, testCase := range testCases {
		target := testCase.GetTarget()
		outputResponse, err := target.GetTags(testCase.Query)

		if err != nil {
			t.Errorf("Test %s\n Unexpected error: %v\n", testCase.Name, err)
		}

		compareResponse(t, testCase.Name, sortResponse(outputResponse), sortResponse(testCase.ExpectedResponse))
	}
}

func TestMakeTag_successes(t *testing.T) {
	os.Unsetenv("LAYER0_SQLITE_DB_PATH")
	testCases := []struct {
		Name      string
		Tag       models.EntityTag
		GetTarget func() (*TagLogicLayer, error)
	}{
		{
			Name:      "Legal name tag",
			Tag:       makeTag("name", "superSimple1", "123", "service"),
			GetTarget: defaultTagGetTarget,
		}, {
			Name:      "Limited special characters",
			Tag:       makeTag("name", "limited_special-char.acters", "123", "service"),
			GetTarget: defaultTagGetTarget,
		},
	}

	for _, testCase := range testCases {
		target, err := testCase.GetTarget()
		if err != nil {
			t.Errorf("Test %s\n Unexpected test setup error: %v\n", testCase.Name, err)
			continue
		}

		err = target.Make(testCase.Tag)

		if err != nil {
			t.Errorf("Test %s\n Unexpected make error: %v\n", testCase.Name, err)
			continue
		}
		// check that our tag was actually created
		tags, err := target.GetTags(map[string]string{})
		if err != nil {
			t.Errorf("Test %s\n Unexpected get error: %v\n", testCase.Name, err)
		}
		expected := []models.EntityWithTags{
			makeExpected(testCase.Tag.EntityID,
				testCase.Tag.EntityType,
				[]models.EntityTag{testCase.Tag}),
		}
		compareResponse(t, testCase.Name, tags, expected)
	}
}

func TestMakeTag_errors(t *testing.T) {
	os.Unsetenv("LAYER0_SQLITE_DB_PATH")
	testCases := []struct {
		Name              string
		ExpectedErrorCode errors.ErrorCode
		Tag               models.EntityTag
		GetTarget         func() (*TagLogicLayer, error)
	}{
		{
			Name:              "Should prevent key name_prefix",
			ExpectedErrorCode: errors.InvalidTagKey,
			Tag:               makeTag("name_prefix", "bogus", "123", "service"),
			GetTarget:         defaultTagGetTarget,
		}, {
			Name:              "Should prevent key id_prefix",
			ExpectedErrorCode: errors.InvalidTagKey,
			Tag:               makeTag("id_prefix", "bogus", "123", "service"),
			GetTarget:         defaultTagGetTarget,
		},
		{
			Name:              "Should prevent key fuzz",
			ExpectedErrorCode: errors.InvalidTagKey,
			Tag:               makeTag("fuzz", "bogus", "123", "service"),
			GetTarget:         defaultTagGetTarget,
		}, {
			Name:              "Should prevent special characters",
			ExpectedErrorCode: errors.InvalidTagValue,
			Tag:               makeTag("custom", "\\|?!@#$%^&*()", "123", "service"),
			GetTarget:         defaultTagGetTarget,
		}, {
			Name:              "Key too long",
			ExpectedErrorCode: errors.InvalidTagKey,
			Tag:               makeTag(makeString("name", 65), "value", "123", "service"),
			GetTarget:         defaultTagGetTarget,
		}, {
			Name:              "Value too long",
			ExpectedErrorCode: errors.InvalidTagValue,
			Tag:               makeTag("name", makeString("value", 65), "123", "service"),
			GetTarget:         defaultTagGetTarget,
		}, {
			Name:              "Version must be an integer",
			ExpectedErrorCode: errors.InvalidTagValue,
			Tag:               makeTag("version", "best", "123", "deploy"),
			GetTarget:         defaultTagGetTarget,
		},
	}

	for _, testCase := range testCases {
		target, err := testCase.GetTarget()
		if err != nil {
			t.Errorf("Test %s\n Unexpected test setup error: %v\n", testCase.Name, err)
			continue
		}

		err = target.Make(testCase.Tag)

		if err == nil {
			t.Errorf("Test %s\n Unexpected lack of error for: %v\n", testCase.Name, testCase.Tag)
			continue
		}
		l0err, ok := err.(*errors.ServerError)
		if !ok {
			t.Errorf("Test %s\n Expected ServerError, observed:  %v", testCase.Name, err)
			continue
		}
		compareResponse(t, testCase.Name, l0err.Code, testCase.ExpectedErrorCode)

		// check that no tag was actually created
		tags, err := target.GetTags(map[string]string{})
		if err != nil {
			t.Errorf("Test %s\n Unexpected error: %v\n", testCase.Name, err)
		}
		compareResponse(t, testCase.Name, tags, []models.EntityWithTags{})
	}
}

func makeTag(tag_key, tag_value, eid, etype string) models.EntityTag {
	return models.EntityTag{
		EntityID:   eid,
		EntityType: etype,
		Key:        tag_key,
		Value:      tag_value,
	}
}

func makeExpected(eid, etype string, tags []models.EntityTag) models.EntityWithTags {
	return models.EntityWithTags{
		EntityID:   eid,
		EntityType: etype,
		Tags:       tags,
	}
}

func makeString(base string, length int) string {
	var result = []string{base}
	for i := len(base); i < length; i++ {
		result = append(result, strconv.Itoa((i+1)%10))
	}
	return strings.Join(result, "")
}

func compareResponse(t *testing.T, name string, result, expected interface{}) {
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed on case '%s':\n  observed %v (%v)\n  expected %v (%v)",
			name, result, reflect.TypeOf(result), expected, reflect.TypeOf(expected))
	}
}

func defaultTagGetTarget() (*TagLogicLayer, error) {
	dataStore, err := NewTagSQLiteDataStore()
	if err != nil {
		return nil, err
	}
	return NewTagLogicLayer(dataStore), nil
}

// Our Response tags come out of the map unsorted
// so we need to sort them for compareResponse to be useful.
// Lots of boiler-plate code follows, see below:
func sortResponse(entityList []models.EntityWithTags) []models.EntityWithTags {
	var s ByEntity = byId
	s.Sort(entityList)

	var s2 ByTag = byTag
	for _, e := range entityList {
		s2.Sort(e.Tags)
	}
	return entityList
}

// Sort Entity list by comparing first ID, and then Type as a fallback.
func byId(p1, p2 *models.EntityWithTags) bool {
	if p1.EntityID == p2.EntityID {
		return p1.EntityType < p2.EntityType
	}

	return p1.EntityID < p2.EntityID
}

type ByEntity func(p1, p2 *models.EntityWithTags) bool

func (by ByEntity) Sort(entityList []models.EntityWithTags) {
	ps := &entitySorter{
		list: entityList,
		by:   by,
	}
	sort.Sort(ps)
}

type entitySorter struct {
	list []models.EntityWithTags
	by   ByEntity
}

func (s *entitySorter) Len() int {
	return len(s.list)
}

func (s *entitySorter) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[j], s.list[i]
}

func (s *entitySorter) Less(i, j int) bool {
	return s.by(&s.list[i], &s.list[j])
}

// Sort EntityTag by comparing first ID, and then Type as a fallback.
func byTag(p1, p2 *models.EntityTag) bool {
	if p1.Key == p2.Key {
		return p1.Value < p2.Value
	}

	return p1.Key < p2.Key
}

type ByTag func(p1, p2 *models.EntityTag) bool

func (by ByTag) Sort(entityList []models.EntityTag) {
	ps := &tagSorter{
		list: entityList,
		by:   by,
	}
	sort.Sort(ps)
}

type tagSorter struct {
	list []models.EntityTag
	by   ByTag
}

func (s *tagSorter) Len() int {
	return len(s.list)
}

func (s *tagSorter) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[j], s.list[i]
}

func (s *tagSorter) Less(i, j int) bool {
	return s.by(&s.list[i], &s.list[j])
}
