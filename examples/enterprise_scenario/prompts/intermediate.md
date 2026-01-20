Create an ARM template for a multi-tier enterprise application:

**Network:**
- Virtual network (10.0.0.0/16) with subnets for web and app tiers
- Network security groups with appropriate rules
- Public IP for load balancing

**Compute:**
- Web tier VMSS (2 instances) in web-subnet
- Application Gateway for internet-facing traffic

**Storage:**
- Storage account for application data

**Security:**
- NSG rules to restrict traffic between tiers

Location: East US

Generate a single ARM template JSON file with all resources.
