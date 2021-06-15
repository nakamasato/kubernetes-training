package main


test_nodeSelector_deployment {
  deny["Deployment must have nodeSelector with nodegroup as a key and prod as value"] with input as
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

test_nodeSelector_deployment {
  deny["Deployment must have nodeSelector with nodegroup as a key and prod as value"] with input as
  {
    "kind": "Deployment",
    "spec": {
      "template": {
        "spec": {
          "containers": [
          ],
          "nodeSelector": {
              "test": "prod"
          }
        }
      }
    }
  }
}

test_nodeSelector_deployment {
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

test_nodeSelector_deployment {
  deny["Deployment must have nodeSelector with nodegroup as a key and prod as value"] with input as
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
