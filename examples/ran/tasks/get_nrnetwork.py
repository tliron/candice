#!/usr/bin/env python3

import candice

with candice.Task() as task:
    id = task.input.get("id", 100)
    device = candice.Device()

    with device.executor() as executor:
        response = executor.get(f"_3gpp-nr-nrm-nrnetwork:NRNetwork={id}")
        task.output = response
