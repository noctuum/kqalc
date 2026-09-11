# kqalc

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/4fd9e15bc6e34e2c87e7a114963c6f34)](https://app.codacy.com/gh/noctuum/kqalc/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/4fd9e15bc6e34e2c87e7a114963c6f34)](https://app.codacy.com/gh/noctuum/kqalc/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

Full [qalculate](https://qalculate.github.io/) power in KRunner.

KDE Plasma's built-in calculator is limited — no currency conversion with `to`, no equation solving, no unit arithmetic. **kqalc** fixes this by wrapping the `qalc` CLI as a KRunner DBus plugin.

## Features

- **Math**: `qc 2+2`, `qc sqrt(2)`, `qc sin(pi/4)`, `qc 2^64`
- **Currency conversion**: `qc 5000 GEL to USD`, `qc 100 EUR to KZT`
- **Unit conversion**: `qc 100 km/h to mph`, `qc 180 lbs to kg`
- **Equations**: `qc x^2 = 9`, `qc solve(2x+5=15, x)`
- **Exact & approximate** results shown side by side where useful
- **Copy to clipboard** on selection (Wayland & X11)

## Requirements

- KDE Plasma 6
- [libqalculate](https://qalculate.github.io/) (`qalc` CLI)
- `wl-copy` (Wayland) or `xclip` (X11)

## Install

### KDE Store

Open Discover, go to Krunner → System Runners and install [kqalc](https://store.kde.org/p/2371030). On Arch, install `packagekit-qt6` first — Plasma lists it as optional, but the tool that installs KRunner plugins is linked against it and will not start without it.

### From source

```bash
git clone https://github.com/noctuum/kqalc.git
cd kqalc
./install.sh
```

Requires Go 1.22+.

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

```bash
nix run github:noctuum/kqalc
```

## After installing

kqalc runs as a D-Bus service, and the session has to notice it before `qc` works. What that takes depends on where the service file landed.

**Installed from a package — AUR, Debian/Ubuntu.** The service file goes to `/usr/share/dbus-1/services/`, which the message bus already watches. Restart KRunner and you are done:

```bash
kquitapp6 krunner && kstart6 krunner
```

**Installed from source or from the KDE Store.** These write the service file into your home directory instead. If nothing has ever installed a user D-Bus service on this machine, the directory is brand new — and the bus only watches directories that existed when it started, so it will not see the runner. Logging out does not help either, because the user bus outlives the session. **Reboot once.** Every later install and update is picked up immediately, with no reboot.

**Nix.** `nix run` builds and runs the binary without installing it, so KRunner never sees the plugin metadata. Add the package to a profile or to your system or home-manager configuration instead, so its `share/` directory reaches `XDG_DATA_DIRS`, then restart KRunner.

## Usage

Open KRunner (`Alt+Space`) and type:

```
qc <expression>
```

The `qc ` prefix triggers kqalc. Results appear instantly — select to copy to clipboard.

## Uninstall

```bash
./uninstall.sh
```

## License

[GPL-2.0](LICENSE)
