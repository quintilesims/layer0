package tag_store

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/quintilesims/layer0/common/models"
	"time"
)

type DynamoTagStore struct {
	table          dynamo.Table
	consistentRead bool
}

func NewDynamoTagStore(session *session.Session, table string) *DynamoTagStore {
	db := dynamo.New(session)

	return &DynamoTagStore{
		table:          db.Table(table),
		consistentRead: false,
	}
}

func (d *DynamoTagStore) setConsistentRead(v bool) {
	d.consistentRead = v
}

func (d *DynamoTagStore) Init() error {
	d.setConsistentRead(true)
	return nil
}

func (d *DynamoTagStore) Clear() error {
	tags, err := d.SelectAll()
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := d.Delete(tag.TagID); err != nil {
			return err
		}
	}

	return nil
}

func (d *DynamoTagStore) Insert(tag *models.Tag) error {
	// ensure we perform a consistent read after a write
	defer d.setConsistentRead(true)

	tag.TagID = time.Now().UnixNano()
	return d.table.Put(tag).Run()
}

func (d *DynamoTagStore) Delete(tagID int64) error {
	// ensure we perform a consistent read after a delete
	defer d.setConsistentRead(true)
	return d.table.Delete("TagID", tagID).Run()
}

func (d *DynamoTagStore) SelectAll() (models.Tags, error) {
	defer d.setConsistentRead(false)

	tags := models.Tags{}
	if err := d.table.Scan().Consistent(d.consistentRead).All(&tags); err != nil {
		return nil, err
	}

	return tags, nil
}

func (d *DynamoTagStore) SelectByQuery(entityType, entityID string) (models.Tags, error) {
	defer d.setConsistentRead(false)

	scan := d.table.Scan()

	switch {
	case entityType != "" && entityID == "":
		scan.Filter("'EntityType' = ?", entityType)
	case entityType == "" && entityID != "":
		scan.Filter("'EntityID' = ?", entityID)
	case entityType != "" && entityID != "":
		scan.Filter("'EntityType' = ? AND $ = ?", entityType, "EntityID", entityID)
	default:
		return d.SelectAll()
	}

	tags := models.Tags{}
	if err := scan.Consistent(d.consistentRead).All(&tags); err != nil {
		return nil, err
	}

	return tags, nil
}
