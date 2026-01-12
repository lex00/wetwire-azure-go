// Package compute provides Azure compute resource types
package compute

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVirtualMachine(t *testing.T) {
	vm := NewVirtualMachine("my-vm", "eastus", "Standard_D2s_v3")

	assert.Equal(t, "my-vm", vm.Name)
	assert.Equal(t, "Microsoft.Compute/virtualMachines", vm.Type)
	assert.Equal(t, "2021-07-01", vm.APIVersion)
	assert.Equal(t, "eastus", vm.Location)
	assert.Equal(t, "Standard_D2s_v3", vm.Properties.HardwareProfile.VMSize)
	assert.Equal(t, "FromImage", vm.Properties.StorageProfile.OSDisk.CreateOption)
}

func TestVirtualMachine_WithTags(t *testing.T) {
	vm := NewVirtualMachine("my-vm", "eastus", "Standard_D2s_v3").
		WithTags(map[string]string{"env": "prod", "team": "platform"})

	assert.Equal(t, "prod", vm.Tags["env"])
	assert.Equal(t, "platform", vm.Tags["team"])
}

func TestVirtualMachine_WithImage(t *testing.T) {
	vm := NewVirtualMachine("my-vm", "eastus", "Standard_D2s_v3").
		WithImage("Canonical", "UbuntuServer", "18.04-LTS", "latest")

	require.NotNil(t, vm.Properties.StorageProfile.ImageReference)
	assert.Equal(t, "Canonical", *vm.Properties.StorageProfile.ImageReference.Publisher)
	assert.Equal(t, "UbuntuServer", *vm.Properties.StorageProfile.ImageReference.Offer)
	assert.Equal(t, "18.04-LTS", *vm.Properties.StorageProfile.ImageReference.SKU)
	assert.Equal(t, "latest", *vm.Properties.StorageProfile.ImageReference.Version)
}

func TestVirtualMachine_WithNetworkInterface(t *testing.T) {
	nicID := "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/my-nic"
	vm := NewVirtualMachine("my-vm", "eastus", "Standard_D2s_v3").
		WithNetworkInterface(nicID, true)

	require.Len(t, vm.Properties.NetworkProfile.NetworkInterfaces, 1)
	assert.Equal(t, nicID, vm.Properties.NetworkProfile.NetworkInterfaces[0].ID)
	assert.True(t, *vm.Properties.NetworkProfile.NetworkInterfaces[0].Primary)
}

func TestVirtualMachine_MultipleNetworkInterfaces(t *testing.T) {
	nic1ID := "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/nic1"
	nic2ID := "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/nic2"

	vm := NewVirtualMachine("my-vm", "eastus", "Standard_D2s_v3").
		WithNetworkInterface(nic1ID, true).
		WithNetworkInterface(nic2ID, false)

	require.Len(t, vm.Properties.NetworkProfile.NetworkInterfaces, 2)
	assert.Equal(t, nic1ID, vm.Properties.NetworkProfile.NetworkInterfaces[0].ID)
	assert.True(t, *vm.Properties.NetworkProfile.NetworkInterfaces[0].Primary)
	assert.Equal(t, nic2ID, vm.Properties.NetworkProfile.NetworkInterfaces[1].ID)
	assert.False(t, *vm.Properties.NetworkProfile.NetworkInterfaces[1].Primary)
}

func TestVirtualMachine_ChainedBuilders(t *testing.T) {
	nicID := "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/my-nic"
	vm := NewVirtualMachine("my-vm", "eastus", "Standard_D2s_v3").
		WithTags(map[string]string{"env": "prod"}).
		WithImage("Canonical", "UbuntuServer", "18.04-LTS", "latest").
		WithNetworkInterface(nicID, true)

	assert.Equal(t, "my-vm", vm.Name)
	assert.Equal(t, "prod", vm.Tags["env"])
	assert.Equal(t, "Canonical", *vm.Properties.StorageProfile.ImageReference.Publisher)
	assert.Equal(t, nicID, vm.Properties.NetworkProfile.NetworkInterfaces[0].ID)
}

func TestVirtualMachine_JSON(t *testing.T) {
	vm := NewVirtualMachine("my-vm", "eastus", "Standard_D2s_v3").
		WithTags(map[string]string{"env": "prod"})

	data, err := json.Marshal(vm)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "my-vm", result["name"])
	assert.Equal(t, "Microsoft.Compute/virtualMachines", result["type"])
	assert.Equal(t, "2021-07-01", result["apiVersion"])
	assert.Equal(t, "eastus", result["location"])

	tags := result["tags"].(map[string]interface{})
	assert.Equal(t, "prod", tags["env"])

	props := result["properties"].(map[string]interface{})
	hw := props["hardwareProfile"].(map[string]interface{})
	assert.Equal(t, "Standard_D2s_v3", hw["vmSize"])
}

