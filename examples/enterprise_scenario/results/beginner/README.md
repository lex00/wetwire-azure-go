# Multi-Tier Enterprise Web Application - ARM Template

This ARM template deploys a secure, scalable multi-tier web application infrastructure on Azure.

## Architecture Overview

```
Internet
    |
    v
[Load Balancer (Public IP)]
    |
    v
[Web Tier - VM Scale Set]
    |
    v
[App Tier Subnet]
    |
    v
[Data Tier Subnet]
    |
    v
[Storage Account]
```

## What Gets Deployed

### 1. **Network Infrastructure**
- **Virtual Network (VNet)**: 10.0.0.0/16
  - **Web Subnet**: 10.0.1.0/24 - Hosts web servers
  - **App Subnet**: 10.0.2.0/24 - For application servers
  - **Data Subnet**: 10.0.3.0/24 - For databases

### 2. **Security Components**
- **Web NSG**: Allows HTTP (80), HTTPS (443), and SSH (22) from internet
- **App NSG**: Only allows traffic from web subnet on port 8080, denies internet
- **Data NSG**: Only allows traffic from app subnet on port 3306, denies everything else

### 3. **Web Servers**
- **VM Scale Set**: Automatically scales between 1-10 instances
- **Ubuntu 18.04 LTS** with NGINX pre-installed
- **Load Balancer**: Distributes traffic across all instances
- **Public IP**: Static IP with DNS name for external access

### 4. **Storage**
- **Storage Account**: Secure blob and file storage
- **Encryption**: TLS 1.2 minimum, encrypted at rest
- **Network Rules**: Only accessible from web and app subnets

## Security Features

✅ **Network Isolation**: Three-tier subnet architecture
✅ **NSG Rules**: Granular traffic control between tiers
✅ **Storage Security**: HTTPS-only, encrypted, network-restricted
✅ **Load Balancing**: Health probes and automatic failover
✅ **Scalability**: Auto-scaling web tier

## Prerequisites

- Azure subscription
- Azure CLI installed ([Install Guide](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli))
- Resource group created

## Deployment Instructions

### Option 1: Using Azure CLI

```bash
# Login to Azure
az login

# Create a resource group
az group create \
  --name my-webapp-rg \
  --location eastus

# Deploy the template
az deployment group create \
  --resource-group my-webapp-rg \
  --template-file template.json \
  --parameters vmAdminUsername=azureuser \
  --parameters vmAdminPassword='YourSecurePassword123!' \
  --parameters instanceCount=2
```

### Option 2: Using Azure Portal

