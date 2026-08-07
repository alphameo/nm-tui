class NmTui < Formula
  desc "Lightweight TUI wrapper for NetworkManager"
  homepage "https://github.com/alphameo/nm-tui"
  url "https://github.com/alphameo/nm-tui/archive/refs/tags/v0.2.0.tar.gz"
  sha256 "0000000000000000000000000000000000000000000000000000000000000000"
  license "MIT"
  head "https://github.com/alphameo/nm-tui.git", branch: "main"

  depends_on "go" => :build
  depends_on "networkmanager"

  def install
    ldflags = %W[-s -w -X main.version=#{version}]
    system "go", "build", *std_go_args(ldflags: ldflags), "./cmd/nm-tui/main.go"
    man1.install "docs/nm-tui.1" if File.exist?("docs/nm-tui.1")
  end

  test do
    assert_match "nm-tui #{version}", shell_output("#{bin}/nm-tui --version")
  end
end