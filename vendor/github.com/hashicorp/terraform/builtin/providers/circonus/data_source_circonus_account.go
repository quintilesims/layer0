package circonus

import (
	"fmt"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/circonus-gometrics/api/config"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	accountAddress1Attr      = "address1"
	accountAddress2Attr      = "address2"
	accountCCEmailAttr       = "cc_email"
	accountCityAttr          = "city"
	accountContactGroupsAttr = "contact_groups"
	accountCountryAttr       = "country"
	accountCurrentAttr       = "current"
	accountDescriptionAttr   = "description"
	accountEmailAttr         = "email"
	accountIDAttr            = "id"
	accountInvitesAttr       = "invites"
	accountLimitAttr         = "limit"
	accountNameAttr          = "name"
	accountOwnerAttr         = "owner"
	accountRoleAttr          = "role"
	accountStateProvAttr     = "state"
	accountTimezoneAttr      = "timezone"
	accountTypeAttr          = "type"
	accountUIBaseURLAttr     = "ui_base_url"
	accountUsageAttr         = "usage"
	accountUsedAttr          = "used"
	accountUserIDAttr        = "id"
	accountUsersAttr         = "users"
)

var accountDescription = map[schemaAttr]string{
	accountContactGroupsAttr: "Contact Groups in this account",
	accountInvitesAttr:       "Outstanding invites attached to the account",
	accountUsageAttr:         "Account's usage limits",
	accountUsersAttr:         "Users attached to this account",
}

func dataSourceCirconusAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCirconusAccountRead,

		Schema: map[string]*schema.Schema{
			accountAddress1Attr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountAddress1Attr],
			},
			accountAddress2Attr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountAddress2Attr],
			},
			accountCCEmailAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountCCEmailAttr],
			},
			accountIDAttr: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{accountCurrentAttr},
				ValidateFunc: validateFuncs(
					validateRegexp(accountIDAttr, config.AccountCIDRegex),
				),
				Description: accountDescription[accountIDAttr],
			},
			accountCityAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountCityAttr],
			},
			accountContactGroupsAttr: {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: accountDescription[accountContactGroupsAttr],
			},
			accountCountryAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountCountryAttr],
			},
			accountCurrentAttr: {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{accountIDAttr},
				Description:   accountDescription[accountCurrentAttr],
			},
			accountDescriptionAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountDescriptionAttr],
			},
			accountInvitesAttr: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: accountDescription[accountInvitesAttr],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						accountEmailAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: accountDescription[accountEmailAttr],
						},
						accountRoleAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: accountDescription[accountRoleAttr],
						},
					},
				},
			},
			accountNameAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountNameAttr],
			},
			accountOwnerAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountOwnerAttr],
			},
			accountStateProvAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountStateProvAttr],
			},
			accountTimezoneAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountTimezoneAttr],
			},
			accountUIBaseURLAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountUIBaseURLAttr],
			},
			accountUsageAttr: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: accountDescription[accountUsageAttr],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						accountLimitAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: accountDescription[accountLimitAttr],
						},
						accountTypeAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: accountDescription[accountTypeAttr],
						},
						accountUsedAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: accountDescription[accountUsedAttr],
						},
					},
				},
			},
			accountUsersAttr: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: accountDescription[accountUsersAttr],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						accountUserIDAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: accountDescription[accountUserIDAttr],
						},
						accountRoleAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: accountDescription[accountRoleAttr],
						},
					},
				},
			},
		},
	}
}

func dataSourceCirconusAccountRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	var cid string

	var a *api.Account
	var err error
	if v, ok := d.GetOk(accountIDAttr); ok {
		cid = v.(string)
	}

	if v, ok := d.GetOk(accountCurrentAttr); ok {
		if v.(bool) {
			cid = ""
		}
	}

	a, err = c.client.FetchAccount(api.CIDType(&cid))
	if err != nil {
		return err
	}

	invitesList := make([]interface{}, 0, len(a.Invites))
	for i := range a.Invites {
		invitesList = append(invitesList, map[string]interface{}{
			accountEmailAttr: a.Invites[i].Email,
			accountRoleAttr:  a.Invites[i].Role,
		})
	}

	usageList := make([]interface{}, 0, len(a.Usage))
	for i := range a.Usage {
		usageList = append(usageList, map[string]interface{}{
			accountLimitAttr: a.Usage[i].Limit,
			accountTypeAttr:  a.Usage[i].Type,
			accountUsedAttr:  a.Usage[i].Used,
		})
	}

	usersList := make([]interface{}, 0, len(a.Users))
	for i := range a.Users {
		usersList = append(usersList, map[string]interface{}{
			accountUserIDAttr: a.Users[i].UserCID,
			accountRoleAttr:   a.Users[i].Role,
		})
	}

	d.SetId(a.CID)

	d.Set(accountAddress1Attr, a.Address1)
	d.Set(accountAddress2Attr, a.Address2)
	d.Set(accountCCEmailAttr, a.CCEmail)
	d.Set(accountIDAttr, a.CID)
	d.Set(accountCityAttr, a.City)
	d.Set(accountContactGroupsAttr, a.ContactGroups)
	d.Set(accountCountryAttr, a.Country)
	d.Set(accountDescriptionAttr, a.Description)

	if err := d.Set(accountInvitesAttr, invitesList); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to store account %q attribute: {{err}}", accountInvitesAttr), err)
	}

	d.Set(accountNameAttr, a.Name)
	d.Set(accountOwnerAttr, a.OwnerCID)
	d.Set(accountStateProvAttr, a.StateProv)
	d.Set(accountTimezoneAttr, a.Timezone)
	d.Set(accountUIBaseURLAttr, a.UIBaseURL)

	if err := d.Set(accountUsageAttr, usageList); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to store account %q attribute: {{err}}", accountUsageAttr), err)
	}

	if err := d.Set(accountUsersAttr, usersList); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to store account %q attribute: {{err}}", accountUsersAttr), err)
	}

	return nil
}
