package fortios

import (
	"fmt"
	log "github.com/sourcegraph-ce/logrus"

	fmgclient "github.com/fortinetdev/forti-sdk-go/fortimanager/sdkcore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceFortimanagerFirewallObjectAddress() *schema.Resource {
	return &schema.Resource{
		Create: createFMGFirewallObjectAddress,
		Read:   readFMGFirewallObjectAddress,
		Update: updateFMGFirewallObjectAddress,
		Delete: deleteFMGFirewallObjectAddress,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ipmask",
				ValidateFunc: validation.StringInSlice([]string{
					"ipmask", "iprange", "fqdn",
				}, false),
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"fqdn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"associated_intf": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "any",
			},
			"subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.0.0.0 0.0.0.0",
			},
			"start_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.0.0.0",
			},
			"end_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "255.255.255.255",
			},
			"allow_routing": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "disable",
				ValidateFunc: validation.StringInSlice([]string{
					"enable", "disable",
				}, false),
			},
			"adom": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "root",
			},
		},
	}
}

func createFMGFirewallObjectAddress(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).ClientFortimanager
	defer c.Trace("createFMGFirewallObjectAddress")()

	i := &fmgclient.JSONFirewallObjectAddress{
		Name:           d.Get("name").(string),
		Type:           d.Get("type").(string),
		Comment:        d.Get("comment").(string),
		Fqdn:           d.Get("fqdn").(string),
		AssociatedIntf: d.Get("associated_intf").(string),
		Subnet:         d.Get("subnet").(string),
		StartIp:        d.Get("start_ip").(string),
		EndIp:          d.Get("end_ip").(string),
		AllowRouting:   d.Get("allow_routing").(string),
	}

	err := c.CreateUpdateFirewallObjectAddress(i, "add", d.Get("adom").(string))
	if err != nil {
		return fmt.Errorf("Error creating Firewall Object Address: %s", err)
	}

	d.SetId(i.Name)

	return readFMGFirewallObjectAddress(d, m)
}

func readFMGFirewallObjectAddress(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).ClientFortimanager
	defer c.Trace("readFMGFirewallObjectAddress")()

	name := d.Id()
	o, err := c.ReadFirewallObjectAddress(d.Get("adom").(string), name)
	if err != nil {
		return fmt.Errorf("Error reading Firewall Object Address: %s", err)
	}

	if o == nil {
		log.Printf("[WARN] resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", o.Name)
	d.Set("comment", o.Comment)
	d.Set("type", o.Type)
	d.Set("associated_intf", o.AssociatedIntf)
	if o.Type == "fqdn" {
		d.Set("fqdn", o.Fqdn)
	} else if o.Type == "ipmask" {
		d.Set("subnet", o.Subnet)
		d.Set("allow_routing", o.AllowRouting)
	} else if o.Type == "iprange" {
		d.Set("start_ip", o.StartIp)
		d.Set("end_id", o.EndIp)
	}

	return nil
}

func updateFMGFirewallObjectAddress(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).ClientFortimanager
	defer c.Trace("updateFMGFirewallObjectAddress")()

	if d.HasChange("name") {
		return fmt.Errorf("the name argument is the key and should not be modified here")
	}

	i := &fmgclient.JSONFirewallObjectAddress{
		Name:           d.Get("name").(string),
		Type:           d.Get("type").(string),
		Comment:        d.Get("comment").(string),
		Fqdn:           d.Get("fqdn").(string),
		AssociatedIntf: d.Get("associated_intf").(string),
		Subnet:         d.Get("subnet").(string),
		StartIp:        d.Get("start_ip").(string),
		EndIp:          d.Get("end_ip").(string),
		AllowRouting:   d.Get("allow_routing").(string),
	}

	err := c.CreateUpdateFirewallObjectAddress(i, "update", d.Get("adom").(string))
	if err != nil {
		return fmt.Errorf("Error updating Firewall Object Address: %s", err)
	}

	return readFMGFirewallObjectAddress(d, m)
}

func deleteFMGFirewallObjectAddress(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).ClientFortimanager
	defer c.Trace("deleteFMGFirewallObjectAddress")()

	name := d.Id()

	err := c.DeleteFirewallObjectAddress(d.Get("adom").(string), name)
	if err != nil {
		return fmt.Errorf("Error deleting Firewall Object Address: %s", err)
	}

	d.SetId("")

	return nil
}
