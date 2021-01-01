#!/usr/bin/env python3

import candice

from ncclient import manager
#from ncclient.devices.default import DefaultDeviceHandler

host = "router.mynamespace.svc.cluster.local"
port = 830
username = "root"
password = "root"

output = {"caps": []}

# manager = candice.get_executor("router")

with manager.connect_ssh(host=host, port=port, hostkey_verify=False, username=username, password=password) as m:
    for c in m.server_capabilities:
        output["caps"].append(c)

    #c = m.get_schema(identifier="ietf-interfaces")
    #output["schema"] = c.data

    #expr = "/ietf-interfaces:interfaces/"
    #c = m.get_config(source="running", filter=("xpath", expr)).data_xml
    #output["x"] = c

    #c = m.get_config(source="running").data_xml
    #with open("%s.xml" % host, 'w') as f:
    #    f.write(c)

candice.output(output)
