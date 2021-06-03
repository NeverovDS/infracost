provider "azurerm" {
  skip_provider_registration = true
  features {}
}
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_dns_txt_record" "over1B" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300

  record {
    value = "google-site-authenticator"
  }

  record {
    value = "more site information here"
  }

  tags = {
    Environment = "Production"
  }
}
resource "azurerm_dns_txt_record" "first1B" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300

  record {
    value = "google-site-authenticator"
  }

  record {
    value = "more site information here"
  }

  tags = {
    Environment = "Production"
  }
}
resource "azurerm_dns_txt_record" "withoutUsage" {
  name                = "test"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300

  record {
    value = "google-site-authenticator"
  }

  record {
    value = "more site information here"
  }

  tags = {
    Environment = "Production"
  }
}