---
title: "K3s on Civo with Terraform"
created: 2025-07-23
updated: 2025-07-23
tags: "civo, terraform, kubernetes, gitops"
---
Let's go through how I've setup easy deployment of containerized applications using [Civo's K3s](https://www.civo.com/kubernetes), Helm, ArgoCD, all built using Terraform.
Git repository: [alex-rosta/civo-env](https://github.com/alex-rosta/civo-env)
![image](https://www.milesweb.in/blog/wp-content/uploads/2023/01/kubernetes-vs-terraform-pros-cons-and-their-differences.png)

## Philosophy

- **GitOps**: All infrastructure and application deployments are managed through Git, ensuring version control and traceability.
- **Simplicity**: Lightweight K3s.
- **Secrets offloaded**: Nearly all are called from Akeyless during deployment.
- **Split deployment for infrastructure and applications**: Infrastructure is deployed first, followed by the applications layer.

### Secrets and State file

The state file is stored in Cloudflare R2 using this provider:

Infrastructure-layer:

```hcl
terraform {
  required_providers {
    civo = {
      source = "civo/civo"
    }
    akeyless = {
      source = "akeyless-community/akeyless"
    }
  }
  backend "s3" {
    endpoints = {
      s3 = var.s3_endpoint
    }
    key                         = "infra/terraform.tfstate"
    bucket                      = var.s3_bucket_name
    region                      = "auto"
    access_key                  = var.cf_access_key
    secret_key                  = var.cf_secret_key
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true
    use_path_style              = true
  }
}

provider "akeyless" {
  api_gateway_address = "https://api.akeyless.io"
  api_key_login {
    access_id  = var.akeyless_access_id
    access_key = var.akeyless_access_key
  }
}

provider "civo" {
  region = var.region
}
```

Application-layer:

```hcl
terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
    cloudflare = {
      source = "cloudflare/cloudflare"
    }
    akeyless = {
      source = "akeyless-community/akeyless"
    }
  }
  backend "s3" {
    endpoints = {
      s3 = var.s3_endpoint
    }
    key                         = "k3s/terraform.tfstate"
    bucket                      = var.s3_bucket_name
    region                      = "auto"
    access_key                  = var.cf_access_key
    secret_key                  = var.cf_secret_key
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true
    use_path_style              = true
  }
}

provider "akeyless" {
  api_gateway_address = "https://api.akeyless.io"
  api_key_login {
    access_id  = var.akeyless_access_id
    access_key = var.akeyless_access_key
  }
}

provider "kubernetes" {
  config_path = "../infra/kubeconfig"
}
```

So, the secrets needed for the initial setup (terraform.tfvars) are as follows:

```hcl
cf_access_key       = "your_cloudflare_access_key"
cf_secret_key       = "your_cloudflare_secret_key"
s3_bucket_name      = "your_s3_bucket_name"
s3_endpoint         = "your_s3_endpoint"
akeyless_access_id  = "your_akeyless_access_id"
akeyless_access_key = "your_akeyless_access_key"
```

Rest is gathered during deployment.

## Infrastructure Setup `infra/main.tf`

The infrastructure setup is done in the `infra` directory, which includes:

- **Firewall**: your outbound IP is grabbed and whitelisted to run `kubectl` commands towards the K3s instance. All else are blocked except 443.
- **K3s Cluster**: modularized in `modules/cluster`.
- **kubeconfig**: generated and stored for later use with the Kubernetes tf provider.
- **K3s apps**: Civo provides a set of apps that can be installed, of these we install argo-cd, cert-manager, metrics-server and nginx ingress-controller.
- **Nginx Ingress**: single IP routed for all applications, with automatic TLS using cert-manager, for each application.

## Kubernetes Setup `k3s/main.tf`

The Kubernetes setup is done in the `k3s` directory, which includes:

- **cert-manager configuration**: modularized in `modules/cert_manager` referencing nginx as ingress-class (HTTP-01 Solver).
- **Nginx Ingress**: modularized in `modules/ingress` with configuration for automatic TLS.
- **Secrets->Akeyless**: pulled from Akeyless and used in the applications with `kubernetes_secret` resources.
- **DNS**: modularized in `modules/dns` to create DNS records for each application using Cloudflare. All applications will be pointed to the same IP address of the Nginx Ingress controller.

### The fun part

Now let's see how easy it is to deploy applications using this method.
In `gitops/argocd` we have ArgoCD yaml deployment files for each application. Which defines, esentially a helm-chart:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: armory
  namespace: argocd
spec:
  destination:
    server: https://kubernetes.default.svc
    namespace: armory
  project: default
  source:
    repoURL: https://github.com/alex-rosta/armory-helm.git
    targetRevision: HEAD
    path: .
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

Now, that helm-chart contains the actual configuration. Lets check that out [Helm Repo](https://github.com/alex-rosta/armory-helm)
Ingress and autoscaling are configured in the `values.yaml` file;

```yaml
ingress:
  enabled: true
  className: "nginx"
  annotations: 
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/backend-protocol: "HTTP"
  hosts:
    - host: armory.rosta.dev
      paths:
        - path: /
          pathType: Prefix
  tls:
    - hosts:
        - armory.rosta.dev
      secretName: armory-tls
```

autoscaling;

```yaml
autoscaling:
  enabled: enabled
  minReplicas: 4
  maxReplicas: 12
  targetCPUUtilizationPercentage: 80
  targetMemoryUtilizationPercentage: 80
```

redis

```yaml
redis:
  enabled: true
  image: redis:latest
  port: 6379
  hostname: redis
```

etc, etc. This can go on forever but you get the idea. Each app looks different :smiley:

### Now, let's go back to `k3s/main.tf`

This is all that's really needed for getting the app deployed and running:

```hcl
module "armory_dns" {
  source             = "../modules/dns"
  cloudflare_email   = local.cloudflare_secrets.cloudflare_email
  cloudflare_api_key = local.cloudflare_secrets.cloudflare_api_key
  cloudflare_zone_id = local.cloudflare_secrets.cloudflare_zone_id
  content            = data.kubernetes_service.nginx_ingress.status[0].load_balancer[0].ingress[0].ip
  name               = "armory.${local.cloudflare_secrets.domain}"
}

resource "kubernetes_manifest" "app-armory" {
  manifest = yamldecode(file("../gitops/argocd/app-armory.yaml"))
}
```

optional secrets

```hcl
resource "kubernetes_secret" "app-armory-secret" {
  metadata {
    name      = "armory-secret"
    namespace = "armory"
  }
  data = {
    "CLIENT_ID"              = local.armory_secrets.client_id
    "CLIENT_SECRET"          = local.armory_secrets.client_secret
    "WARCRAFTLOGS_API_TOKEN" = local.armory_secrets.warcraftlogs_token
    "REDIS_ADDR"             = local.armory_secrets.redis_addr
    "REDIS_PASSWORD"         = local.armory_secrets.redis_password
    "REDIS_DB"               = local.armory_secrets.redis_db
  }
  depends_on = [kubernetes_manifest.app-armory]

}
```

That was all, rinse and repeat for additional applications.
Reach out for any questions or feedback! :heart:
