---
title: "Infrastructure and Technologies I Use Personally"
tags: "welcome, personal, it-infrastructure"
created: 2025-07-21
---
In this post, I want to share the various technologies and infrastructure I rely on in my personal projects, usually with generous free tiers or open-source options.

## Cloud Providers

### Cloudflare

- **Usage**: DNS management, CDN, and security.
- **Note**: R2, Workers, Pages, Proxies, and geoblocking. It's really a Swiss Army knife for web infrastructure. :shield:

### Civo

- **Usage**: Kubernetes clusters and managed databases.
- **Note**: Focused on simplicity and ease of use for developers. Using it primarily for K3S with their Terraform provider. Check out this [blog post](https://www.civo.com/learn/terraform-kubernetes-cluster) for my implementation. :computer:

### Akeyless

- **Usage**: Secrets management and encryption.
- **Note**: Retrieve your secrets securely using API calls, works great in your pipelines. :closed_lock_with_key:

### Fly.io

- **Usage**: Application hosting and serverless functions.
- **Note**: Great for small apps and services, especially with their free tier. I use it for hosting some of my personal projects. :rocket:

### Proton

- **Usage**: Email hosting :email:
