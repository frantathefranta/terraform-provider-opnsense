// Configure an OSPF interface
resource "opnsense_quagga_ospf_interface" "example0" {
  enabled = true
  interfacename = "opt1"
  authtype = ""
  authkey = ""
  authkey_id = 1
  area = ""
  cost = ""
  cost_demoted = 65535
  hellointerval = 30
  deadinterval = 60
  retransmitinterval = 120
  retransmitdelay = 3
  transmitdelay = 1
  priority = 100
  bfd = false
  networktype = "point-to-point"
}
