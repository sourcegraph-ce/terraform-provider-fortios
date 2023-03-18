package fortios

import (
	"fmt"
	log "github.com/sourcegraph-ce/logrus"

	forticlient "github.com/fortinetdev/forti-sdk-go/fortios/sdkcore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNetworkingInterfacePort() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingInterfacePortCreate,
		Read:   resourceNetworkingInterfacePortRead,
		Update: resourceNetworkingInterfacePortUpdate,
		Delete: resourceNetworkingInterfacePortDelete,

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
				Required: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_identification": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tcp_mss": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"speed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mtu_override": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"allowaccess": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns_server_override": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"defaultgw": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"distance": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Created by Terraform Provider for FortiOS",
			},
			"interface": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vdom": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vlanid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingInterfacePortCreate(d *schema.ResourceData, m interface{}) error {
	// type, name, vdom, ip, interface is Required
	// if type is "physical", then call update(PUT)
	// else call create(POST)
	// 	if type is "vlan", then check  vlanid, vdom, interface
	// 	if type is "loopback", check vdom

	typef := d.Get("type").(string)

	if typef == "physical" {
		//Get Params from d
		portname := d.Get("name").(string)

		d.SetId(portname)

		//Call process by sdk
		err := resourceNetworkingInterfacePortUpdate(d, m)
		if err != nil {
			return fmt.Errorf("Error creating Networking Interface Port: %s", err)
		}
	} else {
		c := m.(*FortiClient).Client

		if c == nil {
			return fmt.Errorf("FortiOS connection did not initialize successfully!")
		}

		c.Retries = 1

		//Get Params from d
		ip := d.Get("ip").(string)
		alias := d.Get("alias").(string)
		status := d.Get("status").(string)
		deviceIdentification := d.Get("device_identification").(string)
		tcpMss := d.Get("tcp_mss").(string)
		speed := d.Get("speed").(string)
		mtuOverride := d.Get("mtu_override").(string)
		mtu := d.Get("mtu").(string)
		role := d.Get("role").(string)
		allowaccess := d.Get("allowaccess").(string)
		mode := d.Get("mode").(string)
		dnsServerOverride := d.Get("dns_server_override").(string)
		defaultgw := d.Get("defaultgw").(string)
		distance := d.Get("distance").(string)
		description := d.Get("description").(string)
		interfacef := d.Get("interface").(string)
		name := d.Get("name").(string)
		vdom := d.Get("vdom").(string)
		vlanid := d.Get("vlanid").(string)

		if typef == "vdom" {
			if vlanid == "" {
				return fmt.Errorf("Error creating Networking Interface Port, vdom is nil")
			}
		}

		if typef == "vlan" {
			if vlanid == "" {
				return fmt.Errorf("Error creating Networking Interface Port, vlanid is nil")
			}

			if interfacef == "" {
				return fmt.Errorf("Error creating Networking Interface Port, interface is nil")
			}
		}

		//Build input data by sdk
		i := &forticlient.JSONNetworkingInterfacePort{
			//Portname:             portname,
			Ipf:                  ip,
			Alias:                alias,
			Status:               status,
			DeviceIdentification: deviceIdentification,
			TCPMss:               tcpMss,
			Speed:                speed,
			MtuOverride:          mtuOverride,
			Mtu:                  mtu,
			Role:                 role,
			Allowaccess:          allowaccess,
			Mode:                 mode,
			DNSServerOverride:    dnsServerOverride,
			Defaultgw:            defaultgw,
			Distance:             distance,
			Description:          description,
			Type:                 typef,
			Interface:            interfacef,
			Name:                 name,
			Vdom:                 vdom,
			Vlanid:               vlanid,
		}

		//Call process by sdk
		o, err := c.CreateNetworkingInterfacePort(i)
		if err != nil {
			return fmt.Errorf("Error creating Networking Interface Port: %s", err)
		}

		//Set index for d
		log.Printf("FOS-id is %v\n", o.Mkey)
		d.SetId(o.Mkey)

	}
	return resourceNetworkingInterfacePortRead(d, m)
}

