
import sys

from ruamel.yaml import YAML

yaml=YAML()

def output(o):
    yaml.dump(o, sys.stdout)
