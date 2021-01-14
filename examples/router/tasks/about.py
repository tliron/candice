#!/usr/bin/env python3

import candice

with candice.Task() as task:
    if "count" in task.store:
        task.store["count"] += 1
    else:
        task.store["count"] = 1

    task.output["count"] = task.store["count"]

    device = candice.Device()

    with device.executor() as executor:
        task.output["inputs"] = task.input

        task.output["capabilities"] = []
        for c in executor.manager.server_capabilities:
            task.output["capabilities"].append(c)

        config = candice.from_xml(executor.manager.get_config(source="running").data_xml)
        task.output["config"] = candice.to_xml(config)

        executor.namespaces.update({
            "s": "urn:ietf:params:xml:ns:yang:ietf-system",
            "k": "urn:ietf:params:xml:ns:yang:ietf-keystore",
            "i": "urn:ietf:params:xml:ns:yang:ietf-interfaces"
        })

        element = executor.get_xpath("//s:current-datetime")
        task.output["time"] = element[0].text

        element = executor.get_xpath("//s:os-name")
        task.output["os"] = element[0].text

        element = executor.get_xpath("/s:system/s:hostname")
        if element:
            task.output["hostname"] = element[0].text
