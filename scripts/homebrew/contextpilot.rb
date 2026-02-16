# Homebrew formula for ContextPilot
# To use: brew tap contextpilot-dev/tap && brew install contextpilot
# Or copy to your own tap repository

class Contextpilot < Formula
  desc "Make every AI tool understand your codebase"
  homepage "https://contextpilot.dev"
  version "0.1.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/contextpilot-dev/contextpilot/releases/download/v#{version}/contextpilot-darwin-amd64.tar.gz"
      sha256 "c3e0c4d94713291f46b98c81964892871dbb12dd2b171f503b45b4c3ba1e227b"
    end
    on_arm do
      url "https://github.com/contextpilot-dev/contextpilot/releases/download/v#{version}/contextpilot-darwin-arm64.tar.gz"
      sha256 "c00db5899c973130fb1d1df7f380212dcc323c160295992f43ae9f9821644429"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/contextpilot-dev/contextpilot/releases/download/v#{version}/contextpilot-linux-amd64.tar.gz"
      sha256 "739fb55136bebf5aa2ce807998f74f84980af097d61561d34f7318fad2a5e8dc"
    end
    on_arm do
      url "https://github.com/contextpilot-dev/contextpilot/releases/download/v#{version}/contextpilot-linux-arm64.tar.gz"
      sha256 "f178ab306159b9cebbfb820e28d7aeb1bf225202b965c15d825e4f5cb0a45436"
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
