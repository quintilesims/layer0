# Developing the AWS Provider

### Rules of Thumb
* Hide ugliness caused by aws-sdk-go by keeping it as encapsulated as possible
* Helper functions should make life easier for the caller (typically, the master function)
* Arguments/variables should be explicit on what they represent.
For example:
```
Bad:  func createCluster(environment string) error
Good: func createCluster(clusterName string) error
```
* The master function should make resource-name relationships explicit. 
For example:
```
	launchConfigName := fqEnvironmentID
	if err := e.createLaunchConfig(launchConfigName, ...); err != nil { ... }
```

* Use `common.go` sparingly
* Helper functions should not call other  functions
* Follow existing patterns/conventions of other entities. Consistency is key!
* Delete operations must be idempotent
* Use a Layer0 error type when assigning a meaningful code to it (anything that isn't an unexpected error): 
```
errors.Newf(errors.InvalidRequest, "Example operation is not supported in Layer0")
```

## Skeletons
The following sections provide skeletons for each entity action. 

### Create 
```
package main

func (e *EntityProvider) Create(req models.CreateEntityRequest) (*models.Entity, error) {
	entityID := createEntityID(req.EntityName)
	fqEntityID := addLayer0Prefix(entityID)

	// if the request has default arguments or complex objects, fetch/convert
	// those objects from the request using  functions
	arg1 := getArg1(req.Field1)
	arg2 := getArg2(req.Field2)

	// setup up args for resource A
	resourceAName := fqEntityID
	if err := e.createResourceA(resourceAName, args...); err != nil {
		return nil, err
	}

	// setup args for resource B. The output is dependent on later calls
	resourceBName := fqEntityID
	resourceB, err := e.createResourceB(resourceBName, args...)
	if err != nil {
		return nil, err
	}

	// create the tags for the entity
	if err := e.createTags(entityID, resourceB.FieldA); err != nil {
		return nil, err
	}

	// return Read() after a create
	return e.Read(entityID)
}

func (e *EntityProvider) createResourceA(args) error {
	input := &aws.Input{}
	input.FieldA(args)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := e.AWS.ECS.Create(input); err != nil {
		return err
	}

	return nil
}

func (e *EntityProvider) createResourceB(args) (*aws.ResourceB, error) {
	input := &aws.Input{}
	input.FieldA(args)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := e.AWS.ECS.Create(input)
	if err != nil {
		return nil, err
	}

	return output.ResourceB, nil
}

func (e *EntityProvider) createTags(entityID, args ...string) error {
	tags := []models.Tag{
		{
			EntityID:   entityID,
			EntityType: "entity_type",
			Key:        "name",
			Value:      arg,
		},
		{
			EntityID:   entityID,
			EntityType: "entity_type",
			Key:        "other",
			Value:      arg,
		},
	}

	for _, tag := range tags {
		if err := e.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}
```

### Delete
```
package main

func (e *EntityProvider) Delete(entityID string) error {
	fqEntityID := addLayer0Prefix(entityID)

	if err := e.deleteResourceA(args...); err != nil {
		return err
	}

	if err := e.deleteResourceB(args...); err != nil {
		return err
	}

	// use the common helper function to delete tags
	if err := deleteEntityTags(e.TagStore, "entity_type", entityID); err != nil {
		return err
	}

	return nil
}

func (e *EntityProvider) DeleteResourceA(args) error {
	input := &aws.Input{}
	input.FieldA(args)

	if err := input.Validate(); err != nil {
		return err
	}

	// catch idempotent errors
	if err := e.AWS.ECS.Delete(input); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "ResourceNotFoundException" {
			return nil
		}
	}

	return nil
}

func (e *EntityProvider) DeleteResourceB(args) error {
	input := &aws.Input{}
	input.FieldA(args)

	if err := input.Validate(); err != nil {
		return err
	}

	// catch idempotent errors
	if err := e.AWS.ECS.Delete(input); err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "ResourceNotFoundException" {
			return nil
		}
	}

	return nil
}
```

### Read
```
package main

func (e *EntityProvider) Read(entityID string) (*models.Entity, error) {
	fqEntityID := addLayer0Prefix(entityID)

	// lookup dependent identifiers if necessary
	environmentID, err := lookupEnvironmentID(e.TagStore, "entity_type", entityID)
	if err != nil {
		return nil, err
	}
	fqEnvironmentID := addLayer0Prefix(environmentID)

	// setup args for readResourceA
	resourceAName := fqEnvironmentID
	resourceA, err := e.readResourceA(resourceAName, fqEntityID)
	if err != nil {
		return nil, err
	}

	resourceB, err := e.readResourceB(args...)
	if err != nil {
		return nil, err
	}

	model, err := e.makeNewEntityModel(args...)
	if err != nil {
                return nil, err
        }

	// make sure to use un-qualifed entity ids in the model
	// safely de-reference pointers using aws-sdk-go helper functions
	model.FieldA = aws.IntValue(resourceA.Field)
	modelFieldB = aws.StringValue(resourceB.Field)

	return model, nil
}

func (e *EntityProvider) readResourceA(args) (*aws.ResourceA, error) {
	input := &aws.Input{}
	input.FieldA(args)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := e.AWS.ECS.Describe(input)
	if err != nil {
		return nil, err
	}

	return output.ResourceA, nil
}

// if the helper function's AWS call returns a slice of objects,
// and we expect `len(slice) == 0` or `len(slice) == 1`, standardize
// on a pattern like this:
func (e *EntityProvider) readResourceB(args) (*aws.ResourceB, error) {
	input := &aws.Input{}
	input.FieldA(args)

	if err := input.Validate(); err != nil {
		return nil, err
	}

    // assuming `Describe()` returns `(ECSOutput, err)`
    // and `ECSOutput.ResourceBs` is `[]*ResourceB`
	output, err := e.AWS.ECS.Describe(input)
	if err != nil {
        if err, ok := err.(awserr.Error); ok && err.Code() == "<Entity>NotFoundException" {
            return nil, errors.Newf(errors.<Entity>DoesNotExist, "<Entity> '%s' does not exist", <entity>ID)
        }

        return nil, err
    }

    if len(output.ResourceBs) == 0 {
        return nil, errors.Newf(errors.<Entity>DoesNotExist, "<Entity> '%s' does not exist", <entity>ID)
    }

	return output.ResourceBs[0], nil
}

func (e *EntityProvider) makeEntityModel(entityID string) (*models.Entity, error) {
	model := &models.Entity{
		EntityID: entityID,
	}

	tags, err := e.TagStore.SelectByTypeAndID("entity_type", entityID)
	if err != nil {
		return nil, err
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.EntityName = tag.Value
	}

	return model, nil
}
```

### List
```
package main

func (e *EntityProvider) List() ([]models.EntitySummary, error) {
	resourceAIDs, err := e.listResourceAIDs()
	if err != nil {
		return nil, err
	}

        entityIDs := make([]string, len(clusterNames))
        for i, clusterName := range clusterNames {
                entityID := delLayer0Prefix(e.Config.Instance(), clusterName)
                entityIDs[i] = entityID
        }

        return e.makeEntitySummaryModels(entityIDs)
}

func (e *EntityProvider) listResourceIDs(args) ([]string, error) {
	resourceIDs = []string{}
	fn := func(output *ecs.Output, lastPage bool) bool {
		for _, resourceID := range output.ResourceIDs {
			// safely de-reference pointers
			resourceIDs = append(resourceIDs, aws.StringValue(resourceID))
			return !lastPage
		}
	}

	if err := e.AWS.ECS.ListResourcePages(&input{}, fn); err != nil {
		return nil, err
	}

	return resourceIDs, nil
}

func (e *EntityProvider) makeEntitySummaryModels(entityIDs []string) ([]models.EntitySummary, error) {
        tags, err := e.TagStore.SelectByType("entity")
        if err != nil {
                return nil, err
        }

        models := make([]models.EntitySummary, len(entityIDs))
        for i, entityID := range entityIDs {
                models[i].EntityID = entityID

                if tag, ok := tags.WithID(entityID).WithKey("name").First(); ok {
                        models[i].EntityName = tag.Value
                }
        }

        return models, nil
}
```

### Update
```
package main

func (e *EntityProvider) Update(req models.UpdateEntityRequest) error {
	entityID := createEntityID(req.EntityName)
	fqEntityID := addLayer0Prefix(entityID)

	// only run functions if the optional params are specified
	if req.FieldA != nil {
		if err := e.updateResourceA(args); err != nil {
			return err
		}
	}

	if req.FieldB != nil {
		if err := e.updateResourceB(args); err != nil {
			return err
		}
	}

	return nil
}

func (e *EntityProvider) updateResourceA(args) error {
	input := &aws.Input{}
	input.FieldA(args)

	if err := input.Validate(); err != nil {
		return err
	}

	output, err := e.AWS.ECS.Update(input)
	if err != nil {
		return err
	}

	return nil
}

func (e *EntityProvider) updateResourceB(args) error {
	input := &aws.Input{}
	input.FieldB(args)

	if err := input.Validate(); err != nil {
		return err
	}

	output, err := e.AWS.ECS.Update(input)
	if err != nil {
		return err
	}

	return nil
}
```
