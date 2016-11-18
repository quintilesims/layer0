package command

import (
	"github.com/urfave/cli"
	"gitlab.imshealth.com/xfra/layer0/cli/entity"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"io/ioutil"
)

type CertificateCommand struct {
	*Command
}

func NewCertificateCommand(command *Command) *CertificateCommand {
	return &CertificateCommand{command}
}

func (d *CertificateCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:  "certificate",
		Usage: "manage layer0 certificates",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new certificate",
				Action:    wrapAction(d.Command, d.Create),
				ArgsUsage: "NAME PUBLIC PRIVATE",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "chain",
						Usage: "The path to your intermediate certificate chain",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a certificate",
				ArgsUsage: "NAME",
				Action:    wrapAction(d.Command, d.Delete),
			},
			{
				Name:      "get",
				Usage:     "describe a certificate",
				Action:    wrapAction(d.Command, d.Get),
				ArgsUsage: "NAME",
			},
			{
				Name:      "list",
				Usage:     "list all certificates",
				Action:    wrapAction(d.Command, d.List),
				ArgsUsage: " ",
			},
		},
	}
}

func (d *CertificateCommand) Create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME", "PUBLIC", "PRIVATE")
	if err != nil {
		return err
	}

	public, err := ioutil.ReadFile(args["PUBLIC"])
	if err != nil {
		return err
	}

	private, err := ioutil.ReadFile(args["PRIVATE"])
	if err != nil {
		return err
	}

	var chain []byte
	if path := c.String("chain"); path != "" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		chain = content
	}

	certificate, err := d.Client.CreateCertificate(args["NAME"], public, private, chain)
	if err != nil {
		return err
	}

	return d.printCertificate(certificate)
}

func (d *CertificateCommand) Delete(c *cli.Context) error {
	return d.delete(c, "certificate", d.Client.DeleteCertificate)
}

func (d *CertificateCommand) Get(c *cli.Context) error {
	return d.get(c, "certificate", func(id string) (entity.Entity, error) {
		certificate, err := d.Client.GetCertificate(id)
		if err != nil {
			return nil, err
		}

		return entity.NewCertificate(certificate), nil
	})
}

func (d *CertificateCommand) List(c *cli.Context) error {
	certificates, err := d.Client.ListCertificates()
	if err != nil {
		return err
	}

	return d.printCertificates(certificates)
}

func (d *CertificateCommand) printCertificate(certificate *models.Certificate) error {
	entity := entity.NewCertificate(certificate)
	return d.Printer.PrintEntity(entity)
}

func (d *CertificateCommand) printCertificates(certificates []*models.Certificate) error {
	entities := []entity.Entity{}
	for _, certificate := range certificates {
		entities = append(entities, entity.NewCertificate(certificate))
	}

	return d.Printer.PrintEntities(entities)
}
