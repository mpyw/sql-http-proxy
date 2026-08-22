## Downloads

### macOS

| Architecture | Download |
|--------------|----------|
| Apple Silicon (M1/M2/M3/M4) | [sql-http-proxy_${VERSION}_darwin_arm64.tar.gz](${BASE_URL}/sql-http-proxy_${VERSION}_darwin_arm64.tar.gz) |
| Intel | [sql-http-proxy_${VERSION}_darwin_amd64.tar.gz](${BASE_URL}/sql-http-proxy_${VERSION}_darwin_amd64.tar.gz) |

> [!TIP]
> If macOS shows "cannot be opened because the developer cannot be verified", run:
> ```bash
> xattr -d com.apple.quarantine /path/to/sql-http-proxy
> ```

### Windows

| Architecture | Download |
|--------------|----------|
| x86_64 | [sql-http-proxy_${VERSION}_windows_amd64.zip](${BASE_URL}/sql-http-proxy_${VERSION}_windows_amd64.zip) |
| ARM64 | [sql-http-proxy_${VERSION}_windows_arm64.zip](${BASE_URL}/sql-http-proxy_${VERSION}_windows_arm64.zip) |

> [!TIP]
> On first run, if Windows SmartScreen shows "Windows protected your PC", click **More info** → **Run anyway**.

### Linux

| Architecture | Tarball | Debian/Ubuntu | RHEL/Fedora |
|--------------|---------|---------------|-------------|
| x86_64 | [sql-http-proxy_${VERSION}_linux_amd64.tar.gz](${BASE_URL}/sql-http-proxy_${VERSION}_linux_amd64.tar.gz) | [sql-http-proxy_${VERSION}-1_amd64.deb](${BASE_URL}/sql-http-proxy_${VERSION}-1_amd64.deb) | [sql-http-proxy-${VERSION}-1.x86_64.rpm](${BASE_URL}/sql-http-proxy-${VERSION}-1.x86_64.rpm) |
| ARM64 | [sql-http-proxy_${VERSION}_linux_arm64.tar.gz](${BASE_URL}/sql-http-proxy_${VERSION}_linux_arm64.tar.gz) | [sql-http-proxy_${VERSION}-1_arm64.deb](${BASE_URL}/sql-http-proxy_${VERSION}-1_arm64.deb) | [sql-http-proxy-${VERSION}-1.aarch64.rpm](${BASE_URL}/sql-http-proxy-${VERSION}-1.aarch64.rpm) |

Every binary is a static, CGO-free build with all database drivers included.

### Checksums

[checksums.txt](${BASE_URL}/checksums.txt)
