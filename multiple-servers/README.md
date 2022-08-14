# Multiple servers

![Architecture of this solution](./readme-assets/architecture.png)

Create file server to serve static HTML files. Create an API server that serves JSON from a database. Run the API and file server as two separate servers. Try to load the website & see CORS issue. Put apache in front of the file server and the API so they are on a single port and hostname. Learn about how to run services in VMs in the cloud. Replicate this local setup in the cloud on a single VM, with all services running on the same host. Route requests to the service.

Timebox: 10 days

Learning objectives:

- Basic microservices ideas, separating concerns of services
- Configure apache to talk to 2-3 copies of the API server
- Some web security ideas (CORS)
- Reverse proxy configuration, routing on path
- Health checks
- Running applications in the cloud on a raw VM
- Using cloud-hosted services like databases
- Multi-environment configuration
