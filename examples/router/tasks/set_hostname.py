#!/usr/bin/env python3

import candice
import lxml

with candice.Task() as task:
    hostname = task.input.get("hostname", "hello")

    router = candice.Device()

    with router.executor() as executor:
        executor.namespaces.update({"s": "urn:ietf:params:xml:ns:yang:ietf-system"})

        element = executor.get_xpath("/s:system/s:hostname")
        if element:
            task.output["old-hostname"] = element[0].text

        config = executor.element("config", ns="base")
        element = executor.element("system", ns="s", parent=config)
        element = executor.element("hostname", parent=element)
        element.text = hostname

        executor.manager.edit_config(target="running", config=candice.to_xml(config))

        element = executor.get_xpath("/s:system/s:hostname")
        if element:
            task.output["new-hostname"] = element[0].text
