package exoscale

import (
	"context"

	"github.com/exoscale/egoscale"
	"github.com/hashicorp/terraform/helper/schema"
)

func securityGroupResource() *schema.Resource {
	return &schema.Resource{
		Create: createSecurityGroup,
		Exists: existsSecurityGroup,
		Read:   readSecurityGroup,
		Update: updateSecurityGroup,
		Delete: deleteSecurityGroup,

		Importer: &schema.ResourceImporter{
			State: importSecurityGroup,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Read:   schema.DefaultTimeout(defaultTimeout),
			Update: schema.DefaultTimeout(defaultTimeout),
			Delete: schema.DefaultTimeout(defaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Removed:  "Tags cannot be set on security groups for the time being",
			},
		},
	}
}

func createSecurityGroup(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()

	client := GetComputeClient(meta)

	resp, err := client.RequestWithContext(ctx, &egoscale.CreateSecurityGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return err
	}

	sg := resp.(*egoscale.SecurityGroup)

	d.SetId(sg.ID.String())
	return readSecurityGroup(d, meta)
}

func existsSecurityGroup(d *schema.ResourceData, meta interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutRead))
	defer cancel()

	client := GetComputeClient(meta)

	id, err := egoscale.ParseUUID(d.Id())
	if err != nil {
		return false, err
	}

	sg := &egoscale.SecurityGroup{
		ID: id,
	}
	if err := client.GetWithContext(ctx, sg); err != nil {
		e := handleNotFound(d, err)
		return d.Id() != "", e
	}
	return true, nil
}

func readSecurityGroup(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutRead))
	defer cancel()

	client := GetComputeClient(meta)

	id, err := egoscale.ParseUUID(d.Id())
	if err != nil {
		return err
	}

	sg := &egoscale.SecurityGroup{
		ID: id,
	}
	if err := client.GetWithContext(ctx, sg); err != nil {
		return handleNotFound(d, err)
	}

	return applySecurityGroup(d, sg)
}

func updateSecurityGroup(d *schema.ResourceData, meta interface{}) error {
	return readSecurityGroup(d, meta)
}

func deleteSecurityGroup(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()

	client := GetComputeClient(meta)
	err := client.BooleanRequestWithContext(ctx, &egoscale.DeleteSecurityGroup{
		Name: d.Get("name").(string),
	})

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func importSecurityGroup(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutRead))
	defer cancel()

	client := GetComputeClient(meta)

	securityGroup := &egoscale.SecurityGroup{}

	id, err := egoscale.ParseUUID(d.Id())
	if err != nil {
		securityGroup.Name = d.Id()
	} else {
		securityGroup.ID = id
	}

	if err := client.GetWithContext(ctx, securityGroup); err != nil {
		return nil, err
	}

	if err := applySecurityGroup(d, securityGroup); err != nil {
		return nil, err
	}

	// Create all the rulez!
	ruleLength := len(securityGroup.EgressRule) + len(securityGroup.IngressRule)
	resources := make([]*schema.ResourceData, 0, 1+ruleLength)
	resources = append(resources, d)

	for _, rule := range securityGroup.EgressRule {
		resource := securityGroupRuleResource()
		d := resource.Data(nil)
		d.SetType("exoscale_security_group_rule")
		d.Set("type", "EGRESS")
		err := applySecurityGroupRule(d, securityGroup, rule)
		if err != nil {
			return nil, err
		}

		resources = append(resources, d)
	}
	for _, rule := range securityGroup.IngressRule {
		resource := securityGroupRuleResource()
		d := resource.Data(nil)
		d.SetType("exoscale_security_group_rule")
		d.Set("type", "INGRESS")
		err := applySecurityGroupRule(d, securityGroup, (egoscale.EgressRule)(rule))
		if err != nil {
			return nil, err
		}

		resources = append(resources, d)
	}

	return resources, nil
}

func applySecurityGroup(d *schema.ResourceData, securityGroup *egoscale.SecurityGroup) error {
	d.SetId(securityGroup.ID.String())
	d.Set("name", securityGroup.Name)
	d.Set("description", securityGroup.Description)
	return nil
}
