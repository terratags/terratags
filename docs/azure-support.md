# Azure Support in Terratags

Terratags now supports Azure resources through both the `azurerm` and `azapi` providers.

## Azurerm Provider

The `azurerm` provider supports tagging at the resource level. Terratags will automatically detect and manage tags for resources that support the `tags` attribute.

Example:

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
  
  tags = {
    environment = "production"
    department  = "finance"
  }
}
```

Note: The `azurerm` provider does not support default_tags at the provider level.

## Azapi Provider

The `azapi` provider supports tagging at both the provider level (via `default_tags`) and at the resource level.

Example with provider-level default tags:

```hcl
provider "azapi" {
  default_tags = {
    managed_by = "terraform"
    project    = "terratags-example"
  }
}
```

Example with resource-level tags:

```hcl
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
  }
}
```

## Supported Resources

Terratags automatically detects which Azure resources support tagging by analyzing the provider schemas. The list of taggable resources is generated during the build process.

For the most up-to-date list of supported resources, refer to the provider documentation:
- [Azurerm Provider](https://registry.terraform.io/providers/hashicorp/azurerm/4.31.0)
- [Azapi Provider](https://registry.terraform.io/providers/Azure/azapi/latest/docs)