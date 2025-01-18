A website to download Youtube, TikTok, X and Instagram reels for free without Ads/watermark.

## Setting up the repository

- Backend (Go, yt-dl)
- Frontend (React, TailwindCSS, Vite)

Ensure you have [`Node JS`](https://github.com/nvm-sh/nvm), [`yt-dlp`](https://github.com/yt-dlp/yt-dlp) and [`Go`](go.dev) installed on your machine.

```bash
go version
node -v
make --v
```

## Running the project

1. Clone the repository

```bash
git clone https://github.com/wathika-eng/ytdl_go --depth 1 && cd ytdl_go
```

2. Open 2 tabs in your terminal and run the following commands:

```bash
cd backend && make run
```

```bash
cd frontend && npm install && npm run dev
```
