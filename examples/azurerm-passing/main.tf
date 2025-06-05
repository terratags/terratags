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

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
  
  tags = {
    Name        = "example-vnet"
    Environment = "Production"
    Owner       = "Network-Team"
    Project     = "Terratags-Demo"
  }
}
