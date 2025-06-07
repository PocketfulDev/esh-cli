class EshCli < Formula
  desc "ESH CLI tool for managing git tags and deployments"
  homepage "https://github.com/PocketfulDev/esh-cli"
  version "1.0.4"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/PocketfulDev/esh-cli/releases/download/v#{version}/esh-cli-darwin-arm64.tar.gz"
    sha256 "REPLACE_WITH_ARM64_SHA256"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/PocketfulDev/esh-cli/releases/download/v#{version}/esh-cli-darwin-amd64.tar.gz"
    sha256 "REPLACE_WITH_AMD64_SHA256"
  elsif OS.linux? && Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
    url "https://github.com/PocketfulDev/esh-cli/releases/download/v#{version}/esh-cli-linux-arm64.tar.gz"
    sha256 "REPLACE_WITH_LINUX_ARM64_SHA256"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/PocketfulDev/esh-cli/releases/download/v#{version}/esh-cli-linux-amd64.tar.gz"
    sha256 "REPLACE_WITH_LINUX_AMD64_SHA256"
  end

  def install
    if OS.mac?
      if Hardware::CPU.arm?
        bin.install "esh-cli-darwin-arm64" => "esh-cli"
      else
        bin.install "esh-cli-darwin-amd64" => "esh-cli"
      end
    elsif OS.linux?
      if Hardware::CPU.arm?
        bin.install "esh-cli-linux-arm64" => "esh-cli"
      else
        bin.install "esh-cli-linux-amd64" => "esh-cli"
      end
    end
  end

  test do
    system "#{bin}/esh-cli", "--help"
    assert_match "ESH CLI tool", shell_output("#{bin}/esh-cli --help")
  end
end
