#!/bin/bash

# Multi-Tier Web Application Deployment Script
# This script deploys the ARM template to Azure

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
RESOURCE_GROUP_NAME="webapp-enterprise-rg"
LOCATION="eastus"
DEPLOYMENT_NAME="webapp-deployment-$(date +%Y%m%d-%H%M%S)"

echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Multi-Tier Web Application Deployment                    ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if Azure CLI is installed
if ! command -v az &> /dev/null; then
    echo -e "${RED}❌ Azure CLI is not installed. Please install it first.${NC}"
    echo "   Visit: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli"
    exit 1
fi

echo -e "${GREEN}✓${NC} Azure CLI found"

# Check if logged in
echo -n "Checking Azure login status... "
if ! az account show &> /dev/null; then
    echo -e "${YELLOW}not logged in${NC}"
    echo "Please log in to Azure:"
    az login
else
    echo -e "${GREEN}✓${NC}"
fi

# Get current subscription
SUBSCRIPTION_NAME=$(az account show --query name -o tsv)
SUBSCRIPTION_ID=$(az account show --query id -o tsv)
echo -e "Current subscription: ${YELLOW}${SUBSCRIPTION_NAME}${NC} (${SUBSCRIPTION_ID})"
echo ""

# Prompt for admin password
echo -e "${YELLOW}Please enter a secure password for VM administrator:${NC}"
echo "Requirements: 12+ characters, uppercase, lowercase, number, and special character"
read -s -p "Password: " ADMIN_PASSWORD
echo ""
read -s -p "Confirm password: " ADMIN_PASSWORD_CONFIRM
echo ""

if [ "$ADMIN_PASSWORD" != "$ADMIN_PASSWORD_CONFIRM" ]; then
    echo -e "${RED}❌ Passwords do not match!${NC}"
    exit 1
fi

# Prompt for instance count
echo ""
echo -e "${YELLOW}How many web server instances do you want? (1-10)${NC}"
read -p "Instance count [2]: " INSTANCE_COUNT
INSTANCE_COUNT=${INSTANCE_COUNT:-2}

# Validate instance count
if ! [[ "$INSTANCE_COUNT" =~ ^[0-9]+$ ]] || [ "$INSTANCE_COUNT" -lt 1 ] || [ "$INSTANCE_COUNT" -gt 10 ]; then
    echo -e "${RED}❌ Invalid instance count. Must be between 1 and 10.${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Deployment Configuration                                 ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo -e "  Resource Group: ${YELLOW}${RESOURCE_GROUP_NAME}${NC}"
echo -e "  Location:       ${YELLOW}${LOCATION}${NC}"
echo -e "  VM Instances:   ${YELLOW}${INSTANCE_COUNT}${NC}"
echo -e "  Deployment:     ${YELLOW}${DEPLOYMENT_NAME}${NC}"
echo ""

read -p "Proceed with deployment? (yes/no): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
    echo -e "${RED}Deployment cancelled.${NC}"
    exit 0
fi

echo ""
echo -e "${GREEN}Starting deployment...${NC}"
echo ""

# Create resource group
echo -n "Creating resource group... "
if az group create \
    --name "$RESOURCE_GROUP_NAME" \
    --location "$LOCATION" \
    --output none; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}❌ Failed${NC}"
    exit 1
fi

# Deploy template
echo "Deploying ARM template... (this may take 10-15 minutes)"
if az deployment group create \
    --resource-group "$RESOURCE_GROUP_NAME" \
    --name "$DEPLOYMENT_NAME" \
    --template-file template.json \
    --parameters vmAdminUsername=azureuser \
    --parameters vmAdminPassword="$ADMIN_PASSWORD" \
    --parameters instanceCount="$INSTANCE_COUNT" \
    --output json > deployment-output.json; then

    echo -e "${GREEN}✓ Deployment completed successfully!${NC}"
    echo ""

    # Extract outputs
    PUBLIC_IP=$(jq -r '.properties.outputs.publicIPAddress.value' deployment-output.json)
    PUBLIC_FQDN=$(jq -r '.properties.outputs.publicFQDN.value' deployment-output.json)
    STORAGE_ACCOUNT=$(jq -r '.properties.outputs.storageAccountName.value' deployment-output.json)

    echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║  Deployment Results                                       ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
    echo -e "  Public IP:        ${YELLOW}${PUBLIC_IP}${NC}"
    echo -e "  Public FQDN:      ${YELLOW}${PUBLIC_FQDN}${NC}"
    echo -e "  Storage Account:  ${YELLOW}${STORAGE_ACCOUNT}${NC}"
    echo ""
    echo -e "${GREEN}Access your application:${NC}"
    echo -e "  ${YELLOW}http://${PUBLIC_FQDN}${NC}"
    echo ""
    echo -e "${YELLOW}Note:${NC} It may take a few minutes for the web servers to fully initialize."
    echo ""

    # Save deployment info
    cat > deployment-info.txt <<EOF
Deployment Information
======================
Date: $(date)
Resource Group: ${RESOURCE_GROUP_NAME}
Location: ${LOCATION}
Deployment Name: ${DEPLOYMENT_NAME}

Outputs:
--------
Public IP: ${PUBLIC_IP}
Public FQDN: ${PUBLIC_FQDN}
Storage Account: ${STORAGE_ACCOUNT}

Access URL: http://${PUBLIC_FQDN}

To delete all resources:
  az group delete --name ${RESOURCE_GROUP_NAME} --yes --no-wait
EOF

    echo -e "${GREEN}✓${NC} Deployment information saved to: ${YELLOW}deployment-info.txt${NC}"

else
    echo -e "${RED}❌ Deployment failed${NC}"
    echo "Check deployment-output.json for details"
    exit 1
fi

echo ""
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Next Steps                                                ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo "1. Wait 5-10 minutes for VMs to initialize"
echo "2. Access your application at http://${PUBLIC_FQDN}"
echo "3. Monitor resources in Azure Portal"
echo "4. Configure SSL/TLS for production use"
echo ""
echo -e "${YELLOW}To delete all resources when done:${NC}"
echo "  az group delete --name ${RESOURCE_GROUP_NAME} --yes"
echo ""
