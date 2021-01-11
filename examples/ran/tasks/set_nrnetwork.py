#!/usr/bin/env python3

import candice

with candice.Task() as task:
    id = task.input.get("id", 22)
    device = candice.Device()

    with device.executor("restconf") as executor:
        response = executor.post(f"_3gpp-nr-nrm-nrnetwork:NRNetwork={id}", {
                                     "_3gpp-nr-nrm-nrnetwork:attributes": {
                                         "dnPrefix": "er",
                                         "priorityLabel": 1
                                     }})
        task.output = response
