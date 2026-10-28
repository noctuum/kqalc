# kqalc

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/4fd9e15bc6e34e2c87e7a114963c6f34)](https://app.codacy.com/gh/noctuum/kqalc/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/4fd9e15bc6e34e2c87e7a114963c6f34)](https://app.codacy.com/gh/noctuum/kqalc/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

Full [qalculate](https://qalculate.github.io/) power in KRunner.

KRunner already converts a single value between units, and between the currencies of the ECB reference rates. Its calculator handles plain arithmetic. That calculator removes the spaces from a query and skips any query with a letter in it, unless the query calls a known function or starts or ends with `=`. **kqalc** is a KRunner D-Bus plugin that passes the query to the `qalc` command of libqalculate as typed. It converts every currency libqalculate has rates for, a conversion can end an expression, an equation needs no prefix, and an exact result appears beside the approximate one.

## Features

- **Math**: `qc 2+2`, `qc sqrt(2)`, `qc sin(pi/4)`, `qc 2^64`
- **Currency conversion**: `qc 5000 GEL to USD`, `qc 100 EUR to KZT`, `qc 100 USD + 50 EUR to GEL`
- **Units**: `qc 100 km/h to mph`, `qc 180 lb to kg`, `qc 5 km + 300 m`
- **Equations**: `qc x^2 = 9`, `qc solve(2x+5=15, x)`
- **Exact and approximate results** side by side where both are useful: `qc 2/6` gives `0.33333333` and `1/3`
- **Copy to clipboard** on selection, on Wayland and X11

## Requirements

- KDE Plasma 6
- [libqalculate](https://qalculate.github.io/) (the `qalc` command)
- `wl-copy` on Wayland or `xclip` on X11

## Install

### KDE Store

In Discover, open Application Addons → Krunner → System Runners and install [kqalc](https://store.kde.org/p/2371030). On Arch, install `packagekit-qt6` first. Plasma lists it as optional, but `krunner-plugininstaller`, the tool that installs KRunner plugins, is linked against it and does not start without it.

### From source

```bash
git clone https://github.com/noctuum/kqalc.git
cd kqalc
./install.sh
```

Requires Go 1.25.7 or later, the version `go.mod` names.

### Debian / Ubuntu

```bash
curl -fsSL https://noctuum.github.io/kqalc/gpg-key.asc | sudo tee /usr/share/keyrings/kqalc.asc > /dev/null
echo "deb [signed-by=/usr/share/keyrings/kqalc.asc] https://noctuum.github.io/kqalc stable main" | sudo tee /etc/apt/sources.list.d/kqalc.list
sudo apt update && sudo apt install kqalc
```

### Arch Linux (AUR)

```bash
paru -S kqalc-bin
```

### Nix

The flake in this repository builds the package. `nix run` starts it without installing it, which is not enough for KRunner; see below.

```bash
nix run github:noctuum/kqalc
```

## After installing

kqalc runs as a D-Bus service. The bus starts it on the first `qc` query. It exits 60 seconds after a KRunner session ends, when every cached result has expired, and the next query starts it again from whatever binary is installed by then.

**From source.** Nothing to do. `install.sh` starts kqalc itself, so `qc 2+2` works right away, and it stops a copy left by an earlier install.

**From a package: AUR, Debian, Ubuntu.** Restart KRunner so it picks up the new plugin:

```bash
systemctl --user restart plasma-krunner.service
```

Versions up to 0.2.2 never exit on their own. After an upgrade from one of them, stop the old copy once, or the old binary keeps answering:

```bash
pkill -x kqalc
```

**From the KDE Store.** Usually nothing. If `qc` finds nothing at all, the bus has not picked up the new service file: it watches only the service directories that existed when it started, and this install created one. Start kqalc by hand once. Discover unpacks the archive under `~/.local/share/krunner-sources/`, and the `kqalc` file inside it is the one to run. It stays up until the bus can start it on demand.

**Nix.** `nix run` builds and runs the binary without installing it, so KRunner never sees the plugin metadata. Add the package to a profile or to your system or home-manager configuration instead, so its `share/` directory reaches `XDG_DATA_DIRS`, then restart KRunner.

## Usage

Open KRunner with `Alt+Space` and type:

```text
qc <expression>
```

The prefix `qc` followed by a space triggers kqalc. Select a result to copy it to the clipboard.

## Uninstall

A source install is removed with:

```bash
./uninstall.sh
```

A package is removed with its package manager, and a KDE Store install from Discover.

## License

[GPL-2.0](LICENSE)
