package fortios

import (
	"fmt"
	log "github.com/sourcegraph-ce/logrus"

	fmgclient "github.com/fortinetdev/forti-sdk-go/fortimanager/sdkcore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceFortimanagerSystemDNS() *schema.Resource {
	return &schema.Resource{
		Create: setFMGSystemDNS,
		Read:   readFMGSystemDNS,
		Update: setFMGSystemDNS,
		Delete: deleteFMGSystemDNS,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"primary": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"secondary": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func setFMGSystemDNS(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).ClientFortimanager
	defer c.Trace("setFMGSystemDNS")()

	i := &fmgclient.JSONSystemDNS{
		Primary:   d.Get("primary").(string),
		Secondary: d.Get("secondary").(string),
	}

	err := c.SetSystemDNS(i)
	if err != nil {
		return fmt.Errorf("Error setting System DNS : %s", err)
	}

	d.SetId("fortimanager-sys-dns")

	return readFMGSystemDNS(d, m)
}

func readFMGSystemDNS(d *schema.ResourceData, m interface{}) error {
	c := m.(*FortiClient).ClientFortimanager
	defer c.Trace("readFMGSystemDNS")()

	o, err := c.ReadSystemDNS()
	if err != nil {
		return fmt.Errorf("Error reading System DNS : %s", err)
	}

	if o == nil {
		log.Printf("[WARN] resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if d.Get("primary") != "" {
		d.Set("primary", o.Primary)
	}
	if d.Get("secondary") != "" {
		d.Set("secondary", o.Secondary)
	}

	return nil
}

func deleteFMGSystemDNS(d *schema.ResourceData, m interface{}) error {
	return nil
}
