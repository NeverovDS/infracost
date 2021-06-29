package azure

import (
	"strings"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

func GetAzureRMKubernetesClusterRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_kubernetes_cluster",
		RFunc: NewAzureRMKubernetesCluster,
	}
}

func NewAzureRMKubernetesCluster(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := lookupRegion(d, []string{})
	var costComponents []*schema.CostComponent
	var subResources []*schema.Resource

	skuTier := "Free"
	if d.Get("sku_tier").Type != gjson.Null {
		skuTier = d.Get("sku_tier").String()
	}

	if skuTier == "Paid" {
		costComponents = append(costComponents, &schema.CostComponent{
			Name:           "Uptime SLA",
			Unit:           "hours",
			UnitMultiplier: 1,
			HourlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(region),
				Service:       strPtr("Azure Kubernetes Service"),
				ProductFamily: strPtr("Compute"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "skuName", Value: strPtr("Standard")},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Consumption"),
			},
		})
	}

	nodeCount := decimal.NewFromInt(1)
	if d.Get("default_node_pool.0.node_count").Type != gjson.Null {
		nodeCount = decimal.NewFromInt(d.Get("default_node_pool.0.node_count").Int())
	}
	if u != nil && u.Get("default_node_pool.nodes").Exists() {
		nodeCount = decimal.NewFromInt(u.Get("default_node_pool.nodes").Int())
	}

	subResources = []*schema.Resource{
		aksClusterNodePool("default_node_pool", region, d.Get("default_node_pool.0"), nodeCount, u),
	}

	if d.Get("network_profile.0.load_balancer_sku").Type != gjson.Null {
		if strings.ToLower(d.Get("network_profile.0.load_balancer_sku").String()) == "standard" {
			location := region
			if strings.Contains(strings.ToLower(region), "usgov") {
				location = "US Gov"
			} else if strings.Contains(strings.ToLower(region), "china") {
				location = "Сhina"
			} else {
				location = "Global"
			}
			var monthlyDataProcessedGb *decimal.Decimal
			if u != nil && u.Get("load_balancer.monthly_data_processed_gb").Type != gjson.Null {
				monthlyDataProcessedGb = decimalPtr(decimal.NewFromInt(u.Get("load_balancer.monthly_data_processed_gb").Int()))
			}
			lbResource := schema.Resource{
				Name:           "Load Balancer",
				CostComponents: []*schema.CostComponent{dataProcessedCostComponent(location, monthlyDataProcessedGb)},
			}
			subResources = append(subResources, &lbResource)
		}
	}
	if d.Get("addon_profile.0.http_application_routing").Type != gjson.Null {
		if strings.ToLower(d.Get("addon_profile.0.http_application_routing.0.enabled").String()) == "true" {
			location := region
			if strings.HasPrefix(strings.ToLower(region), "usgov") {
				location = "US Gov Zone 1"
			} else if strings.HasPrefix(strings.ToLower(region), "germany") {
				location = "DE Zone 1"
			} else if strings.HasPrefix(strings.ToLower(region), "china") {
				location = "Zone 1 (China)"
			} else {
				location = "Zone 1"
			}

			dnsResource := schema.Resource{
				Name:           "DNS",
				CostComponents: []*schema.CostComponent{hostedPublicZoneCostComponent(location)},
			}
			subResources = append(subResources, &dnsResource)
		}
	}

	return &schema.Resource{
		Name:           d.Address,
		CostComponents: costComponents,
		SubResources:   subResources,
	}
}
