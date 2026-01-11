# Virtual Machine Example

This example demonstrates deploying an Azure Linux Virtual Machine using wetwire-azure-go.

## Resources Created

- **VirtualMachine**: An Ubuntu 22.04 LTS Linux VM

## Prerequisites

Before deploying this VM, you need:

1. A Virtual Network with a subnet
2. A Network Interface (NIC) attached to the subnet
3. (Optional) A Public IP address attached to the NIC

## Configuration

| Property | Value | Description |
|----------|-------|-------------|
| Name | my-linux-vm | VM name |
| Location | eastus | Azure region |
| VM Size | Standard_B2s | 2 vCPUs, 4 GB RAM |
| OS | Ubuntu 22.04 LTS | Canonical's Ubuntu Server |
| Disk Type | Premium_LRS | Premium SSD |

## Build

Generate the ARM template:

```bash
wetwire-azure build -o template.json
```

## Deploy

### 1. Create Prerequisites

```bash
# Login to Azure
az login

# Create resource group
az group create --name vm-example-rg --location eastus

# Create virtual network
az network vnet create \
  --resource-group vm-example-rg \
  --name my-vnet \
  --address-prefix 10.0.0.0/16 \
  --subnet-name my-subnet \
  --subnet-prefix 10.0.1.0/24

# Create public IP
az network public-ip create \
  --resource-group vm-example-rg \
  --name my-public-ip \
  --allocation-method Static

# Create network interface
az network nic create \
  --resource-group vm-example-rg \
  --name my-nic \
  --vnet-name my-vnet \
  --subnet my-subnet \
  --public-ip-address my-public-ip
```

### 2. Update Network Interface ID

Update the `networkIfaceID` in `main.go` with your actual NIC resource ID:

```bash
az network nic show --resource-group vm-example-rg --name my-nic --query id -o tsv
```

### 3. Deploy the VM

```bash
az deployment group create \
  --resource-group vm-example-rg \
  --template-file template.json
```

## Customization

### Change VM Size

```go
HardwareProfile: compute.HardwareProfile{
    VMSize: "Standard_D2s_v3",  // 2 vCPUs, 8 GB RAM
}
```

Common VM sizes:
- `Standard_B1s` - Burstable, 1 vCPU, 1 GB RAM
- `Standard_B2s` - Burstable, 2 vCPUs, 4 GB RAM
- `Standard_D2s_v3` - General purpose, 2 vCPUs, 8 GB RAM
- `Standard_D4s_v3` - General purpose, 4 vCPUs, 16 GB RAM

### Use SSH Key Authentication

Replace password authentication with SSH keys:

```go
disablePwd := true
sshPath := "/home/azureuser/.ssh/authorized_keys"
sshKey := "ssh-rsa AAAAB3NzaC1yc2E..."

OSProfile: &compute.OSProfile{
    ComputerName:  &computerName,
    AdminUsername: &adminUsername,
    LinuxConfiguration: &compute.LinuxConfiguration{
        DisablePasswordAuthentication: &disablePwd,
        SSH: &compute.SSHConfiguration{
            PublicKeys: []compute.SSHPublicKey{
                {
                    Path:    &sshPath,
                    KeyData: &sshKey,
                },
            },
        },
    },
}
```

### Add Tags

```go
var LinuxVM = compute.VirtualMachine{
    Name:     "my-linux-vm",
    Location: "eastus",
    Tags: map[string]string{
        "environment": "development",
        "os":          "ubuntu",
    },
    Properties: compute.VirtualMachineProperties{...},
}
```

## Cleanup

```bash
az group delete --name vm-example-rg --yes
```

## Security Notes

- The example uses password authentication for simplicity
- For production, use SSH key authentication
- Store credentials in Azure Key Vault, not in code
- Use managed identities for Azure resource access