1. Login to [Azure Portal](https://portal.azure.com)
2. Search for "Deploy a custom template"
3. Click "Build your own template in the editor"
4. Copy and paste the contents of `template.json`
5. Click "Save"
6. Fill in the parameters:
   - **Resource Group**: Select or create new
   - **Location**: East US (recommended)
   - **VM Admin Username**: azureuser
   - **VM Admin Password**: Your secure password
   - **Instance Count**: 2 (or more)
7. Click "Review + create" then "Create"

### Option 3: Using PowerShell

```powershell
# Login to Azure
Connect-AzAccount

# Create a resource group
New-AzResourceGroup -Name my-webapp-rg -Location eastus

# Deploy the template
New-AzResourceGroupDeployment `
  -ResourceGroupName my-webapp-rg `
  -TemplateFile template.json `
  -vmAdminUsername azureuser `
  -vmAdminPassword (ConvertTo-SecureString "YourSecurePassword123!" -AsPlainText -Force) `
  -instanceCount 2
```

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `vmAdminUsername` | string | azureuser | Administrator username for VMs |
| `vmAdminPassword` | securestring | - | Administrator password (required) |
| `vmSize` | string | Standard_B2s | VM size for web servers |
| `instanceCount` | int | 2 | Number of web server instances (1-10) |

## Outputs

After deployment completes, you'll receive:

- **publicIPAddress**: The public IP address of your load balancer
- **publicFQDN**: The fully qualified domain name (e.g., webapp-abc123.eastus.cloudapp.azure.com)
- **storageAccountName**: Name of the created storage account
- **vnetName**: Name of the virtual network
- **loadBalancerName**: Name of the load balancer

## Accessing Your Application

Once deployed, access your web application:

```bash
# Get the public FQDN
az deployment group show \
  --resource-group my-webapp-rg \
  --name <deployment-name> \
  --query properties.outputs.publicFQDN.value -o tsv

# Open in browser
http://<your-fqdn>
```

You should see: "Web Server Instance: webvm000000"

## Customization

### Change Instance Count

```bash
az vmss scale \
  --resource-group my-webapp-rg \
  --name web-vmss \
  --new-capacity 5
```

### Update VM Size

Modify the `vmSize` parameter during deployment:

```bash
--parameters vmSize=Standard_D2s_v3
```

### Add Custom Application

The template includes a basic NGINX installation. To deploy your application:

1. SSH into a VM instance (use the public IP)
2. Upload your application files
3. Configure NGINX to serve your app

Or modify the `customData` section in the template to automate deployment.

## Cost Estimation

Approximate monthly costs (East US region):

- **VM Scale Set** (2x Standard_B2s): ~$60/month
- **Load Balancer** (Standard): ~$22/month
- **Storage Account** (LRS, 100GB): ~$2/month
- **Public IP** (Static): ~$4/month
- **Bandwidth**: Variable

**Total**: ~$88-100/month (approximate)

## Monitoring

Enable monitoring in Azure Portal:

1. Navigate to your VM Scale Set
2. Click "Diagnostic settings"
3. Enable metrics and logs
4. View in Azure Monitor

## Cleanup

To delete all resources:

```bash
az group delete --name my-webapp-rg --yes --no-wait
```

**Warning**: This will permanently delete all resources in the resource group.

## Troubleshooting

### Deployment Fails

- **Check password**: Must meet complexity requirements (uppercase, lowercase, number, symbol)
- **Quota limits**: Ensure your subscription has available VM quota
- **Region availability**: Try a different region if resources aren't available

### Can't Access Web Application

- **Wait for deployment**: VMs can take 5-10 minutes to fully initialize
- **Check NSG rules**: Ensure port 80/443 are open
- **Verify health probe**: Load balancer needs healthy backends

### VMSS Not Scaling

- **Manual mode**: This template uses manual upgrade policy
- **Scale manually**: Use Azure CLI or Portal to adjust capacity
- **Auto-scale**: Add auto-scale rules in Portal after deployment

## Next Steps

1. **Add SSL/TLS**: Configure HTTPS with Azure Key Vault certificates
2. **Deploy Application**: Add your web application code
3. **Configure Monitoring**: Set up Azure Monitor and alerts
4. **Backup Strategy**: Enable Azure Backup for VMs and storage
5. **CI/CD Pipeline**: Integrate with Azure DevOps or GitHub Actions

## Architecture Diagram

```
┌─────────────────────────────────────────────────────┐
│                    Internet                          │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│  Load Balancer + Public IP (Static)                 │
│  - HTTP: 80  HTTPS: 443                             │
└──────────────────────┬──────────────────────────────┘
                       │
        ┌──────────────┴──────────────┐
        │   enterprise-vnet 10.0.0.0/16│
        │                              │
        │  ┌────────────────────────┐  │
        │  │ Web Subnet (10.0.1.0/24)│ │
        │  │ [Web NSG]               │  │
        │  │ ┌──────────────────┐   │  │
        │  │ │ VM Scale Set     │   │  │
        │  │ │ - NGINX Servers  │   │  │
        │  │ │ - Auto-scaling   │   │  │
        │  │ └──────────────────┘   │  │
        │  └────────────────────────┘  │
        │            │                  │
        │            ▼                  │
        │  ┌────────────────────────┐  │
        │  │ App Subnet (10.0.2.0/24)│ │
        │  │ [App NSG]               │  │
        │  │ (Ready for app servers) │  │
        │  └────────────────────────┘  │
        │            │                  │
        │            ▼                  │
        │  ┌────────────────────────┐  │
        │  │Data Subnet (10.0.3.0/24)│ │
        │  │ [Data NSG]              │  │
        │  │ (Ready for databases)   │  │
        │  └────────────────────────┘  │
        └──────────────────────────────┘
                       │
                       ▼
        ┌──────────────────────────────┐
        │   Storage Account            │
        │   - Encrypted (TLS 1.2)      │
        │   - VNet restricted          │
        └──────────────────────────────┘
```

## Support

For issues or questions:
- Azure Documentation: https://docs.microsoft.com/azure
- Azure Support: https://azure.microsoft.com/support
- Community Forums: https://docs.microsoft.com/answers

## License

This template is provided as-is for educational and deployment purposes.
