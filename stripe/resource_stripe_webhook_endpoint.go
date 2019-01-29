package stripe

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"

	"log"
)

func resourceStripeWebhookEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceStripeWebhookEndpointCreate,
		Read:   resourceStripeWebhookEndpointRead,
		Update: resourceStripeWebhookEndpointUpdate,
		Delete: resourceStripeWebhookEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled_events": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceStripeWebhookEndpointCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	webhookEndpointURL := d.Get("url").(string)
	rawEnabledEvents := d.Get("enabled_events").([]interface{})
	webhookEndpointEnabledEvents := make([]*string, len(rawEnabledEvents))
	for i, v := range rawEnabledEvents {
		stringEvent := v.(string)
		webhookEndpointEnabledEvents[i] = &stringEvent
	}

	params := &stripe.WebhookEndpointParams{
		URL:           stripe.String(webhookEndpointURL),
		EnabledEvents: webhookEndpointEnabledEvents,
	}

	webhookEndpoint, err := client.WebhookEndpoints.New(params)

	if err != nil {
		return err
	} else {
		log.Printf("[INFO] Create wehbook endpoint: %s", webhookEndpointURL)
		d.SetId(webhookEndpoint.ID)
	}
	return nil
}

func resourceStripeWebhookEndpointRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	webhookEndpoint, err := client.WebhookEndpoints.Get(d.Id(), nil)

	if err != nil {
		return err
	} else {
		d.Set("url", webhookEndpoint.URL)
		d.Set("enabled_events", webhookEndpoint.EnabledEvents)
	}
	return nil
}

func resourceStripeWebhookEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	params := stripe.WebhookEndpointParams{}

	if d.HasChange("name") {
		params.URL = stripe.String(d.Get("url").(string))
	}
	if d.HasChange("enabled_events") {
		rawEnabledEvents := d.Get("enabled_events").([]interface{})
		enabledEvents := make([]*string, len(rawEnabledEvents))
		for i, v := range rawEnabledEvents {
			stringEvent := v.(string)
			enabledEvents[i] = &stringEvent
		}
		params.EnabledEvents = enabledEvents
	}

	_, err := client.WebhookEndpoints.Update(d.Id(), &params)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func resourceStripeWebhookEndpointDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.API)
	_, err := client.WebhookEndpoints.Del(d.Id(), nil)

	if err != nil {
		return err
	} else {
		d.SetId("")
		return nil
	}
}
