
import sys, threading

from ruamel.yaml import YAML
import json
import lxml

import ZODB, ZODB.FileStorage, BTrees.OOBTree, transaction

import ncclient.manager
import requests

yaml=YAML()

threadlocals = threading.local()

args = yaml.load(sys.stdin) or {}

def from_xml(xml):
    """
    Returns an lxml.Element
    """
    if isinstance(xml, str):
        xml = bytes(xml, encoding="utf-8")
    return lxml.etree.fromstring(xml)

def to_xml(element, pretty=True):
    return lxml.etree.tostring(element, pretty_print=pretty).decode("utf-8")

#
# Store
#

class Store:
    def __init__(self):
        self._transaction_manager = transaction.TransactionManager()

        storage = ZODB.FileStorage.FileStorage("/store/candice.fs")
        db = ZODB.DB(storage)
        connection = db.open(self._transaction_manager)

        self.root = connection.root 

    def get_dict(self, *names):
        dict_ = self.root
        for name in names:
            try:
                return getattr(dict_, name)
            except AttributeError:
                setattr(dict_, name, BTrees.OOBTree.BTree())
                return getattr(dict_, name)
        return dict_

    def commit(self):
        self._transaction_manager.commit()

    def abort(self):
        self._transaction_manager.abort()

#
# Task
#

class Task:
    def __init__(self):
        self.name = args["task"]["name"]
        self.input = args["task"]["input"]
        self.component = args["component"]["name"]
        self.output = {}
        self.store = None

        self._store = None

    @staticmethod
    def current():
        try:
            return threadlocals.task
        except AttributeError:
            raise Exception("not in task")

    def error(self, error):
        """
        'error' can be a string or a YAML-serializable data structure
        """
        if not isinstance(self.output, dict):
            if self.output is not None:
                self.output = {"output": self.output}
            else:
                self.output = {}
        if isinstance(error, str):
            self.output["error"] = {"message": error}
        else:
            self.output["error"] = error

    def abort(self):
        if self._store is None:
            raise Exception("not in 'with'")
        self._store.abort()
        self._store = None

    def __enter__(self):
        threadlocals.task = self
        self._store = Store()
        self.store = self._store.get_dict(self.component, self.name)
        return self

    def __exit__(self, type, value, traceback):
        if self._store is not None:
            self._store.commit()
            self._store = None
            self.store = None
        if (self.output is not None) and (not isinstance(self.output, dict)):
            # Make sure the the output is always a dict
            self.output = {"output": self.output}
        if self.output:
            yaml.dump(self.output, sys.stdout)

#
# Device
#

class Device:
    def __init__(self, name=None):
        """
        Will default to use the component's name
        """
        self.task = Task.current()
        self.name = name or self.task.component
        self.device = args["devices"][self.name]

    def executor(self, type_=None):
        """
        Will device to use the device's protocol
        """
        if type_ is None:
            type_ = self.device["protocol"]
        if type_ == "netconf":
            return NetconfExecutor(self)
        elif type_ == "restconf":
            return RestconfExecutor(self)
        else:
            raise Exception(f"unsupported executor: {type_}")

#
# Executor
#

class Executor:
    def __init__(self, device):
        self.device = device

    def __enter__(self):
        return self

    def __exit__(self, type, value, traceback):
        pass

#
# NetconfExecutor
#

class NetconfExecutor(Executor):
    def __init__(self, device):
        super(NetconfExecutor, self).__init__(device)
        host, port = device.device["direct"]["host"].split(":", 2) 
        self.manager = ncclient.manager.connect_ssh(host=host,
                                                    port=port,
                                                    hostkey_verify=False,
                                                    username="root",
                                                    password="root")
        self.namespaces = {
            "base": "urn:ietf:params:xml:ns:netconf:base:1.0"
        }

    def element(self, name, ns=None, parent=None):
        if parent is None:
            if ns is None:
                return lxml.etree.Element(name)
            else:
                ns = self.namespaces.get(ns, ns)
                return lxml.etree.Element(f"{{{ns}}}{name}")
        else:
            if ns is None:
                if parent.prefix:
                    ns = parent.nsmap[parent.prefix]
                    return lxml.etree.SubElement(parent, f"{{{ns}}}{name}")
                else:
                    return lxml.etree.SubElement(parent, name)
            else:
                ns = self.namespaces.get(ns, ns)
                return lxml.etree.SubElement(parent, f"{{{ns}}}{name}")

    def get_xpath_raw(self, expression):
        # Note: xpath is implemented by the NETCONF agent
        return self.manager.get(filter=("xpath", (self.namespaces, expression))).data_xml

    def get_xpath(self, expression):
        # Note: get_xpath_raw is implemented by the NETCONF agent
        # Sysrepo's implementation is lax about namespaces, allowing child tags to inherit from
        # the parent, but lxml does require all tags to be properly namespaced,
        # thus you do want to want to namespace all tags for this method to work properly
        xml = self.get_xpath_raw(expression)
        root = from_xml(xml)
        return root.xpath("/base:data" + expression, namespaces=self.namespaces)

    def __enter__(self):
        self.manager.__enter__()
        return self

    def __exit__(self, type, value, traceback):
        self.manager.__exit__(type, value, traceback)

#
# RestconfExecutor
#

class RestconfExecutor(Executor):
    def __init__(self, device):
        super(RestconfExecutor, self).__init__(device)
        host = device.device["direct"]["host"]
        self.url = f"http://{host}/restconf/data" 
        self.headers = {
            "Accept" : "application/yang-data+json", 
            "Content-Type" : "application/yang-data+json"
        }

    def get(self, module):
        response = requests.get(f"{self.url}/{module}", headers=self.headers)
        if not response.ok:
            raise RestconfError(response)
        try:
            return response.json()
        except json.decoder.JSONDecodeError:
            return response.text

    def put(self, module, data):
        response = requests.put(f"{self.url}/{module}", headers=self.headers, data=json.dumps(data))
        if not response.ok:
            raise RestconfError(response)
        try:
            return response.json()
        except json.decoder.JSONDecodeError:
            return response.text

    def delete(self, module):
        response = requests.delete(f"{self.url}/{module}", headers=self.headers)
        if not response.ok:
            raise RestconfError(response)
        try:
            return response.json()
        except json.decoder.JSONDecodeError:
            return response.text

    def post(self, module, data):
        response = requests.post(f"{self.url}/{module}", headers=self.headers, data=json.dumps(data))
        if not response.ok:
            raise RestconfError(response)
        try:
            return response.json()
        except json.decoder.JSONDecodeError:
            return response.text

class RestconfError(Exception):
    def __init__(self, response):
        try:
            error = response.json()["ietf-restconf:errors"]["error"]
            message = error["error-message"]
            super(RestconfError, self).__init__(message)
            self.error = error
        except:
            super(RestconfError, self).__init__(response.text)
