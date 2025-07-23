---
title: "Host this Blog"
tags: "go, web development, tutorial, backend, docker"
created: 2025-07-21
updated: 2025-07-23
---
In this tutorial, I'll walk you through the steps to host this blog on your own server or PC using Go and Docker. This blog is built with modern web technologies and is designed to be easy to deploy.

![image](https://perfectmediaserver.com/images/logos/docker-logo-horizontal.png)

### Prerequisites

- **Docker**: Make sure you have Docker installed on your machine.
- **Hardware**: Somewhere to deploy to :smiley:

### Step 1: Clone the Repository

```bash
git clone https://github.com/alex-rosta/tech.rosta.dev.git
cd tech.rosta.dev
```

### Step 2: Build the Docker Image

```bash
docker build -t goblog .
```

### Step 3: Run the Docker Container

```bash
docker run -p 3000:3000 goblog -d
```

### Step 4: Access the Blog

Open your web browser and go to `http://localhost:3000` to see the blog in action.

### Adding Your Own Content

To add your own posts, simply create Markdown files in the `posts` directory. The blog will automatically pick them up and display them after a restart with `docker restart goblog`.
The posts should follow the same format as the existing ones, with a YAML front matter for metadata like title, tags, and creation date and update date if applicable.

### Customizing the Blog

Customizing is as easy as editing the different html files in the `templates` directory. You can change the layout, colors, and styles to fit your personal taste. Make sure to add your own links.
