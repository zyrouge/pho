# Phở

[![Latest Version](https://img.shields.io/github/v/release/zyrouge/pho?label=latest)](https://github.com/zyrouge/pho/releases/latest)
[![Build](https://github.com/zyrouge/pho/actions/workflows/build.yml/badge.svg)](https://github.com/zyrouge/pho/actions/workflows/build.yml)
[![Release](https://github.com/zyrouge/pho/actions/workflows/release.yml/badge.svg)](https://github.com/zyrouge/pho/actions/workflows/release.yml)

<div align="center">
    <img src="./media/banner.png">
</div>

## Features

-   Manage your AppImages by organizing them in a single folder.
-   Integrates your AppImages seamlessly.
-   Ability to download AppImages from Github Releases and URLs.
-   Supports updation of AppImages.
-   You can manually edit configuration files of Pho to further customize functionality.

##### Notes:

-   AppImages must follow AppImage Specification to be integrated with desktop.
-   Only for AppImages fetched from Github Releases support updating.

## Installation

1. All releases can be found [here](https://github.com/zyrouge/pho/releases). Choose a valid release.

2. Binaries are built for 32-bit/64-bit AMD and ARM separately. Download the appropriate one.

| Binary name | Platform   |
| ----------- | ---------- |
| `pho-amd`   | 32-bit AMD |
| `pho-amd64` | 64-bit AMD |
| `pho-arm`   | 32-bit ARM |
| `pho-arm64` | 32-bit ARM |

3. Place your downloaded binary in a folder that is available in the `PATH` environmental variable. Typically this would be `~/.local/bin`.

4. Run `pho init` to setup necessary configuration.

5. Have fun! 🎉

## License

[MIT](./LICENSE)
