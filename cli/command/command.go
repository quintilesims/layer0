package command

import (
	"github.com/urfave/cli"
	"gitlab.imshealth.com/xfra/layer0/cli/client"
	"gitlab.imshealth.com/xfra/layer0/cli/entity"
	"gitlab.imshealth.com/xfra/layer0/cli/printer"
)

type CommandGroup interface {
	SetPrinter(printer.Printer)
	GetCommand() cli.Command
}

type Command struct {
	Client   client.Client
	Printer  printer.Printer
	Resolver Resolver
}

func (cm *Command) SetPrinter(printer printer.Printer) {
	cm.Printer = printer
}

func (cm *Command) resolveSingleID(entityType, target string) (string, error) {
	ids, err := cm.Resolver.Resolve(entityType, target)
	if err != nil {
		return "", err
	}

	return assertSingleID(entityType, target, ids)
}

func (cm *Command) handleError(c *cli.Context, err error) {
	if _, ok := err.(*UsageError); ok {
		handleUsageError(c, err)
	}

	text := err.Error()
	if suggestion, hasSuggestion := errorSuggestion(err); hasSuggestion {
		text = suggestion
	}

	code := int64(1)
	if serverErr, ok := err.(*client.ServerError); ok {
		code = serverErr.ErrorCode
	}

	cm.Printer.Fatalf(code, text)
}

// deleteWithJob will fetch the NAME arg and use it to resolve the entity id of the specified type
// the deleteEntity function should wrap a Client Delete<Entity> call and return the jobID
// if the 'wait' flag is specified, the Client.WaitForJob() function will be ran
func (cm *Command) deleteWithJob(c *cli.Context, entityType string, deleteEntity func(string) (string, error)) error {
	return cm.delete(c, entityType, func(id string) error {
		jobID, err := deleteEntity(id)
		if err != nil {
			return err
		}

		if !c.Bool("wait") {
			cm.Printer.Printf("This operation is running as a job. Run `l0 job get %s` to see progress\n", jobID)
			return nil
		}

		timeout, err := getTimeout(c)
		if err != nil {
			return err
		}

		cm.Printer.StartSpinner("Deleting")
		if err := cm.Client.WaitForJob(jobID, timeout); err != nil {
			return err
		}

		return nil
	})
}

// delete will fetch the NAME arg and use it to resolve the entity id of the specified type
// the deleteEntity function should wrap a Client Delete<Entity> call
func (cm *Command) delete(c *cli.Context, entityType string, deleteEntity func(string) error) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	id, err := cm.resolveSingleID(entityType, args["NAME"])
	if err != nil {
		return err
	}

	if err := deleteEntity(id); err != nil {
		return err
	}

	return nil
}

// get will fetch the NAME arg and use it to resolve the entity id of the specified type
// the getEntity function should wrap a Client Get<Entity> call and convert the returned model
// into an entity.Entity
func (cm *Command) get(c *cli.Context, entityType string, getEntity func(string) (entity.Entity, error)) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	ids, err := cm.Resolver.Resolve(entityType, args["NAME"])
	if err != nil {
		return err
	}

	entities := []entity.Entity{}
	for _, id := range ids {
		entity, err := getEntity(id)
		if err != nil {
			return err
		}

		entities = append(entities, entity)
	}

	return cm.Printer.PrintEntities(entities)
}
