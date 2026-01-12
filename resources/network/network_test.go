// Package network provides Azure network resource types
package network

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVirtualNetwork(t *testing.T) {
	vnet := NewVirtualNetwork("my-vnet", "eastus", []string{"10.0.0.0/16"})

	assert.Equal(t, "my-vnet", vnet.Name)
	assert.Equal(t, "Microsoft.Network/virtualNetworks", vnet.Type)
	assert.Equal(t, "2021-05-01", vnet.APIVersion)
	assert.Equal(t, "eastus", vnet.Location)
	assert.Equal(t, []string{"10.0.0.0/16"}, vnet.Properties.AddressSpace.AddressPrefixes)
}

func TestVirtualNetwork_WithTags(t *testing.T) {
	vnet := NewVirtualNetwork("my-vnet", "eastus", []string{"10.0.0.0/16"}).
		WithTags(map[string]string{"env": "prod"})

	assert.Equal(t, "prod", vnet.Tags["env"])
}

func TestVirtualNetwork_WithSubnet(t *testing.T) {
	vnet := NewVirtualNetwork("my-vnet", "eastus", []string{"10.0.0.0/16"}).
		WithSubnet("default", "10.0.0.0/24")

	require.Len(t, vnet.Properties.Subnets, 1)
	assert.Equal(t, "default", vnet.Properties.Subnets[0].Name)
	assert.Equal(t, "10.0.0.0/24", vnet.Properties.Subnets[0].Properties.AddressPrefix)
}

func TestVirtualNetwork_JSON(t *testing.T) {
	vnet := NewVirtualNetwork("my-vnet", "eastus", []string{"10.0.0.0/16"}).
		WithSubnet("default", "10.0.0.0/24")

	data, err := json.Marshal(vnet)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "my-vnet", result["name"])
	assert.Equal(t, "Microsoft.Network/virtualNetworks", result["type"])
}

func TestNewSubnet(t *testing.T) {
	subnet := NewSubnet("my-subnet", "10.0.1.0/24")

	assert.Equal(t, "my-subnet", subnet.Name)
	assert.Equal(t, "10.0.1.0/24", subnet.Properties.AddressPrefix)
}

func TestSubnet_WithNSG(t *testing.T) {
	subnet := NewSubnet("my-subnet", "10.0.1.0/24").
		WithNSG("/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/my-nsg")

	require.NotNil(t, subnet.Properties.NetworkSecurityGroup)
	assert.Equal(t, "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/my-nsg",
		*subnet.Properties.NetworkSecurityGroup.ID)
}

func TestNewNetworkInterface(t *testing.T) {
	nic := NewNetworkInterface("my-nic", "eastus")

	assert.Equal(t, "my-nic", nic.Name)
	assert.Equal(t, "Microsoft.Network/networkInterfaces", nic.Type)
	assert.Equal(t, "2021-05-01", nic.APIVersion)
	assert.Equal(t, "eastus", nic.Location)
}

func TestNetworkInterface_WithIPConfiguration(t *testing.T) {
	nic := NewNetworkInterface("my-nic", "eastus").
		WithIPConfiguration("ipconfig1", "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/default", true)

	require.Len(t, nic.Properties.IPConfigurations, 1)
	assert.Equal(t, "ipconfig1", nic.Properties.IPConfigurations[0].Name)
	assert.Equal(t, "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/default",
		*nic.Properties.IPConfigurations[0].Properties.Subnet.ID)
	assert.True(t, *nic.Properties.IPConfigurations[0].Properties.Primary)
}

func TestNetworkInterface_WithPublicIP(t *testing.T) {
	nic := NewNetworkInterface("my-nic", "eastus").
		WithIPConfiguration("ipconfig1", "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/default", true).
		WithPublicIP("ipconfig1", "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/my-pip")

	require.Len(t, nic.Properties.IPConfigurations, 1)
	require.NotNil(t, nic.Properties.IPConfigurations[0].Properties.PublicIPAddress)
	assert.Equal(t, "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/my-pip",
		*nic.Properties.IPConfigurations[0].Properties.PublicIPAddress.ID)
}

func TestNetworkInterface_WithNSG(t *testing.T) {
	nic := NewNetworkInterface("my-nic", "eastus").
		WithNSG("/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/my-nsg")

	require.NotNil(t, nic.Properties.NetworkSecurityGroup)
	assert.Equal(t, "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/my-nsg",
		*nic.Properties.NetworkSecurityGroup.ID)
}

func TestNewPublicIPAddress(t *testing.T) {
	pip := NewPublicIPAddress("my-pip", "eastus", "Static", "Standard")

	assert.Equal(t, "my-pip", pip.Name)
	assert.Equal(t, "Microsoft.Network/publicIPAddresses", pip.Type)
	assert.Equal(t, "2021-05-01", pip.APIVersion)
	assert.Equal(t, "eastus", pip.Location)
	assert.Equal(t, "Static", pip.Properties.PublicIPAllocationMethod)
	assert.Equal(t, "Standard", pip.SKU.Name)
}

func TestPublicIPAddress_WithDNSLabel(t *testing.T) {
	pip := NewPublicIPAddress("my-pip", "eastus", "Static", "Standard").
		WithDNSLabel("myapp")

	require.NotNil(t, pip.Properties.DNSSettings)
	assert.Equal(t, "myapp", *pip.Properties.DNSSettings.DomainNameLabel)
}

func TestNewNetworkSecurityGroup(t *testing.T) {
	nsg := NewNetworkSecurityGroup("my-nsg", "eastus")

	assert.Equal(t, "my-nsg", nsg.Name)
	assert.Equal(t, "Microsoft.Network/networkSecurityGroups", nsg.Type)
	assert.Equal(t, "2021-05-01", nsg.APIVersion)
	assert.Equal(t, "eastus", nsg.Location)
}

func TestNetworkSecurityGroup_WithRule(t *testing.T) {
	nsg := NewNetworkSecurityGroup("my-nsg", "eastus").
		WithRule("allow-ssh", 100, "Inbound", "Allow", "Tcp", "*", "22", "*", "*")

	require.Len(t, nsg.Properties.SecurityRules, 1)
	rule := nsg.Properties.SecurityRules[0]
	assert.Equal(t, "allow-ssh", rule.Name)
	assert.Equal(t, 100, rule.Properties.Priority)
	assert.Equal(t, "Inbound", rule.Properties.Direction)
	assert.Equal(t, "Allow", rule.Properties.Access)
	assert.Equal(t, "Tcp", rule.Properties.Protocol)
	assert.Equal(t, "*", rule.Properties.SourcePortRange)
	assert.Equal(t, "22", rule.Properties.DestinationPortRange)
	assert.Equal(t, "*", rule.Properties.SourceAddressPrefix)
	assert.Equal(t, "*", rule.Properties.DestinationAddressPrefix)
}

func TestNetworkSecurityGroup_JSON(t *testing.T) {
	nsg := NewNetworkSecurityGroup("my-nsg", "eastus").
		WithRule("allow-http", 100, "Inbound", "Allow", "Tcp", "*", "80", "*", "*").
		WithRule("allow-https", 110, "Inbound", "Allow", "Tcp", "*", "443", "*", "*")

	data, err := json.Marshal(nsg)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "my-nsg", result["name"])
	assert.Equal(t, "Microsoft.Network/networkSecurityGroups", result["type"])

	props := result["properties"].(map[string]interface{})
	rules := props["securityRules"].([]interface{})
	assert.Len(t, rules, 2)
}
