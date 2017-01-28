package ignition

import (
	"github.com/coreos/ignition/config/types"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Delete: resourceUserDelete,
		Exists: resourceUserExists,
		Read:   resourceUserRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password_hash": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ssh_authorized_keys": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"uid": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"gecos": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"home_dir": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"no_create_home": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"primary_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"groups": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"no_user_group": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"no_log_init": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"shell": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	id, err := buildUser(d, meta.(*cache))
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourceUserExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildUser(d, meta.(*cache))
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func buildUser(d *schema.ResourceData, c *cache) (string, error) {
	return c.addUser(&types.User{
		Name:              d.Get("name").(string),
		PasswordHash:      d.Get("password_hash").(string),
		SSHAuthorizedKeys: castSliceInterface(d.Get("ssh_authorized_keys").([]interface{})),
		Create: &types.UserCreate{
			Uid:          getUInt(d, "uid"),
			GECOS:        d.Get("gecos").(string),
			Homedir:      d.Get("home_dir").(string),
			NoCreateHome: d.Get("no_create_home").(bool),
			PrimaryGroup: d.Get("primary_group").(string),
			Groups:       castSliceInterface(d.Get("groups").([]interface{})),
			NoUserGroup:  d.Get("no_user_group").(bool),
			NoLogInit:    d.Get("no_log_init").(bool),
			Shell:        d.Get("shell").(string),
		},
	}), nil
}
