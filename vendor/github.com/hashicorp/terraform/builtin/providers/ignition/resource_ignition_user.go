package ignition

import (
	"reflect"

	"github.com/coreos/ignition/config/types"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Exists: resourceUserExists,
		Read:   resourceUserRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password_hash": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ssh_authorized_keys": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"uid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"gecos": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"home_dir": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"no_create_home": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"primary_group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"no_user_group": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"no_log_init": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"shell": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildUser(d, meta.(*cache))
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceUserExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildUser(d, meta.(*cache))
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildUser(d *schema.ResourceData, c *cache) (string, error) {
	uc := types.UserCreate{
		Uid:          getUInt(d, "uid"),
		GECOS:        d.Get("gecos").(string),
		Homedir:      d.Get("home_dir").(string),
		NoCreateHome: d.Get("no_create_home").(bool),
		PrimaryGroup: d.Get("primary_group").(string),
		Groups:       castSliceInterface(d.Get("groups").([]interface{})),
		NoUserGroup:  d.Get("no_user_group").(bool),
		NoLogInit:    d.Get("no_log_init").(bool),
		Shell:        d.Get("shell").(string),
	}

	puc := &uc
	if reflect.DeepEqual(uc, types.UserCreate{}) { // check if the struct is empty
		puc = nil
	}

	user := types.User{
		Name:              d.Get("name").(string),
		PasswordHash:      d.Get("password_hash").(string),
		SSHAuthorizedKeys: castSliceInterface(d.Get("ssh_authorized_keys").([]interface{})),
		Create:            puc,
	}

	return c.addUser(&user), nil
}
