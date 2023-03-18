package fortios

import (
	"fmt"
	log "github.com/sourcegraph-ce/logrus"

	fmgclient "github.com/fortinetdev/forti-sdk-go/fortimanager/sdkcore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceFortimanagerFirewallObjectIppool() *schema.Resource {
	return &schema.Resource{
		Create: createFMGFirewallObjectIppool,
		Read:   readFMGFirewallObjectIppool,
		Update: updateFMGFirewallObjectIppool,
		Delete: deleteFMGFirewallObjectIppool,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"startip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"endip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "overload",
				ValidateFunc: validation.StringInSlice([]string{
					"overload", "one-to-one",
				}, false),
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"arp_reply": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "enable",
				ValidateFunc: validation.StringInSlice([]string{
					"disable", "enable",
				}, false),
			},
			"arp_intf": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"associated_intf": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"adom": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "root",
			},
		},
	}
}

func createFMGFirewallObjectIppool(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).ClientFortimanager
	defer c.Trace("createFMGFirewallObjectIppool")()

	i := &fmgclient.JSONFirewallObjectIppool{
		Name:           d.Get("name").(string),
		Comment:        d.Get("comment").(string),
		Type:           d.Get("type").(string),
		ArpIntf:        d.Get("arp_intf").(string),
		ArpReply:       d.Get("arp_reply").(string),
		AssociatedIntf: d.Get("associated_intf").(string),
		StartIp:        d.Get("startip").(string),
		EndIp:          d.Get("endip").(string),
	}

	err := c.CreateUpdateFirewallObjectIppool(i, "add", d.Get("adom").(string))
	if err != nil {
		return fmt.Errorf("Error creating Firewall Object Ippool: %s", err)
	}

	d.SetId(i.Name)

	return readFMGFirewallObjectIppool(d, m)
}

func readFMGFirewallObjectIppool(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).ClientFortimanager
	defer c.Trace("readFMGFirewallObjectIppool")()

	name := d.Id()
	o, err := c.ReadFirewallObjectIppool(d.Get("adom").(string), name)
	if err != nil {
		return fmt.Errorf("Error reading Firewall Object Ippool: %s", err)
	}

	if o == nil {
		log.Printf("[WARN] resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", o.Name)
	d.Set("comment", o.Comment)
	d.Set("type", o.Type)
	d.Set("arp-intf", o.ArpIntf)
	d.Set("arp-reply", o.ArpReply)
	d.Set("associated-intf", o.AssociatedIntf)
	d.Set("startip", o.StartIp)
	d.Set("endip", o.EndIp)

	return nil
}

func updateFMGFirewallObjectIppool(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).ClientFortimanager
	defer c.Trace("updateFMGFirewallObjectIppool")()

	if d.HasChange("name") {
		return fmt.Errorf("the name argument is the key and should not be modified here")
	}

	i := &fmgclient.JSONFirewallObjectIppool{
		Name:           d.Get("name").(string),
		Comment:        d.Get("comment").(string),
		Type:           d.Get("type").(string),
		ArpIntf:        d.Get("arp_intf").(string),
		ArpReply:       d.Get("arp_reply").(string),
		AssociatedIntf: d.Get("associated_intf").(string),
		StartIp:        d.Get("startip").(string),
		EndIp:          d.Get("endip").(string),
	}

	err := c.CreateUpdateFirewallObjectIppool(i, "update", d.Get("adom").(string))
	if err != nil {
		return fmt.Errorf("Error updating Firewall Object Ippool: %s", err)
	}

	return readFMGFirewallObjectIppool(d, m)
}

func deleteFMGFirewallObjectIppool(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).ClientFortimanager
	defer c.Trace("deleteFMGFirewallObjectIppool")()

	name := d.Id()

	err := c.DeleteFirewallObjectIppool(d.Get("adom").(string), name)
	if err != nil {
		return fmt.Errorf("Error deleting Firewall Object Ippool: %s", err)
	}

	d.SetId("")

	return nil
}
