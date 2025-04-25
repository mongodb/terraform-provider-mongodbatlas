# Example - Microsoft Azure and MongoDB Atlas Streams Private Endpoint

This example shows how to use Azure PrivateLink Endpoints with EventHub for Atlas Streams PrivateLink.

You must set the following variables for Atlas in main.tf:

- `public_key`: Public API key to authenticate to Atlas
- `private_key`: Private API key to authenticate to Atlas
- `project_id`: Unique 24-hexadecimal digit string that identifies your atlas project
- `atlas_region`: Atlas region where you want to create the Streams PrivateLink resources. To learn more, see `Atlas Region` column in https://www.mongodb.com/docs/atlas/reference/microsoft-azure/#stream-processing-instances. 

- Additional required fields in main.tf:
- `dns_domain`: Hostname of the Event Hub Namespace in Azure, which is the dns_domain.
- `service_endpoint_id`: Service Endpoint ID for the EventHub Namespace. You can find this in the Azure portal under the EventHub Namespace properties. It typically looks like `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}`.

The following setup is for Azure PrivateLink with EventHub example in azure.tf. To learn more, see documentation https://learn.microsoft.com/en-us/azure/event-hubs/private-link-service#add-a-private-endpoint-using-azure-portal

- `azure_region`: The Azure region where you want to create the Azure PrivateLink resources. `Azure Region` column in https://www.mongodb.com/docs/atlas/reference/microsoft-azure/#stream-processing-instances.
- `azure_resource_group`: The name of the Azure Resource Group where you want to create the PrivateLink resources. 
- `vnet_name`: The name of the Azure Virtual Network (VNet) where you want to create the PrivateLink resources.
- `subnet_name`: The name of the subnet within the VNet where you want to create the PrivateLink resources. 
- `eventhub_namespace_name`: The name of the Azure EventHub Namespace that you want to use for the PrivateLink connection. Must be globally unique. 
- `eventhub_name`: The name of the Azure EventHub that you want to use for the PrivateLink connection. 
- `vnet_address_space`: The address space for the Azure Virtual Network. 
- `subnet_address_prefix`: The address prefix for the Azure Subnet.

## Usage

**1\. Ensure that your Azure credentials are set up.**

1. Install the Azure CLI by following the steps from the [official Azure documentation](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli).
2. Run the command `az login` and authenticate using the default browser.
3. Once authenticated, Azure returns the following user details:

**2\. Set up your MongoDB Atlas API keys.**
1. Log in to your MongoDB Atlas account.
2. Navigate to the "Project Access" section of your project.
3. Create a new API key with the necessary permissions (Project Owner or similar).
4. Copy the Public and Private keys to use with the Terraform configuration.

**3\. Create a terraform.tfvars file.**
1. Create a file named `terraform.tfvars` in the same directory as your `main.tf`.
2. Defining the required variables in the `terraform.tfvars` file.

**4\. Optional: Create the Azure resources with EventHub.**
1. If you don't have an existing Azure EventHub Namespace and EventHub, you can create them using the provided `azure.tf` file.

**5\. Initialize Terraform.**
1. Run the following command to initialize Terraform and download the required providers:
   ```bash
   terraform init
   ```
**6\. Plan the Terraform deployment.**
1. Run the following command to see the execution plan and verify the resources that will be created:
   ```bash
   terraform plan
   ```
   
**7\. Apply the Terraform configuration.**
1. If the plan looks good, run the following command to create the resources:
   ```bash
   terraform apply
   ```

**8\. Destroy the Terraform resources.**
1. When you no longer need the resources, you can destroy them by running the following command:
   ```bash
   terraform destroy
   ```

