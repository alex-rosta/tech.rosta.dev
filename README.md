# GoBlog - Personal Blog & Technical Documentation

A modern, feature-rich blog application built with Go, designed for sharing technical documentation, tutorials, and personal thoughts.

## 🚀 Features

### ✅ Core Functionality
- **Markdown-powered Posts** - Write posts in Markdown with front matter support
- **Dark/Light Theme Toggle** - Seamless theme switching with localStorage persistence
- **Responsive Design** - Modern, clean design that works on all devices
- **Search Functionality** - Search posts by title, tags, or headers
- **Tag-based Organization** - Browse posts by topics/tags
- **Recent Posts** - Homepage showcases newly updated or created content

### ✅ User Experience
- **Fast Performance** - Built with Go for excellent speed
- **Clean Typography** - Easy-to-read content with proper styling
- **Social Media Integration** - Links to GitHub, Instagram, etc.
- **SEO Friendly** - Proper meta tags and semantic HTML

### ✅ Admin Features
- **GitHub OAuth Authentication** - Secure admin access (structure ready)
- **Post Management** - Create, edit, and delete posts
- **User Whitelist** - Restrict admin access to specified GitHub users

## 🛠 Technology Stack

- **Backend**: Go 1.21 with Chi router
- **Frontend**: HTML templates with Tailwind CSS
- **Markdown Processing**: Goldmark with extensions
- **Storage**: Filesystem-based (with interface for easy extension)
- **Deployment**: Docker + Fly.io + GitHub Actions
- **Authentication**: GitHub OAuth (structure ready)

## 📁 Project Structure

```
goblog/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/                    # GitHub OAuth handling (ready for implementation)
│   ├── config/                  # Configuration management
│   │   └── config.go
│   ├── handlers/                # HTTP handlers
│   │   ├── admin.go            # Admin post management
│   │   ├── auth.go             # Authentication handlers
│   │   ├── blog.go             # Blog post handlers
│   │   ├── handlers.go         # Main handler setup
│   │   └── home.go             # Homepage handler
│   ├── models/                  # Data models
│   │   └── post.go             # Post structure
│   └── storage/                 # Storage interface & implementation
│       ├── interface.go        # Storage interface
│       └── filesystem.go       # File-based storage
├── pkg/
│   └── markdown/                # Markdown processing
│       └── parser.go
├── web/
│   ├── templates/               # HTML templates
│   │   ├── layouts/
│   │   │   └── base.html
│   │   └── pages/
│   │       ├── home.html
│   │       ├── post.html
│   │       └── search.html
│   └── static/                  # Static assets (CSS, JS, images)
├── posts/                       # Markdown posts directory
│   ├── welcome-to-my-blog.md
│   └── building-go-web-apps.md
├── .github/
│   └── workflows/
│       └── deploy.yml           # GitHub Actions deployment
├── Dockerfile                   # Multi-stage Docker build
├── fly.toml                     # Fly.io configuration
├── go.mod
└── go.sum
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or later
- Git

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd goblog
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

4. **Open your browser**
   ```
   http://localhost:8080
   ```

### Adding Content

Create markdown files in the `posts/` directory with front matter:

```markdown
---
title: "Your Post Title"
tags: "tag1, tag2, tag3"
created: 2024-01-15
updated: 2024-01-20
---

# Your Post Title

Your content here...
```

## 🐳 Deployment

### Docker

1. **Build the image**
   ```bash
   docker build -t goblog .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 -v $(pwd)/posts:/data/posts goblog
   ```

### Fly.io Deployment

1. **Install Fly CLI** and login
   ```bash
   fly auth login
   ```

2. **Deploy the application**
   ```bash
   fly deploy
   ```

3. **Create persistent volume** (first time only)
   ```bash
   fly volumes create goblog_data --region ams --size 1
   ```

### GitHub Actions

The repository includes automated deployment:
- Pushes to `main` branch trigger deployment
- Tests run before deployment
- Automatic deployment to Fly.io

Set up the `FLY_API_TOKEN` secret in your GitHub repository settings.

## ⚙️ Configuration

Configure the application using environment variables:

```env
PORT=8080                               # Server port
POSTS_DIR=posts                         # Posts directory
GITHUB_CLIENT_ID=your_client_id         # GitHub OAuth (optional)
GITHUB_SECRET=your_client_secret        # GitHub OAuth (optional)
SESSION_SECRET=your_session_secret      # Session encryption
ALLOWED_USERS=username1,username2       # Allowed GitHub users
BASE_URL=https://yourdomain.com         # Base URL for OAuth
```

## 🎨 Features in Detail

### Theme System
- **Auto-detection** of system preference
- **Manual toggle** with persistent storage
- **Smooth transitions** between themes

### Search Functionality
- **Full-text search** in titles and headers
- **Tag-based filtering**
- **Real-time results**

### Markdown Support
- **GitHub Flavored Markdown** (GFM)
- **Syntax highlighting** for code blocks
- **Front matter** for post metadata
- **Auto-generated** table of contents

### Storage System
- **Interface-based design** for easy extension
- **Filesystem storage** with efficient parsing
- **Persistent volume** support for deployment

## 🛡️ Security Features

- **Authentication middleware** ready for GitHub OAuth
- **User whitelisting** for admin access
- **Session management** structure
- **Input validation** and sanitization

## 🧪 Testing

Run the test suite:

```bash
go test -v ./...
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

If you have any questions or need help, please open an issue or reach out via the social media links in the navigation.

---

**Built with ❤️ using Go and modern web technologies**
