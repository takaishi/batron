# batron

batron is a deployment tool for AWS Batch inpired by ecspresso.

## Install

```bash
go install github.com/takaishi/batron/cmd/batron
```

## Usage

```
% batron --help
Usage: batron <command> [flags]

Flags:
  -h, --help       Show context-sensitive help.
      --version    show version

Commands:
  render [flags]
    Render job definition

  deploy [flags]
    Deploy job definition

  deregister [flags]
    Deregister old job definitions

Run "batron <command> --help" for more information on a command.
```
