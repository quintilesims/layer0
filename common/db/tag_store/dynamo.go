package tag_store

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/quintilesims/layer0/common/models"
)

// todo: this assume entity_id, type, and key are unique
type DynamoTagSchema struct {
	EntityType string
	EntityID   string
	Tags       map[string]string
}

func (s DynamoTagSchema) ToTags() models.Tags {
	tags := models.Tags{}
	for k, v := range s.Tags {
		tag := &models.Tag{
			EntityType: s.EntityType,
			EntityID:   s.EntityID,
			Key:        k,
			Value:      v,
		}

		tags = append(tags, tag)
	}

	return tags
}

type DynamoTagStore struct {
	table dynamo.Table
}

func NewDynamoTagStore(session *session.Session, table string) *DynamoTagStore {
	db := dynamo.New(session)

	return &DynamoTagStore{
		table: db.Table(table),
	}
}

func (d *DynamoTagStore) Init() error {
	return nil
}

func (d *DynamoTagStore) Clear() error {
	return nil
}

func (d *DynamoTagStore) Insert(tag models.Tag) error {
	schema := DynamoTagSchema{
		EntityType: tag.EntityType,
		EntityID:   tag.EntityID,
		Tags:       map[string]string{tag.Key: tag.Value},
	}

	if err := d.table.Put(schema).If("attribute_not_exists(EntityType)").Run(); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "ConditionalCheckFailedException" {
			return d.insertKey(tag)
		}

		return err
	}

	return nil
}

func (d *DynamoTagStore) insertKey(tag models.Tag) error {
	schema, err := d.selectByQuery(tag.EntityType, tag.EntityID)
	if err != nil {
		return err
	}

	schema.Tags[tag.Key] = tag.Value
	return d.table.Update("EntityType", tag.EntityType).
		Range("EntityID", tag.EntityID).
		Set("Tags", schema.Tags).
		Run()
}

// todo: don't allow empty args
func (d *DynamoTagStore) SelectByQuery(entityType, entityID string) (models.Tags, error) {
	schema, err := d.selectByQuery(entityType, entityID)
	if err != nil {
		return nil, err
	}

	return schema.ToTags(), nil
}

func (d *DynamoTagStore) selectByQuery(entityType, entityID string) (*DynamoTagSchema, error) {
	if entityType == "" {
		return nil, fmt.Errorf("EntityType is required")
	}

	if entityID == "" {
		return nil, fmt.Errorf("EntityID is required")
	}

	var schema *DynamoTagSchema
	if err := d.table.Get("EntityType", entityType).
		Range("EntityID", dynamo.Equal, entityID).
		Consistent(true).
		One(&schema); err != nil {
		return nil, err
	}

	if schema.Tags == nil {
		schema.Tags = map[string]string{}
	}

	return schema, nil
}

func (d *DynamoTagStore) SelectByType(entityType string) (models.Tags, error) {
	schemas, err := d.selectByType(entityType)
	if err != nil {
		return nil, err
	}

	tags := models.Tags{}
	for _, schema := range schemas {
		tags = append(tags, schema.ToTags()...)
	}

	return tags, nil
}

func (d *DynamoTagStore) selectByType(entityType string) ([]*DynamoTagSchema, error) {
	var schemas []*DynamoTagSchema

	if err := d.table.Get("EntityType", entityType).
		//Range("EntityID", dynamo.Equal, entityID).
		Consistent(true).
		All(&schemas); err != nil {
		return nil, err
	}

	return schemas, nil
}

// todo: we should delete by all tags, or by key
func (d *DynamoTagStore) Delete(tagID int64) error {
	return nil
}
