# ls2pwsh

`ls2pwsh` converts `$LS_COLORS` strings into PowerShell compatible [PSStyle.FileInfo](https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.psstyle.fileinfoformatting) assignments, which control the styling used by `Get-ChildItem`/`dir`.

## Installation

```shell
go install github.com/CatsDeservePets/ls2pwsh@latest
```

## Usage

```
usage: ls2pwsh LS_COLORS
```

## Example

Use [vivid](https://github.com/sharkdp/vivid) to generate a theme and apply it:

```powershell
Invoke-Expression (vivid generate nord | ls2pwsh | Out-String)
```

## Disclaimer

`ls2pwsh` is not affiliated with, endorsed by, or sponsored by Microsoft or the PowerShell project.
