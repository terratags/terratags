terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
  
  tags = {
    Name        = "example-resource-group"
    Environment = "Production"
    Owner       = "DevOps-Team"
    Project     = "Terratags-Demo"
  }
}

# This resource supports tags and has all required tags
resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacct"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  
  tags = {
    Name        = "example-storage-account"
    Environment = "Production"
    Owner       = "DevOps-Team"
    Project     = "Terratags-Demo"
  }
}

# This resource does not support tags but will pass validation
# because Terratags recognizes it doesn't support tagging
resource "azurerm_role_definition" "example" {
  name        = "custom-role-definition"
  scope       = azurerm_resource_group.example.id
  description = "This is a custom role definition"

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    azurerm_resource_group.example.id
  ]
  
  # Note: This resource does not support tags
}