func resourceNetworkingInterfacePortUpdate(d *schema.ResourceData, m interface{}) error {
	mkey := d.Id()

	c := m.(*FortiClient).Client

	if c == nil {
		return fmt.Errorf("FortiOS connection did not initialize successfully!")
	}

	c.Retries = 1

	//Get Params from d
	ip := d.Get("ip").(string)
	alias := d.Get("alias").(string)
	status := d.Get("status").(string)
	deviceIdentification := d.Get("device_identification").(string)
	tcpMss := d.Get("tcp_mss").(string)
	speed := d.Get("speed").(string)
	mtuOverride := d.Get("mtu_override").(string)
	mtu := d.Get("mtu").(string)
	role := d.Get("role").(string)
	allowaccess := d.Get("allowaccess").(string)
	mode := d.Get("mode").(string)
	dnsServerOverride := d.Get("dns_server_override").(string)
	defaultgw := d.Get("defaultgw").(string)
	distance := d.Get("distance").(string)
	description := d.Get("description").(string)
	typef := d.Get("type").(string)
	interfacef := d.Get("interface").(string)
	name := d.Get("name").(string)
	vdom := d.Get("vdom").(string)
	vlanid := d.Get("vlanid").(string)

	if typef == "physical" {
		if name == "" {
			name = mkey
		}

		if vdom == "" {
			vdom = "root"
		}
	}

	//Build input data by sdk
	i := &forticlient.JSONNetworkingInterfacePort{
		//Portname:             portname,
		Ipf:                  ip,
		Alias:                alias,
		Status:               status,
		DeviceIdentification: deviceIdentification,
		TCPMss:               tcpMss,
		Speed:                speed,
		MtuOverride:          mtuOverride,
		Mtu:                  mtu,
		Role:                 role,
		Allowaccess:          allowaccess,
		Mode:                 mode,
		DNSServerOverride:    dnsServerOverride,
		Defaultgw:            defaultgw,
		Distance:             distance,
		Description:          description,
		Type:                 typef,
		Interface:            interfacef,
		Name:                 name,
		Vdom:                 vdom,
		Vlanid:               vlanid,
	}

	//Call process by sdk
	_, err := c.UpdateNetworkingInterfacePort(i, mkey)
	if err != nil {
		return fmt.Errorf("Error updating Networking Interface Port: %s", err)
	}

	return resourceNetworkingInterfacePortRead(d, m)
}

func resourceNetworkingInterfacePortDelete(d *schema.ResourceData, m interface{}) error {
	mkey := d.Id()

	c := m.(*FortiClient).Client

	if c == nil {
		return fmt.Errorf("FortiOS connection did not initialize successfully!")
	}

	c.Retries = 1

	typef := d.Get("type").(string)

	if typef == "physical" {
		d.SetId("")

		return nil
	}

	//Call process by sdk
	err := c.DeleteNetworkingInterfacePort(mkey)
	if err != nil {
		return fmt.Errorf("Error deleting Networking Interface Port: %s", err)
	}

	//Set index for d
	d.SetId("")

	return nil
}

func resourceNetworkingInterfacePortRead(d *schema.ResourceData, m interface{}) error {
	mkey := d.Id()

	c := m.(*FortiClient).Client

	if c == nil {
		return fmt.Errorf("FortiOS connection did not initialize successfully!")
	}

	c.Retries = 1

	//Call process by sdk
	o, err := c.ReadNetworkingInterfacePort(mkey)
	if err != nil {
		return fmt.Errorf("Error reading Networking Interface Port: %s", err)
	}

	if o == nil {
		log.Printf("[WARN] resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	//Refresh property
	//d.Set("portname", o.Portname)
	d.Set("ip", validateConvIPMask2CDIR(d.Get("ip").(string), o.Ipf))
	d.Set("alias", o.Alias)
	d.Set("status", o.Status)
	d.Set("device_identification", o.DeviceIdentification)
	d.Set("tcp_mss", o.TCPMss)
	d.Set("speed", o.Speed)
	d.Set("mtu_override", o.MtuOverride)
	d.Set("mtu", o.Mtu)
	d.Set("role", o.Role)
	d.Set("allowaccess", o.Allowaccess)
	d.Set("mode", o.Mode)
	d.Set("dns_server_override", o.DNSServerOverride)
	d.Set("defaultgw", o.Defaultgw)
	d.Set("distance", o.Distance)
	d.Set("description", o.Description)
	d.Set("type", o.Type)
	d.Set("interface", o.Interface)
	d.Set("name", o.Name)
	d.Set("vdom", o.Vdom)
	d.Set("vlanid", o.Vlanid)

	return nil
}
