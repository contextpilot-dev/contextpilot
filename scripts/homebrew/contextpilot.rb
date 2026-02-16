# Homebrew formula for ContextPilot
# To use: brew tap jitin-nhz/tap && brew install contextpilot
# Or copy to your own tap repository

class Contextpilot < Formula
  desc "Make every AI tool understand your codebase"
  homepage "https://contextpilot.dev"
  version "0.1.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/jitin-nhz/contextpilot/releases/download/v#{version}/contextpilot-darwin-amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_DARWIN_AMD64"
    end
    on_arm do
      url "https://github.com/jitin-nhz/contextpilot/releases/download/v#{version}/contextpilot-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_DARWIN_ARM64"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/jitin-nhz/contextpilot/releases/download/v#{version}/contextpilot-linux-amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"
    end
    on_arm do
      url "https://github.com/jitin-nhz/contextpilot/releases/download/v#{version}/contextpilot-linux-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"
    end
  end

  def install
    binary_name = "contextpilot"
    if OS.mac?
      binary_name += "-darwin"
    else
      binary_name += "-linux"
    end
    binary_name += Hardware::CPU.arm? ? "-arm64" : "-amd64"
    
    bin.install binary_name => "contextpilot"
  end

  test do
    assert_match "contextpilot", shell_output("#{bin}/contextpilot --help")
  end
end
