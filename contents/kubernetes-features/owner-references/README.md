# Owner References

## Example1

1. Prepare CRD `Person` and `Pet`

    <details><summary>Person</summary>

    ```yaml
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      name: persons.example.com
    spec:
      group: example.com
      names:
        kind: Person
        listKind: PersonList
        plural: persons
        singular: person
      scope: Namespaced
      versions:
        - name: v1
          served: true
          storage: true
          schema:
            openAPIV3Schema:
              type: object
              properties:
                spec:
                  type: object
                  properties:
                    name:
                      type: string
                    age:
                      type: integer
                    address:
                      type: string
                    email:
                      type: string
                    phone:
                      type: string
                    hobbies:
                      type: array
                      items:
                        type: string
                  required:
                    - name
    ```

    </details>

    <details><summary>Pet</summary>

    ```yaml
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      name: pets.example.com
    spec:
      group: example.com
      versions:
        - name: v1
          served: true
          storage: true
          schema:
            openAPIV3Schema:
              type: object
              properties:
                spec:
                  type: object
                  properties:
                    name:
                      type: string
                    age:
                      type: integer
                    breed:
                      type: string
                    owner:
                      type: string
                status:
                  type: object
                  properties:
                    message:
                      type: string
      scope: Namespaced
      names:
        plural: pets
        singular: pet
        kind: Pet
        listKind: PetList
    ```

    </details>

1. Install the CRD

    ```
    kubectl apply -f crd/person.yaml,crd/pet.yaml
    ```

1. Create `Person` named `Alice`
    ```
    kubectl apply -f example/person-alice.yaml
    ```
1. Get Alice's uid and update `ownerReferences` in  `example/pet-dog.yaml`
    ```
    uid=$(kubectl get -f example/person-alice.yaml -o yaml | yq .metadata.uid); echo $uid
    ```

    ```
    yq e -i ".metadata.ownerReferences[0].uid = \"$uid\"" example/pet-dog.yaml
    ```
1. Create a `Pet` `dog`
    ```
    kubectl apply -f example/pet-dog.yaml
    ```
1. Get Person and Dog
    ```
    kubectl get person,pet
    NAME                       AGE
    person.example.com/alice   5m25s

    NAME                  AGE
    pet.example.com/dog   8s
    ```
1. Delete alice
    ```
    kubectl delete person alice
    person.example.com "alice" deleted
    ```
1. The owned resource `dog` is also deleted. (cascading deletion)
    ```
    kubectl get person,pet
    No resources found in default namespace.
    ```

## Example2: OwnerReferences & Finalizer

1. Create `Person` named `Alice` (same as above)
    ```
    kubectl apply -f example/person-alice.yaml
    ```
1. Get Alice's uid and update `ownerReferences` in  `example/pet-dog.yaml`
    ```
    uid=$(kubectl get -f example/person-alice.yaml -o yaml | yq .metadata.uid); echo $uid
    ```

    ```
    yq e -i ".metadata.ownerReferences[0].uid = \"$uid\"" example/pet-dog-with-finalizer.yaml
    ```
1. Create a `Pet` `dog`
    ```
    kubectl apply -f example/pet-dog-with-finalizer.yaml
    ```
1. Delete alice
    ```
    kubectl delete person alice
    person.example.com "alice" deleted
    ```
1. You can check alice is deleted, while `Pet` still exists as it has finalizer.
    ```
    kubectl get person alice
    Error from server (NotFound): persons.example.com "alice" not found
    ```

    ```
    kubectl get pet
    NAME   AGE
    dog    43s
    ```
1. To clean up the service, `finalizers` must be removed.
    ```
    kubectl patch pet dog -p '{"metadata":{"finalizers": []}}' --type=merge
    ```

**Garbage Collection does not wait to delete the owner object until the dependent object is actually deleted.**


## Example3: OwnerReferences + BlockOwnerDeletion=true & Finalizer

1. Create `Person` named `Alice` (same as above)
    ```
    kubectl apply -f example/person-alice.yaml
    ```
1. Get Alice's uid and update `ownerReferences` in  `example/pet-dog.yaml`
    ```
    uid=$(kubectl get -f example/person-alice.yaml -o yaml | yq .metadata.uid); echo $uid
    ```

    ```
    yq e -i ".metadata.ownerReferences[0].uid = \"$uid\"" example/pet-dog-with-finalizer-and-blockownerdeletion.yaml
    ```
1. Create a `Pet` `dog`
    ```
    kubectl apply -f example/pet-dog-with-finalizer-and-blockownerdeletion.yaml
    ```
1. Delete alice
    ```
    kubectl delete person alice
    person.example.com "alice" deleted
    ```
1. You can check alice is deleted, while `Pet` still exists as it has finalizer.
    ```
    kubectl get person alice
    Error from server (NotFound): persons.example.com "alice" not found
    ```

    ```
    kubectl get pet
    NAME   AGE
    dog    43s
    ```
1. To clean up the service, `finalizers` must be removed.
    ```
    kubectl patch pet dog -p '{"metadata":{"finalizers": []}}' --type=merge
    ```

**Garbage Collection does not wait to delete the owner object until the dependent object is actually deleted.**

## Example4: OwnerReferences + BlockOwnerDeletion=true & Finalizer + delete --cascade=foreground

1. Create `Person` named `Alice` (same as above)
    ```
    kubectl apply -f example/person-alice.yaml
    ```
1. Get Alice's uid and update `ownerReferences` in  `example/pet-dog.yaml`
    ```
    uid=$(kubectl get -f example/person-alice.yaml -o yaml | yq .metadata.uid); echo $uid
    ```

    ```
    yq e -i ".metadata.ownerReferences[0].uid = \"$uid\"" example/pet-dog-with-finalizer-and-blockownerdeletion.yaml
    ```
1. Create a `Pet` `dog`
    ```
    kubectl apply -f example/pet-dog-with-finalizer-and-blockownerdeletion.yaml
    ```
1. Delete alice with `--cascade=foreground`
    ```
    kubectl delete person alice --cascade=foreground
    ```
    The command gets stuck as it waits until all the dependent objects are removed.
1. Remove `finalizers` manually in another terminal.
    ```
    kubectl patch pet dog -p '{"metadata":{"finalizers": []}}' --type=merge
    ```

    Once the finalizer is removed, the command above completes the deletion of the owner object. The pet is also deleted when the finalizer is removed.

**the owner object is deleted after the dependent object deletion completed.**

## Example5: OwnerReferences + BlockOwnerDeletion=true & Finalizer + CRD with propagationPolicy: Background

1. Install the CRD

    ```
    kubectl apply -f crd/person-deletionpolicy-foreground.yaml,crd/pet.yaml
    ```

## Ref

1. https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
1. https://stackoverflow.com/questions/60153700/foreground-cascading-deletion-not-working-as-documentation-suggests
