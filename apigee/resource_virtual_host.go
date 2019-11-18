package apigee

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tibers/go-apigee-edge"
)

func resourceVirtualHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualHostCreate,
		Read:   resourceVirtualHostRead,
		Update: resourceVirtualHostUpdate,
		Delete: resourceVirtualHostDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVirtualHostImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hostAliases": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"baseUrl": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sSLInfo": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"port": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"properties": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"retryOptions": {
				Type:     schema.TypeString,
				Required: false,
			},
			"listenOptions": {
				Type:     schema.TypeString,
				Required: false,
			},
		},
	}
}

func resourceVirtualHostCreate(d *schema.ResourceData, meta interface{}) error {

	log.Print("[DEBUG] resourceVirtualHostCreate START")

	client := meta.(*apigee.EdgeClient)

	u1, _ := uuid.NewV4()
	d.SetId(u1.String())

	VirtualHostData, err := setVirtualHostData(d)
	if err != nil {
		log.Printf("[ERROR] resourceVirtualHostCreate error in setVirtualHostData: %s", err.Error())
		return fmt.Errorf("[ERROR] resourceVirtualHostCreate error in setVirtualHostData: %s", err.Error())
	}

	_, _, e := client.VirtualHosts.Create(VirtualHostData, d.Get("env").(string))
	if e != nil {
		log.Printf("[ERROR] resourceVirtualHostCreate error in create: %s", e.Error())
		return fmt.Errorf("[ERROR] resourceVirtualHostCreate error in create: %s", e.Error())
	}

	return resourceVirtualHostRead(d, meta)
}

func resourceVirtualHostImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	log.Print("[DEBUG] resourceVirtualHostImport START")
	client := meta.(*apigee.EdgeClient)
	splits := strings.Split(d.Id(), "_")
	if len(splits) < 1 {
		return []*schema.ResourceData{}, fmt.Errorf("[ERR] Wrong format of resource: %s. Please follow '{name}_{env}'", d.Id())
	}

	nameOffset := len(splits[len(splits)-1])
	name := d.Id()[:(len(d.Id())-nameOffset)-1]
	IDEnv := splits[len(splits)-1]

	VirtualHostData, _, err := client.VirtualHosts.Get(name, IDEnv)
	if err != nil {
		log.Printf("[ERROR] resourceVirtualHostImport error getting target servers: %s", err.Error())
		if strings.Contains(err.Error(), "404 ") {
			return []*schema.ResourceData{}, fmt.Errorf("[Error] resourceVirtualHostImport 404 encountered.  Removing state for target server: %#v", name)
		}
		return []*schema.ResourceData{}, fmt.Errorf("[ERROR] resourceVirtualHostImport error getting target servers: %s", err.Error())
	}

	d.Set("name", VirtualHostData.Name)
	d.Set("host", VirtualHostData.Host)
	d.Set("enabled", VirtualHostData.Enabled)
	d.Set("port", VirtualHostData.Port)
	d.Set("env", IDEnv)

	if VirtualHostData.SSLInfo != nil {
		protocols := flattenStringList(VirtualHostData.SSLInfo.Protocols)
		ciphers := flattenStringList(VirtualHostData.SSLInfo.Ciphers)

		d.Set("ssl_info.0.ssl_enabled", VirtualHostData.SSLInfo.SSLEnabled)
		d.Set("ssl_info.0.client_auth_enabled", VirtualHostData.SSLInfo.ClientAuthEnabled)
		d.Set("ssl_info.0.key_store", VirtualHostData.SSLInfo.KeyStore)
		d.Set("ssl_info.0.trust_store", VirtualHostData.SSLInfo.TrustStore)
		d.Set("ssl_info.0.key_alias", VirtualHostData.SSLInfo.KeyAlias)
		d.Set("ssl_info.0.ciphers", ciphers)
		d.Set("ssl_info.0.ignore_validation_errors", VirtualHostData.SSLInfo.IgnoreValidationErrors)
		d.Set("ssl_info.0.protocols", protocols)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceVirtualHostRead(d *schema.ResourceData, meta interface{}) error {

	log.Print("[DEBUG] resourceVirtualHostRead START")
	client := meta.(*apigee.EdgeClient)

	VirtualHostData, _, err := client.VirtualHosts.Get(d.Get("name").(string), d.Get("env").(string))
	if err != nil {
		log.Printf("[ERROR] resourceVirtualHostRead error getting target servers: %s", err.Error())
		if strings.Contains(err.Error(), "404 ") {
			log.Printf("[DEBUG] resourceVirtualHostRead 404 encountered.  Removing state for target server: %#v", d.Get("name").(string))
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("[ERROR] resourceVirtualHostRead error getting target servers: %s", err.Error())
		}
	}

	port_str := strconv.Itoa(VirtualHostData.Port)

	d.Set("name", VirtualHostData.Name)
	d.Set("host", VirtualHostData.Host)
	d.Set("enabled", VirtualHostData.Enabled)
	d.Set("port", port_str)

	if VirtualHostData.SSLInfo != nil {
		protocols := flattenStringList(VirtualHostData.SSLInfo.Protocols)
		ciphers := flattenStringList(VirtualHostData.SSLInfo.Ciphers)

		d.Set("ssl_info.0.ssl_enabled", VirtualHostData.SSLInfo.SSLEnabled)
		d.Set("ssl_info.0.client_auth_enabled", VirtualHostData.SSLInfo.ClientAuthEnabled)
		d.Set("ssl_info.0.key_store", VirtualHostData.SSLInfo.KeyStore)
		d.Set("ssl_info.0.trust_store", VirtualHostData.SSLInfo.TrustStore)
		d.Set("ssl_info.0.key_alias", VirtualHostData.SSLInfo.KeyAlias)
		d.Set("ssl_info.0.ciphers", ciphers)
		d.Set("ssl_info.0.ignore_validation_errors", VirtualHostData.SSLInfo.IgnoreValidationErrors)
		d.Set("ssl_info.0.protocols", protocols)
	}

	return nil
}

func resourceVirtualHostUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Print("[DEBUG] resourceVirtualHostUpdate START")

	client := meta.(*apigee.EdgeClient)

	VirtualHostData, err := setVirtualHostData(d)
	if err != nil {
		log.Printf("[ERROR] resourceVirtualHostUpdate error in setVirtualHostData: %s", err.Error())
		return fmt.Errorf("[ERROR] resourceVirtualHostUpdate error in setVirtualHostData: %s", err.Error())
	}

	_, _, e := client.VirtualHosts.Update(VirtualHostData, d.Get("env").(string))
	if e != nil {
		log.Printf("[ERROR] resourceVirtualHostUpdate error in update: %s", e.Error())
		return fmt.Errorf("[ERROR] resourceVirtualHostUpdate error in update: %s", e.Error())
	}

	return resourceVirtualHostRead(d, meta)
}

