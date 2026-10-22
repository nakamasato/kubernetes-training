package main


test_no_nodeSelector {
  deny["Deployment must have nodeSelector"] with input as
  {
    "kind": "Deployment",
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
  deny["Deployment must have nodeSelector with nodegroup as a key"] with input as
  {
    "kind": "Deployment",
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
  deny["Deployment must have nodeSelector with nodegroup as a key and prod or dev as value"] with input as
  {
    "kind": "Deployment",
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
    "kind": "Deployment",
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
    "kind": "Deployment",
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
    "kind": "Deployment",
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
    "kind": "Deployment",
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
