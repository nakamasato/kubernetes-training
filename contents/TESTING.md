# Testing (WIP)

Using [KUTTL](https://kuttl.dev/) to test the behaviors of yamls

## Prerequisite

- Docker
- Kind

## Coverage

1. prometheus-operator

## How to add a test

1. Create a directory in `test`. e.g. `prometheus-operator`
1. Write your `TestStep` in the directory.
1. Run in your local
    ```
    kubectl kuttl test --start-kind=false
    ```