func resourceVirtualHostDelete(d *schema.ResourceData, meta interface{}) error {

	log.Print("[DEBUG] resourceVirtualHostDelete START")
	client := meta.(*apigee.EdgeClient)

	_, err := client.VirtualHosts.Delete(d.Get("name").(string), d.Get("env").(string))
	if err != nil {
		log.Printf("[ERROR] resourceVirtualHostDelete error in delete: %s", err.Error())
		return fmt.Errorf("[ERROR] resourceVirtualHostDelete error in delete: %s", err.Error())
	}

	return nil
}

func setVirtualHostData(d *schema.ResourceData) (apigee.VirtualHost, error) {

	log.Print("[DEBUG] setVirtualHostData START")

	port_int, _ := strconv.Atoi(d.Get("port").(string))

	var ssl_info *apigee.SSLInfo
	if d.Get("ssl_info") != nil {
		ciphers := []string{""}
		if d.Get("ssl_info.0.ciphers") != nil {
			ciphers = getStringList("ssl_info.0.ciphers", d)
		}

		protocols := []string{""}
		if d.Get("ssl_info.0.protocols") != nil {
			protocols = getStringList("ssl_info.0.protocols", d)
		}

		ssl_info = &apigee.SSLInfo{
			SSLEnabled:        d.Get("ssl_info.0.ssl_enabled").(string),
			ClientAuthEnabled: d.Get("ssl_info.0.client_auth_enabled").(string),
			KeyStore:          d.Get("ssl_info.0.key_store").(string),
			TrustStore:        d.Get("ssl_info.0.trust_store").(string),
			KeyAlias:          d.Get("ssl_info.0.key_alias").(string),
			Ciphers:           ciphers,
			//Ciphers: d.Get("ssl_info.0.ciphers").([]string),
			IgnoreValidationErrors: d.Get("ssl_info.0.ignore_validation_errors").(bool),
			Protocols:              protocols,
		}
	}

	VirtualHost := apigee.VirtualHost{
		Name:    d.Get("name").(string),
		Host:    d.Get("host").(string),
		Enabled: d.Get("enabled").(bool),
		Port:    port_int,
		SSLInfo: ssl_info,
	}

	return VirtualHost, nil
}
