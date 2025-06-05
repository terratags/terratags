terraform {
  required_providers {
    azapi = {
      source  = "Azure/azapi"
      version = "~> 1.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
  # Missing required tags: Owner and Project
  default_tags = {
    Environment = "Production"
    Name        = "Default-Name"
  }
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

# This resource will fail validation because provider default_tags are missing required tags
resource "azapi_resource" "example_storage" {
  type      = "Microsoft.Storage/storageAccounts@2022-05-01"
  name      = "examplestorageacct"
  parent_id = azurerm_resource_group.example.id
  location  = azurerm_resource_group.example.location
  
  # No additional tags specified at resource level
  # Inherits incomplete tags from provider default_tags
}
