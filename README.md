# kqalc

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/1d662b12055f47fa84ce602bb1da2aab)](https://app.codacy.com/gh/noctuum/kqalc?utm_source=github.com&utm_medium=referral&utm_content=noctuum/kqalc&utm_campaign=Badge_Grade)

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

### From source

```bash
git clone https://github.com/noctuum/kqalc.git
cd kqalc
./install.sh
```

Requires Go 1.22+. Restart KRunner after install:

```bash
kquitapp6 krunner && kstart6 krunner
```

### Arch Linux (AUR)

```bash
paru -S kqalc
```

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
