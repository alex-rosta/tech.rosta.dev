---
title: "Host this Blog Yourself"
tags: "go, web development, tutorial, backend"
created: 2025-07-21
updated: 2025-07-21
---

## Host this Blog Yourself

In this tutorial, I'll walk you through the steps to host this blog on your own server or PC using Go and Docker. This blog is built with modern web technologies and is designed to be easy to deploy.

### Prerequisites

- **Docker**: Make sure you have Docker installed on your machine.
- **Hardware**: Somewhere to deploy to :smiley:
- **Domain**: Optional, but of course recommended.

### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd goblog
```

### Step 2: Build the Docker Image

```bash
docker build -t goblog .
```

### Step 3: Run the Docker Container

```bash
docker run -p 3000:3000 goblog
```

### Step 4: Access the Blog

Open your web browser and go to `http://localhost:3000` to see the blog in action.

### Adding Your Own Content

To add your own posts, simply create Markdown files in the `posts` directory. The blog will automatically pick them up and display them after a restart with `docker restart goblog`.
