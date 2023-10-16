---
name: Bug report
about: Something unexpected happen? Report an error, crash, or an incorrect behavior here.

---

Hello!

Thank you for opening an issue. **These GitHub issues** are only for community reporting and assistance; as such, we don't have a guaranteed SLA.

**If your issue relates to Terraform itself**, please open it in the Terraform repository https://github.com/hashicorp/terraform/issues.

**If you have an active MongoDB Atlas Support contract**, the best way for us to assist you with the Terraform MongoDB Atlas Provider is through a [support ticket](https://support.mongodb.com/).

**Please note:** In order for us to provide the best experience in supporting our customers, we kindly ask to make sure that all the following sections are correctly filled with all the required information. Our support will **prioritise** issues that contain **all the required** information by using the **"one-click reproducible issues" principle**.

**Please also note:** We try to keep the Terraform MongoDB Atlas Provider issue tracker **reserved for bug reports**. Please ensure you **check open and closed issues first** to ensure your issue hasn't already been reported (if it has been reported add a reaction, i.e. +1, to the issue).

### Guidelines

In order to follow the **"one-click reproducible issues" principle**, please follow these guidelines:

* We should be able to make no changes to your provided script and **be able to run a local execution reproducing the issue**.
  * This means that you should kindly provide us the whole script contatining all the required instructions. This also includes but not limited to:
    * Terraform Atlas provider version used to reproduce the issue
    * Terraform version used to reproduce the issue
  * Configurations that cannot be properly executed will be de-prioritised in favour of the ones that succeed.
* Share your configuration by **using variables**. Create a `variables.tf` file and share its content.
* Before opening an issue, you have to try to specifically isolate it to **Terraform MongoDB Atlas** provider by **removing as many dependencies** as possible. If the issue only happens with other dependencies, then:
  * If other terraform providers are required, please make sure you also include those. Same "one-click reproducible issue" principle applies.
  * If external components are required to replicate it, please make sure you also provides instructions on those parts.

* If the issue requires **changes from the UI**, please document them preferably with screenshots
* If the issue involves **version changes** (e.g. `run plan with version 1.k.x then re-run with version 1.k.y`), please make sure you report the detailed sequence on commands to execute.
* Please confirm if the platform being used is Terraform OSS, Terraform Cloud, or Terraform Enterprise deployment


### Terraform CLI and Terraform MongoDB Atlas Provider Version
<!---
Run `terraform version` to show the version, and paste the result for Terraform and the MongoDB Atlas Provider between the ``` marks below.

If you are not running the latest version of Terraform and the MongoDB Atlas Provider, please try to reproduce the issue in a more recent version in case it has been fixed already. 
-->

```hcl
# Copy-paste your version info here
```

### Terraform Configuration File 
<!--
Paste here all the required parts of your Terraform MongoDB Atlas Provider configuration file below. You can alternatively consider creating a [**GitHub Gist**](https://gist.github.com/) with the information and share it here.

Please be sure to redact any sensitive information; common examples include API Keys, passwords, other secrets, or any items your project/business may consider sensitive such as hostnames or usernames.
-->

```hcl
# Copy-paste your configuration info here
```

### Steps to Reproduce
<!--
Please list the full steps required to reproduce the issue, for example:
1. `terraform init`
2. `terraform apply`
-->

### Expected Behavior
<!--
What should have happened?
-->

### Actual Behavior
<!--
What actually happened?
-->

### Debug Output 
<!--
Note: Debug output can be incredibly helpful in narrowing down an issue but is not required.

Full debug output can be obtained by running Terraform with the environment variable `TF_LOG=trace`. Please create either a GitHub Gist or attach a file containing the debug output. Please do _not_ paste the debug output in the issue, since debug output can be very long.

Debug output may contain sensitive information. Please review it before posting publicly, and if you are concerned feel free to redact it.
-->

### Crash Output
<!--
If the console output indicates that Terraform crashed, please either share a link to a GitHub Gist containing the output of the `crash.log` file or attach the file.
-->

### Additional Context
<!--
Are there anything atypical about your situation that we should know? 
-->

### References
<!--
Are there any other related GitHub issues (open or closed) or Pull Requests that should be linked here? 
-->
