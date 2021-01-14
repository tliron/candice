#!/usr/bin/env python3

import candice

with candice.Task() as task:
    router = candice.Device()

    with router.executor() as executor:
        executor.namespaces.update({
            "ietf-system": "urn:ietf:params:xml:ns:yang:ietf-system",
            "ietf-keystore": "urn:ietf:params:xml:ns:yang:ietf-keystore",
            "ietf-interfaces": "urn:ietf:params:xml:ns:yang:ietf-interfaces"
        })

        names = executor.get_xpath("//ietf-keystore:name")
        task.output["names"] = [n.text for n in names]

        i = executor.get_xpath("//ietf-interfaces:interfaces")
        task.output["interfaces"] = i 

        #c = executor.manager.get_schema(identifier="ietf-interfaces")
        #output["schema"] = c.data
