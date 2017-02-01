package ignition

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sync"

	"github.com/coreos/go-systemd/unit"
	"github.com/coreos/ignition/config/types"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"ignition_config":        resourceConfig(),
			"ignition_disk":          resourceDisk(),
			"ignition_raid":          resourceRaid(),
			"ignition_filesystem":    resourceFilesystem(),
			"ignition_file":          resourceFile(),
			"ignition_systemd_unit":  resourceSystemdUnit(),
			"ignition_networkd_unit": resourceNetworkdUnit(),
			"ignition_user":          resourceUser(),
			"ignition_group":         resourceGroup(),
		},
		ConfigureFunc: func(*schema.ResourceData) (interface{}, error) {
			return &cache{
				disks:         make(map[string]*types.Disk, 0),
				arrays:        make(map[string]*types.Raid, 0),
				filesystems:   make(map[string]*types.Filesystem, 0),
				files:         make(map[string]*types.File, 0),
				systemdUnits:  make(map[string]*types.SystemdUnit, 0),
				networkdUnits: make(map[string]*types.NetworkdUnit, 0),
				users:         make(map[string]*types.User, 0),
				groups:        make(map[string]*types.Group, 0),
			}, nil
		},
	}
}

type cache struct {
	disks         map[string]*types.Disk
	arrays        map[string]*types.Raid
	filesystems   map[string]*types.Filesystem
	files         map[string]*types.File
	systemdUnits  map[string]*types.SystemdUnit
	networkdUnits map[string]*types.NetworkdUnit
	users         map[string]*types.User
	groups        map[string]*types.Group

	sync.Mutex
}

func (c *cache) addDisk(g *types.Disk) string {
	c.Lock()
	defer c.Unlock()

	id := id(g)
	c.disks[id] = g

	return id
}

func (c *cache) addRaid(r *types.Raid) string {
	c.Lock()
	defer c.Unlock()

	id := id(r)
	c.arrays[id] = r

	return id
}

func (c *cache) addFilesystem(f *types.Filesystem) string {
	c.Lock()
	defer c.Unlock()

	id := id(f)
	c.filesystems[id] = f

	return id
}

func (c *cache) addFile(f *types.File) string {
	c.Lock()
	defer c.Unlock()

	id := id(f)
	c.files[id] = f

	return id
}

func (c *cache) addSystemdUnit(u *types.SystemdUnit) string {
	c.Lock()
	defer c.Unlock()

	id := id(u)
	c.systemdUnits[id] = u

	return id
}

func (c *cache) addNetworkdUnit(u *types.NetworkdUnit) string {
	c.Lock()
	defer c.Unlock()

	id := id(u)
	c.networkdUnits[id] = u

	return id
}

func (c *cache) addUser(u *types.User) string {
	c.Lock()
	defer c.Unlock()

	id := id(u)
	c.users[id] = u

	return id
}

func (c *cache) addGroup(g *types.Group) string {
	c.Lock()
	defer c.Unlock()

	id := id(g)
	c.groups[id] = g

	return id
}

func id(input interface{}) string {
	b, _ := json.Marshal(input)
	return hash(string(b))
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}

func castSliceInterface(i []interface{}) []string {
	var o []string
	for _, value := range i {
		o = append(o, value.(string))
	}

	return o
}

func getUInt(d *schema.ResourceData, key string) *uint {
	var uid *uint
	if value, ok := d.GetOk(key); ok {
		u := uint(value.(int))
		uid = &u
	}

	return uid
}

func validateUnit(content string) error {
	r := bytes.NewBuffer([]byte(content))

	u, err := unit.Deserialize(r)
	if len(u) == 0 {
		return fmt.Errorf("invalid or empty unit content")
	}

	if err == nil {
		return nil
	}

	if err == io.EOF {
		return fmt.Errorf("unexpected EOF reading unit content")
	}

	return err
}

func buildURL(raw string) (types.Url, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return types.Url{}, err
	}

	return types.Url(*u), nil
}

func buildHash(raw string) (types.Hash, error) {
	h := types.Hash{}
	err := h.UnmarshalJSON([]byte(fmt.Sprintf("%q", raw)))

	return h, err
}