func TestVirtualMachine_JSONWithImage(t *testing.T) {
	vm := NewVirtualMachine("my-vm", "eastus", "Standard_D2s_v3").
		WithImage("Canonical", "UbuntuServer", "18.04-LTS", "latest")

	data, err := json.Marshal(vm)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	props := result["properties"].(map[string]interface{})
	storage := props["storageProfile"].(map[string]interface{})
	imageRef := storage["imageReference"].(map[string]interface{})

	assert.Equal(t, "Canonical", imageRef["publisher"])
	assert.Equal(t, "UbuntuServer", imageRef["offer"])
	assert.Equal(t, "18.04-LTS", imageRef["sku"])
	assert.Equal(t, "latest", imageRef["version"])
}

func TestHardwareProfile(t *testing.T) {
	hp := HardwareProfile{
		VMSize: "Standard_D4s_v3",
	}

	data, err := json.Marshal(hp)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "Standard_D4s_v3", result["vmSize"])
}

func TestStorageProfile(t *testing.T) {
	publisher := "Canonical"
	offer := "UbuntuServer"
	sku := "18.04-LTS"
	version := "latest"

	sp := StorageProfile{
		ImageReference: &ImageReference{
			Publisher: &publisher,
			Offer:     &offer,
			SKU:       &sku,
			Version:   &version,
		},
		OSDisk: OSDisk{
			CreateOption: "FromImage",
		},
	}

	data, err := json.Marshal(sp)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	imageRef := result["imageReference"].(map[string]interface{})
	assert.Equal(t, "Canonical", imageRef["publisher"])

	osDisk := result["osDisk"].(map[string]interface{})
	assert.Equal(t, "FromImage", osDisk["createOption"])
}

func TestOSProfile(t *testing.T) {
	computerName := "my-vm"
	adminUser := "azureuser"
	adminPass := "P@ssw0rd123!"

	profile := OSProfile{
		ComputerName:  &computerName,
		AdminUsername: &adminUser,
		AdminPassword: &adminPass,
	}

	data, err := json.Marshal(profile)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "my-vm", result["computerName"])
	assert.Equal(t, "azureuser", result["adminUsername"])
	assert.Equal(t, "P@ssw0rd123!", result["adminPassword"])
}

func TestLinuxConfiguration(t *testing.T) {
	disablePass := true
	provisionAgent := true
	path := "/home/azureuser/.ssh/authorized_keys"
	keyData := "ssh-rsa AAAAB..."

	config := LinuxConfiguration{
		DisablePasswordAuthentication: &disablePass,
		ProvisionVMAgent:              &provisionAgent,
		SSH: &SSHConfiguration{
			PublicKeys: []SSHPublicKey{
				{Path: &path, KeyData: &keyData},
			},
		},
	}

	data, err := json.Marshal(config)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, true, result["disablePasswordAuthentication"])
	assert.Equal(t, true, result["provisionVMAgent"])

	ssh := result["ssh"].(map[string]interface{})
	keys := ssh["publicKeys"].([]interface{})
	require.Len(t, keys, 1)

	key := keys[0].(map[string]interface{})
	assert.Equal(t, "/home/azureuser/.ssh/authorized_keys", key["path"])
}

func TestNetworkProfile(t *testing.T) {
	primary := true
	np := NetworkProfile{
		NetworkInterfaces: []NetworkInterfaceReference{
			{
				ID:      "/subscriptions/xxx/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/nic1",
				Primary: &primary,
			},
		},
	}

	data, err := json.Marshal(np)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	nics := result["networkInterfaces"].([]interface{})
	require.Len(t, nics, 1)

	nic := nics[0].(map[string]interface{})
	assert.Contains(t, nic["id"], "nic1")
	assert.Equal(t, true, nic["primary"])
}

func TestIdentity(t *testing.T) {
	id := Identity{
		Type: "SystemAssigned",
	}

	data, err := json.Marshal(id)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "SystemAssigned", result["type"])
}

func TestBootDiagnostics(t *testing.T) {
	enabled := true
	storageUri := "https://mystorageacct.blob.core.windows.net/"

	bd := BootDiagnostics{
		Enabled:    &enabled,
		StorageURI: &storageUri,
	}

	data, err := json.Marshal(bd)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, true, result["enabled"])
	assert.Equal(t, "https://mystorageacct.blob.core.windows.net/", result["storageUri"])
}

func TestDataDisk(t *testing.T) {
	diskName := "data-disk-1"
	dd := DataDisk{
		Name:         &diskName,
		Lun:          0,
		CreateOption: "Empty",
	}

	data, err := json.Marshal(dd)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "data-disk-1", result["name"])
	assert.Equal(t, float64(0), result["lun"])
	assert.Equal(t, "Empty", result["createOption"])
}
