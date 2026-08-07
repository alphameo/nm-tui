Name:           nm-tui
Version:        0.2.0
Release:        1%{?dist}
Summary:        Lightweight TUI wrapper for NetworkManager

License:        MIT
URL:            https://github.com/alphameo/nm-tui
Source0:        %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang >= 1.26

%description
Lightweight TUI wrapper for NetworkManager, built with Bubbletea.

%prep
%autosetup -n nm-tui-%{version}

%build
CGO_ENABLED=0 go build -trimpath \
    -ldflags "-s -w -X main.version=%{version}" \
    -o bin/nm-tui ./cmd/nm-tui/main.go

%install
install -Dm755 bin/nm-tui %{buildroot}%{_bindir}/nm-tui
install -Dm644 LICENSE %{buildroot}%{_licensedir}/nm-tui/LICENSE
install -Dm644 README.md %{buildroot}%{_datadir}/doc/nm-tui/README.md

%files
%license LICENSE
%doc README.md
%{_bindir}/nm-tui

%changelog
* %{?date} alphameo - 0.2.0-1
- Initial package build