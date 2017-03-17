package google

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
)

// resourceGoogleProject returns a *schema.Resource that allows a customer
// to declare a Google Cloud Project resource.
//
// This example shows a project with a policy declared in config:
//
// resource "google_project" "my-project" {
//    project = "a-project-id"
//    policy = "${data.google_iam_policy.admin.policy}"
// }
func resourceGoogleProject() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceGoogleProjectCreate,
		Read:   resourceGoogleProjectRead,
		Update: resourceGoogleProjectUpdate,
		Delete: resourceGoogleProjectDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState: resourceGoogleProjectMigrateState,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "The id field has unexpected behaviour and probably doesn't do what you expect. See https://www.terraform.io/docs/providers/google/r/google_project.html#id-field for more information. Please use project_id instead; future versions of Terraform will remove the id field.",
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// This suppresses the diff if project_id is not set
					if new == "" {
						return true
					}
					return false
				},
			},
			"skip_delete": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_data": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Deprecated:       "Use the 'google_project_iam_policy' resource to define policies for a Google Project",
				DiffSuppressFunc: jsonPolicyDiffSuppress,
			},
			"policy_etag": &schema.Schema{
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use the the 'google_project_iam_policy' resource to define policies for a Google Project",
			},
			"number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_account": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGoogleProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var pid string
	var err error
	pid = d.Get("project_id").(string)
	if pid == "" {
		pid, err = getProject(d, config)
		if err != nil {
			return fmt.Errorf("Error getting project ID: %v", err)
		}
		if pid == "" {
			return fmt.Errorf("'project_id' must be set in the config")
		}
	}

	// we need to check if name and org_id are set, and throw an error if they aren't
	// we can't just set these as required on the object, however, as that would break
	// all configs that used previous iterations of the resource.
	// TODO(paddy): remove this for 0.9 and set these attributes as required.
	name, org_id := d.Get("name").(string), d.Get("org_id").(string)
	if name == "" {
		return fmt.Errorf("`name` must be set in the config if you're creating a project.")
	}
	if org_id == "" {
		return fmt.Errorf("`org_id` must be set in the config if you're creating a project.")
	}

	log.Printf("[DEBUG]: Creating new project %q", pid)
	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      d.Get("name").(string),
		Parent: &cloudresourcemanager.ResourceId{
			Id:   d.Get("org_id").(string),
			Type: "organization",
		},
	}

	op, err := config.clientResourceManager.Projects.Create(project).Do()
	if err != nil {
		return fmt.Errorf("Error creating project %s (%s): %s.", project.ProjectId, project.Name, err)
	}

	d.SetId(pid)

	// Wait for the operation to complete
	waitErr := resourceManagerOperationWait(config, op, "project to create")
	if waitErr != nil {
		return waitErr
	}

	// Apply the IAM policy if it is set
	if pString, ok := d.GetOk("policy_data"); ok {
		// The policy string is just a marshaled cloudresourcemanager.Policy.
		// Unmarshal it to a struct.
		var policy cloudresourcemanager.Policy
		if err := json.Unmarshal([]byte(pString.(string)), &policy); err != nil {
			return err
		}
		log.Printf("[DEBUG] Got policy from config: %#v", policy.Bindings)

		// Retrieve existing IAM policy from project. This will be merged
		// with the policy defined here.
		p, err := getProjectIamPolicy(pid, config)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Got existing bindings from project: %#v", p.Bindings)

		// Merge the existing policy bindings with those defined in this manifest.
		p.Bindings = mergeBindings(append(p.Bindings, policy.Bindings...))

		// Apply the merged policy
		log.Printf("[DEBUG] Setting new policy for project: %#v", p)
		_, err = config.clientResourceManager.Projects.SetIamPolicy(pid,
			&cloudresourcemanager.SetIamPolicyRequest{Policy: p}).Do()

		if err != nil {
			return fmt.Errorf("Error applying IAM policy for project %q: %s", pid, err)
		}
	}

	// Set the billing account
	if v, ok := d.GetOk("billing_account"); ok {
		name := v.(string)
		ba := cloudbilling.ProjectBillingInfo{
			BillingAccountName: "billingAccounts/" + name,
		}
		_, err = config.clientBilling.Projects.UpdateBillingInfo(prefixedProject(pid), &ba).Do()
		if err != nil {
			d.Set("billing_account", "")
			if _err, ok := err.(*googleapi.Error); ok {
				return fmt.Errorf("Error setting billing account %q for project %q: %v", name, prefixedProject(pid), _err)
			}
			return fmt.Errorf("Error setting billing account %q for project %q: %v", name, prefixedProject(pid), err)
		}
	}

	return resourceGoogleProjectRead(d, meta)
}

func resourceGoogleProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid := d.Id()

	// Read the project
	p, err := config.clientResourceManager.Projects.Get(pid).Do()
	if err != nil {
		if v, ok := err.(*googleapi.Error); ok && v.Code == http.StatusNotFound {
			return fmt.Errorf("Project %q does not exist.", pid)
		}
		return fmt.Errorf("Error checking project %q: %s", pid, err)
	}

	d.Set("project_id", pid)
	d.Set("number", strconv.FormatInt(int64(p.ProjectNumber), 10))
	d.Set("name", p.Name)

	if p.Parent != nil {
		d.Set("org_id", p.Parent.Id)
	}

	// Read the billing account
	ba, err := config.clientBilling.Projects.GetBillingInfo(prefixedProject(pid)).Do()
	if err != nil {
		return fmt.Errorf("Error reading billing account for project %q: %v", prefixedProject(pid), err)
	}
	if ba.BillingAccountName != "" {
		// BillingAccountName is contains the resource name of the billing account
		// associated with the project, if any. For example,
		// `billingAccounts/012345-567890-ABCDEF`. We care about the ID and not
		// the `billingAccounts/` prefix, so we need to remove that. If the
		// prefix ever changes, we'll validate to make sure it's something we
		// recognize.
		_ba := strings.TrimPrefix(ba.BillingAccountName, "billingAccounts/")
		if ba.BillingAccountName == _ba {
			return fmt.Errorf("Error parsing billing account for project %q. Expected value to begin with 'billingAccounts/' but got %s", prefixedProject(pid), ba.BillingAccountName)
		}
		d.Set("billing_account", _ba)
	}
	return nil
}

func prefixedProject(pid string) string {
	return "projects/" + pid
}
func resourceGoogleProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid := d.Id()

	// Read the project
	// we need the project even though refresh has already been called
	// because the API doesn't support patch, so we need the actual object
	p, err := config.clientResourceManager.Projects.Get(pid).Do()
	if err != nil {
		if v, ok := err.(*googleapi.Error); ok && v.Code == http.StatusNotFound {
			return fmt.Errorf("Project %q does not exist.", pid)
		}
		return fmt.Errorf("Error checking project %q: %s", pid, err)
	}

	// Project name has changed
	if ok := d.HasChange("name"); ok {
		p.Name = d.Get("name").(string)
		// Do update on project
		p, err = config.clientResourceManager.Projects.Update(p.ProjectId, p).Do()
		if err != nil {
			return fmt.Errorf("Error updating project %q: %s", p.Name, err)
		}
	}

	// Billing account has changed
	if ok := d.HasChange("billing_account"); ok {
		name := d.Get("billing_account").(string)
		ba := cloudbilling.ProjectBillingInfo{
			BillingAccountName: "billingAccounts/" + name,
		}
		_, err = config.clientBilling.Projects.UpdateBillingInfo(prefixedProject(pid), &ba).Do()
		if err != nil {
			d.Set("billing_account", "")
			if _err, ok := err.(*googleapi.Error); ok {
				return fmt.Errorf("Error updating billing account %q for project %q: %v", name, prefixedProject(pid), _err)
			}
			return fmt.Errorf("Error updating billing account %q for project %q: %v", name, prefixedProject(pid), err)
		}
	}
	return updateProjectIamPolicy(d, config, pid)
}

func resourceGoogleProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	// Only delete projects if skip_delete isn't set
	if !d.Get("skip_delete").(bool) {
		pid := d.Id()
		_, err := config.clientResourceManager.Projects.Delete(pid).Do()
		if err != nil {
			return fmt.Errorf("Error deleting project %q: %s", pid, err)
		}
	}
	d.SetId("")
	return nil
}

func updateProjectIamPolicy(d *schema.ResourceData, config *Config, pid string) error {
	// Policy has changed
	if ok := d.HasChange("policy_data"); ok {
		// The policy string is just a marshaled cloudresourcemanager.Policy.
		// Unmarshal it to a struct that contains the old and new policies
		oldP, newP := d.GetChange("policy_data")
		oldPString := oldP.(string)
		newPString := newP.(string)

		// JSON Unmarshaling would fail
		if oldPString == "" {
			oldPString = "{}"
		}
		if newPString == "" {
			newPString = "{}"
		}

		log.Printf("[DEBUG]: Old policy: %q\nNew policy: %q", oldPString, newPString)

		var oldPolicy, newPolicy cloudresourcemanager.Policy
		if err := json.Unmarshal([]byte(newPString), &newPolicy); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(oldPString), &oldPolicy); err != nil {
			return err
		}

		// Find any Roles and Members that were removed (i.e., those that are present
		// in the old but absent in the new
		oldMap := rolesToMembersMap(oldPolicy.Bindings)
		newMap := rolesToMembersMap(newPolicy.Bindings)
		deleted := make(map[string]map[string]bool)

		// Get each role and its associated members in the old state
		for role, members := range oldMap {
			// Initialize map for role
			if _, ok := deleted[role]; !ok {
				deleted[role] = make(map[string]bool)
			}
			// The role exists in the new state
			if _, ok := newMap[role]; ok {
				// Check each memeber
				for member, _ := range members {
					// Member does not exist in new state, so it was deleted
					if _, ok = newMap[role][member]; !ok {
						deleted[role][member] = true
					}
				}
			} else {
				// This indicates an entire role was deleted. Mark all members
				// for delete.
				for member, _ := range members {
					deleted[role][member] = true
				}
			}
		}
		log.Printf("[DEBUG] Roles and Members to be deleted: %#v", deleted)

		// Retrieve existing IAM policy from project. This will be merged
		// with the policy in the current state
		// TODO(evanbrown): Add an 'authoritative' flag that allows policy
		// in manifest to overwrite existing policy.
		p, err := getProjectIamPolicy(pid, config)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Got existing bindings from project: %#v", p.Bindings)

		// Merge existing policy with policy in the current state
		log.Printf("[DEBUG] Merging new bindings from project: %#v", newPolicy.Bindings)
		mergedBindings := mergeBindings(append(p.Bindings, newPolicy.Bindings...))

		// Remove any roles and members that were explicitly deleted
		mergedBindingsMap := rolesToMembersMap(mergedBindings)
		for role, members := range deleted {
			for member, _ := range members {
				delete(mergedBindingsMap[role], member)
			}
		}

		p.Bindings = rolesToMembersBinding(mergedBindingsMap)
		dump, _ := json.MarshalIndent(p.Bindings, " ", "  ")
		log.Printf("[DEBUG] Setting new policy for project: %#v:\n%s", p, string(dump))

		_, err = config.clientResourceManager.Projects.SetIamPolicy(pid,
			&cloudresourcemanager.SetIamPolicyRequest{Policy: p}).Do()

		if err != nil {
			return fmt.Errorf("Error applying IAM policy for project %q: %s", pid, err)
		}
	}
	return nil
}
