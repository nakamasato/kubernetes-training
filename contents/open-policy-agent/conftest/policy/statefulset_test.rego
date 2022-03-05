package main


test_no_nodeSelector {
  deny["StatefulSet must have nodeSelector"] with input as
  {
    "kind": "StatefulSet",
    "spec": {
      "template": {
        "spec": {
          "containers": [
          ],
        }
      }
    }
  }
}

test_nodeSelector_with_invalid_key {
  deny["StatefulSet must have nodeSelector with nodegroup as a key"] with input as
  {
    "kind": "StatefulSet",
    "spec": {
      "template": {
        "spec": {
          "containers": [
          ],
          "nodeSelector": {
              "test": "dev"
          }
        }
      }
    }
  }
}

test_nodeSelector_with_invalid_value {
  deny["StatefulSet must have nodeSelector with nodegroup as a key and prod or dev as value"] with input as
  {
    "kind": "StatefulSet",
    "spec": {
      "template": {
        "spec": {
          "containers": [
          ],
          "nodeSelector": {
              "nodegroup": "not_exist"
          }
        }
      }
    }
  }
}


test_nodeSelector_with_valid_pair_prod {
  deny with input as
  {
    "kind": "StatefulSet",
    "spec": {
      "template": {
        "spec": {
          "containers": [
          ],
          "nodeSelector": {
              "nodegroup": "prod"
          }
        }
      }
    }
  }
}

test_nodeSelector_with_valid_pair_dev {
  deny with input as
  {
    "kind": "StatefulSet",
    "spec": {
      "template": {
        "spec": {
          "containers": [
          ],
          "nodeSelector": {
              "nodegroup": "dev"
          }
        }
      }
    }
  }
}

test_nodeSelector_with_invalid_value_for_dev_namespace {
  deny["nodegroup must be dev in non-prod namespace"] with input as
  {
    "kind": "StatefulSet",
    "metadata": {
      "namespace": "dev"
    },
    "spec": {
      "template": {
        "spec": {
          "containers": [
          ],
          "nodeSelector": {
              "nodegroup": "prod"
          }
        }
      }
    }
  }
}

test_nodeSelector_with_invalid_value_for_dev_namespace {
  deny["nodegroup must be prod in prod namespace"] with input as
  {
    "kind": "StatefulSet",
    "metadata": {
      "namespace": "prod"
    },
    "spec": {
      "template": {
        "spec": {
          "containers": [
          ],
          "nodeSelector": {
              "nodegroup": "dev"
          }
        }
      }
    }
  }
}
