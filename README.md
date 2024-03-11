# Authenticator

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/authenticator-cli)](https://github.com/nhatthm/authenticator-cli/releases/latest)
[![Build Status](https://github.com/nhatthm/authenticator-cli/actions/workflows/release-edge.yaml/badge.svg)](https://github.com/nhatthm/authenticator-cli/actions/workflows/release-edge.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/authenticator-cli)](https://goreportcard.com/report/github.com/nhatthm/authenticator-cli)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/authenticator-cli)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

<!--
[![codecov](https://codecov.io/gh/nhatthm/authenticator-cli/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/authenticator-cli)
-->

CLI tool for managing and generating one-time passwords for multiple accounts.

## Prerequisites

- `Go >= 1.22`

### Keyring

Support **OS X**, **Linux/BSD (dbus)** and **Windows**.

#### OS X

The OS X implementation depends on the `/usr/bin/security` binary for
interfacing with the OS X keychain. It should be available by default.

#### Linux and *BSD

The Linux and *BSD implementation depends on the [Secret Service][SecretService] dbus
interface, which is provided by [GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring).

It's expected that the default collection `login` exists in the keyring, because
it's the default in most distros. If it doesn't exist, you can create it through the
keyring frontend program [Seahorse](https://wiki.gnome.org/Apps/Seahorse):

* Open `seahorse`
* Go to **File > New > Password Keyring**
* Click **Continue**
* When asked for a name, use: **login**

## Install

You can download the [latest stable version](https://github.com/nhatthm/authenticator-cli/releases/latest) or
the [nightly build](https://github.com/nhatthm/authenticator-cli/releases/tag/edge) (`edge` version).

Once downloaded, the binary can be run from anywhere. Ideally, though, you should move it into your `$PATH` for easy use. `/usr/local/bin` is a popular location
for this.

To update the tool to the newest version, run `authenticator self-update`.

### Install from source

If you have `go` installed, you can run the following command to install the latest version:

```bash
$ /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/nhatthm/authenticator-cli/HEAD/install.sh)"
```

The binary will be installed to `$GOBIN` or `$GOPATH/bin` (when `$GOBIN` is empty) directory. If you don't know where it is, you can run `go env` to find out.

In order run the binary from anywhere, the `$GOBIN` or `$GOPATH/bin` directory should be added to your `$PATH`.

## Usage

### Add a new account

Run `authenticator account add` to add a new account. If you don't specify the namespace, the tool will ask you to input it.

If you have a QR code, you can use `--qr` flag to tell the tool to [scan the QR code](#scan-a-qr-code). Otherwise, [input the account and the totp secret manually](#manually-input-all-the-information).

```bash
$ authenticator account add -h
Add an account

Usage:
  authenticator account add [-n <namespace>] [--qr </path/to/qr-code-image>] [flags]

Flags:
  -h, --help               help for add
  -n, --namespace string   namespace
      --qr string          qr code

Global Flags:
  -d, --debug     debug output
  -v, --verbose   verbose output
```

### Generate an OTP

After [adding an account](#add-a-new-account), you can generate an OTP by running `authenticator otp <namespace> <account>`. For example:

```bash
$ authenticator otp -n demo john.doe@example.com
103281
```

If you want to generate and copy the OTP to the clipboard, use the `--copy` flag, for example:

```bash
$ authenticator otp -n demo john.doe@example.com --copy
```

> [!TIP]
> You can add an alias to your shell startup file, such as `.bashrc` or `.zshrc`, to make it easier to use. For example:
>
> ```bash
> alias cotp='authenticator otp demo john.doe@example.com --copy'
> ```
>
> Later, you can just run `cotp` to generate and copy the OTP to the clipboard.
>
> If your terminal allows you to customize the shortcuts, you can also create a shortcut for the alias
>
> <p align="center">
>     <img width="70%" alt="image" src="https://github.com/nhatthm/authenticator-cli/assets/1154587/41488032-a691-49fd-81aa-5bee9aea306a">
> </p>

## Examples

### Manually input all the information

<p align="center">
    <img src="./resources/docs/demo1.gif" alt="demo 1" width="100%" height="auto" /><br/>
    <i>Manually input all the information</i>
</p>

### Scan a QR code

Example QR code:

<p align="center">
    <img src="./resources/fixtures/qr.png" alt="qr" height="auto" /><br/>
    <code>totpauth://otp/john.doe%40example.com?secret=NBSWY3DP&issuer=example.com</code>
    <br/>
    <br/>
    <img src="./resources/docs/demo2.gif" alt="demo 2" width="100%" height="auto" /><br/>
</p>

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
