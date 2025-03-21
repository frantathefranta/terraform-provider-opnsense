# OPNsense Terraform Provider

This provider seeks to support the *entire* OPNsense API.
However, this provider does not, and will not, support resources
not currently supported by the OPNsense API. If required, see if
[dalet-oss/opnsense](https://github.com/dalet-oss/terraform-provider-opnsense)
will support your needs.


⚠️ Please note that this provider is under active development, and makes no
guarantee to be stable. For that reason, it is not currently recommended
to use this provider in any production environment. If a feature is missing,
but is documented in the OPNsense API, please raise an issue to indicate interest.

## Fork comments
I've forked this from [browningluke/terraform-provider-opnsense]( https://github.com/browningluke/terraform-provider-opnsense) because I wanted to see if I could add functionality to it. I've managed to add `OSPFInterface` resource that works sort-of well, we'll see if I want to continue with it.
