class SqlHttpProxy < Formula
  desc "YAML configuration-based HTTP to SQL proxy server"
  homepage "https://github.com/mpyw/sql-http-proxy"
  license "MIT"
  version "${VERSION}"

  on_macos do
    on_arm do
      url "https://github.com/mpyw/sql-http-proxy/releases/download/v${VERSION}/sql-http-proxy_${VERSION}_darwin_arm64.tar.gz"
      sha256 "${SHA256_DARWIN_ARM64}"
    end
    on_intel do
      url "https://github.com/mpyw/sql-http-proxy/releases/download/v${VERSION}/sql-http-proxy_${VERSION}_darwin_amd64.tar.gz"
      sha256 "${SHA256_DARWIN_AMD64}"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/mpyw/sql-http-proxy/releases/download/v${VERSION}/sql-http-proxy_${VERSION}_linux_arm64.tar.gz"
      sha256 "${SHA256_LINUX_ARM64}"
    end
    on_intel do
      url "https://github.com/mpyw/sql-http-proxy/releases/download/v${VERSION}/sql-http-proxy_${VERSION}_linux_amd64.tar.gz"
      sha256 "${SHA256_LINUX_AMD64}"
    end
  end

  def install
    bin.install "sql-http-proxy"
  end

  test do
    system bin/"sql-http-proxy", "--version"
  end
end
