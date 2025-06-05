terraform {
  required_providers {
    azapi = {
      source  = "Azure/azapi"
      version = "~> 1.0"
    }
  }
}

provider "azapi" {
  # azapi supports default_tags at provider level
  default_tags = {
    managed_by = "terraform"
    project    = "terratags-example"
    owner      = "devops-team"
  }
}

# azapi_resource supports tags at resource level
resource "azapi_resource" "example" {
  type      = "Microsoft.Storage/storageAccounts@2022-05-01"
  name      = "examplestorageaccount"
  parent_id = azurerm_resource_group.example.id
  location  = azurerm_resource_group.example.location
  
  body = jsonencode({
    kind = "StorageV2"
    sku = {
      name = "Standard_LRS"
    }
  })

  tags = {
    environment = "development"
    department  = "research"
    name        = "example-storage"
  }
}
