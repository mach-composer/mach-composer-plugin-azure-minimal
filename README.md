# mach-composer-plugin-azure-minimal

An unopinionated plugin for Azure


## Usage

```yaml
mach_composer:
  version: 1
  plugins:
    azure:
      source: mach-composer/azure-minimal
      version: 0.1.0

global:
  environment: test
  cloud: azure
  terraform_config:
    azure_remote_state:
      resource_group: resourcegroupid
      storage_account: storageaccount
      container_name: container-name
    providers:
      azure: =3.43.0 # override
  azure:
    subscription_id: "subscription ID"
    tenant_id: "Tenant ID"
    resource_tags:
      My: "ABC"
      Tag: "def"

sites:
  - identifier: my-site
    azure:
      resource_group: my-rg #resource group in which the site is deployed, usually one per site
      resource_prefix: my-resource-prefix #prefix used for created resources in this group
```
