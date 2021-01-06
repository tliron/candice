*This is an early release. Some features are not yet fully implemented.*

Candice
=======

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/candice.svg)](https://github.com/tliron/candice/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/candice)](https://goreportcard.com/report/github.com/tliron/candice)

Cloud-Native Network Device Configurator.


Get It
------

[![Download](assets/media/download.png "Download")](https://github.com/tliron/candice/releases)


Rationale
---------

Managing network services requires configuring their network functions. Whether those reside in
clouds (CNFs and VNFs) or in boxes (PNFs) they may require sending the configuration data via
NETCONF or similar protocols.

This can be subtle work, as much of this data is contextual and conditional, e.g. you must query for
what network interfaces are available and what features they have before setting them up to achieve
the required functionality. Often this back-and-forth is handled by a controller or orchestrator
where you can design "workflows", "playbooks", "recipes", etc., in the context of the network
service more generally.

The problem is that as our deployments grow into thousands and tens of thousands of sites, with
hundreds of thousands of network functions, these big and faraway orchestrators struggle to handle
the scale.

Candice can make this easier by delegating the function's configuration work to its site, indeed it
should run within CNF's cluster. Work is encapsulated into goal-oriented "tasks", which can be
understood as "micro-workflows" within the general orchestration context. Each task can make many
different NETCONF calls, even on several devices, and can rely on relevant state managed locally
instead of at the faraway orchestrator. The orchestrator can then just tell Candice to run the
whole task with custom inputs and outputs.


How It Works
------------

Candice is a Kubernetes operator that:

1. lets you to create "Device" custom resources, which contain connectivity information
   (addresses, protocols, credentials) for a device. Currently the NETCONF protocol is supported.
2. lets you attach "tasks" to the operator. A task is a Python program that handles a specific
   configuration workflow on one or more devices. Tasks can accept arbitrary inputs and emit
   arbitrary outputs, allowing you to customize them as necessary for integration with management
   and orchestration solutions, such as MANO.
3. provides a rich Python API to make it easier to write tasks, e.g. setting up NETCONF connectivity
   and allowing you to read and write configuration data. Also included is persistent transactional
   state, a "scratch" area that you can use for whatever state is necessary for your tasks,
   including for tasks to share state with each other.

[Here](examples/router/tasks) are examples of a few tasks.

Usage example:

    candice device create router --host router.mynamespace.svc.cluster.local:830
    
    candice task set router about --file=examples/router/tasks/about.py
    candice task set router set-hostname --file=examples/router/tasks/set_hostname.py
    
    candice task run router set-hostname --input=hostname=myhost
    candice task run router about
