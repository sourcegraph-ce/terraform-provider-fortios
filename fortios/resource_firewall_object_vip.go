package fortios

import (
	"fmt"
	log "github.com/sourcegraph-ce/logrus"

	"github.com/fortinetdev/forti-sdk-go/fortios/sdkcore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceFirewallObjectVip() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallObjectVipCreate,
		Read:   resourceFirewallObjectVipRead,
		Update: resourceFirewallObjectVipUpdate,
		Delete: resourceFirewallObjectVipDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Created by Terraform Provider for FortiOS",
			},
			"extip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mappedip": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"extintf": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"portforward": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"extport": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mappedport": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceFirewallObjectVipCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).Client

	if c == nil {
		return fmt.Errorf("FortiOS connection did not initialize successfully!")
	}

	c.Retries = 1

	//Get Params from d
	name := d.Get("name").(string)
	comment := d.Get("comment").(string)
	extip := d.Get("extip").(string)
	mappedip := d.Get("mappedip").([]interface{})
	extintf := d.Get("extintf").(string)
	portforward := d.Get("portforward").(string)
	protocol := d.Get("protocol").(string)
	extport := d.Get("extport").(string)
	mappedport := d.Get("mappedport").(string)

	var mappedips []forticlient.VIPMultValue

	for _, v := range mappedip {
		if v == nil {
			return fmt.Errorf("null value")
		}
		mappedips = append(mappedips,
			forticlient.VIPMultValue{
				Range: v.(string),
			})
	}

	//Build input data by sdk
	i := &forticlient.JSONFirewallObjectVip{
		Name:        name,
		Comment:     comment,
		Extip:       extip,
		Mappedip:    mappedips,
		Extintf:     extintf,
		Portforward: portforward,
		Protocol:    protocol,
		Extport:     extport,
		Mappedport:  mappedport,
	}

	//Call process by sdk
	o, err := c.CreateFirewallObjectVip(i)
	if err != nil {
		return fmt.Errorf("Error creating Firewall Object Vip: %s", err)
	}

	//Set index for d
	d.SetId(o.Mkey)

	return resourceFirewallObjectVipRead(d, m)
}

func resourceFirewallObjectVipUpdate(d *schema.ResourceData, m interface{}) error {
	mkey := d.Id()

	c := m.(*FortiClient).Client

	if c == nil {
		return fmt.Errorf("FortiOS connection did not initialize successfully!")
	}

	c.Retries = 1

	//Get Params from d
	name := d.Get("name").(string)
	comment := d.Get("comment").(string)
	extip := d.Get("extip").(string)
	mappedip := d.Get("mappedip").([]interface{})
	extintf := d.Get("extintf").(string)
	portforward := d.Get("portforward").(string)
	protocol := d.Get("protocol").(string)
	extport := d.Get("extport").(string)
	mappedport := d.Get("mappedport").(string)

	var mappedips []forticlient.VIPMultValue

	for _, v := range mappedip {
		if v == nil {
			return fmt.Errorf("null value")
		}
		mappedips = append(mappedips,
			forticlient.VIPMultValue{
				Range: v.(string),
			})
	}

	//Build input data by sdk
	i := &forticlient.JSONFirewallObjectVip{
		Name:        name,
		Comment:     comment,
		Extip:       extip,
		Mappedip:    mappedips,
		Extintf:     extintf,
		Portforward: portforward,
		Protocol:    protocol,
		Extport:     extport,
		Mappedport:  mappedport,
	}

	//Call process by sdk
	_, err := c.UpdateFirewallObjectVip(i, mkey)
	if err != nil {
		return fmt.Errorf("Error updating Firewall Object Vip: %s", err)
	}

	return resourceFirewallObjectVipRead(d, m)
}

func resourceFirewallObjectVipDelete(d *schema.ResourceData, m interface{}) error {
	mkey := d.Id()

	c := m.(*FortiClient).Client

	if c == nil {
		return fmt.Errorf("FortiOS connection did not initialize successfully!")
	}

	c.Retries = 1

	//Call process by sdk
	err := c.DeleteFirewallObjectVip(mkey)
	if err != nil {
		return fmt.Errorf("Error deleting Firewall Object Vip: %s", err)
	}

	//Set index for d
	d.SetId("")

	return nil
}

func resourceFirewallObjectVipRead(d *schema.ResourceData, m interface{}) error {
	mkey := d.Get("name").(string)

	if mkey == "" {
		mkey = d.Id()
	}

	c := m.(*FortiClient).Client

	if c == nil {
		return fmt.Errorf("FortiOS connection did not initialize successfully!")
	}

	c.Retries = 1

	//Call process by sdk
	o, err := c.ReadFirewallObjectVip(mkey)
	if err != nil {
		return fmt.Errorf("Error reading Firewall Object Vip: %s", err)
	}

	if o == nil {
		log.Printf("[WARN] resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	//Refresh property
	d.SetId(o.Name)
	d.Set("name", o.Name)
	d.Set("comment", o.Comment)
	d.Set("extip", o.Extip)

	vs := make([]string, 0, len(o.Mappedip))
	for _, v := range o.Mappedip {
		c := v.Range
		vs = append(vs, c)
	}

	if err := d.Set("mappedip", vs); err != nil {
		log.Printf("[WARN] Error setting Firewall Object Vip for (%s): %s", d.Id(), err)
	}
	d.Set("extintf", o.Extintf)
	d.Set("portforward", o.Portforward)
	d.Set("protocol", o.Protocol)
	d.Set("extport", o.Extport)
	d.Set("mappedport", o.Mappedport)

	return nil
}
