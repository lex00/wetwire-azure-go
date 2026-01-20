Generate ARM template for enterprise Azure deployment:

- VNet: 10.0.0.0/16, subnets: web (10.0.1.0/24), app (10.0.2.0/24)
- NSGs: web-nsg (80/443 ingress), app-nsg (8080 from web)
- VMSS: web-vmss (B2s, 2 instances, web-subnet)
- Storage: Standard_LRS, StorageV2, HTTPS only
- Public IP: Standard SKU, static

Location: eastus. Single JSON file. No documentation.
