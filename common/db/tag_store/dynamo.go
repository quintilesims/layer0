package tag_store

// !!! TODO: current vendored package is actually the forked zpatrick/dynamo
// should re-vendor package if/when the PR is merged: https://github.com/guregu/dynamo/pull/30
// !!!
import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/quintilesims/layer0/common/models"
	"time"
)

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
	tag.TagID = time.Now().UnixNano()
	return d.table.Put(tag).Run()
}

func (d *DynamoTagStore) Delete(tagID int64) error {
	return d.table.Delete("TagID", tagID).Run()
}

func (d *DynamoTagStore) SelectAll() (models.Tags, error) {
	var tags models.Tags
	if err := d.table.Scan().Consistent(true).All(&tags); err != nil {
		return nil, err
	}

	return tags, nil
}

func (d *DynamoTagStore) SelectByQuery(entityType, entityID string) (models.Tags, error) {
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

	var tags models.Tags
	if err := scan.Consistent(true).All(&tags); err != nil {
		return nil, err
	}

	return tags, nil
}
