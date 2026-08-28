**这个仓库用户信息的交流和资源分享，有好的可以写到README，然后资源自建文件夹放入**
## 前端部署到 Vercel

前端位于 `web/`，可将该目录作为 Vercel 的 Root Directory。

```text
Framework Preset: Vite
Build Command: npm run build
Output Directory: dist
Install Command: npm install
```

在 Vercel 项目环境变量中设置 `VITE_API_BASE_URL` 为线上 Go API 地址，例如 `https://api.example.com`。如果使用相对路径和 Vercel rewrite，请先将 `web/vercel.json` 中的 `api.example.com` 替换为真实后端域名；Vercel 配置文件不会自动读取 `VITE_API_PROXY_URL`。后端也需要允许来自 Vercel 域名的请求。

本地开发：

```powershell
cd web
npm install
npm run dev
```
