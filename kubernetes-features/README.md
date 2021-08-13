# Kubernetes Features

## Namespaces

*Kubernetes supports multiple virtual clusters backed by the same physical cluster. These virtual clusters are called namespaces.* [Namespaces](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)

When to use multiple namespaces seems to be arguable.
- [Kubernetes - Namespaces](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
    - *For clusters with a few to tens of users, you should not need to create or think about namespaces at all.* (in the Kubernetes document.)
- [RedHat - Kubernetes Namespaces Demystified - How To Make The Most of Them](https://cloud.redhat.com/blog/kubernetes-namespaces-demystified-how-to-make-the-most-of-them)
    - Do not overload namespaces with multiple workloads that perform unrelated tasks.
    - Users should **create namespaces for a specific application or microservice and all of the application requirements**. Reasons:
        - Simplified recreation of the entire application
        - Fine-grained network management
        - Greater scalability
        - Greater observability
- [Google Cloud - Kubernetes best practices: Organizing with Namespaces](https://cloud.google.com/blog/products/containers-kubernetes/kubernetes-best-practices-organizing-with-namespaces)
    - Small team: 5~10 microservices → **`default` Namespace**
    - Rapidly growing teams: 10+ microservices → each team owns their own microservices -> **Use multiple clusters or namespaces for production and development** or **Each team may choose to have their own namespace**
    - Large: not everyone knows everyone else. → **each team definitely needs its own namespace. Each team might even opt for multiple namespaces to run its development and production environments.**
