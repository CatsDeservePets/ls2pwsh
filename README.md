# ls2pwsh

`ls2pwsh` converts between `LS_COLORS` strings and PowerShell [PSStyle.FileInfo](https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.psstyle.fileinfoformatting) assignments, which control the styling used by `Get-ChildItem`/`dir`.

> [!NOTE]
> The BSD variant `LSCOLORS` is currently not supported.

## Installation

```shell
go install github.com/CatsDeservePets/ls2pwsh@latest
```

## Usage

```
usage: ls2pwsh LS_COLORS | PSStyle.FileInfo

Convert color strings between LS_COLORS and PowerShell PSStyle.FileInfo format

If the input is a single dash ('-') or absent, ls2pwsh reads from the standard input.
```

> [!NOTE]
> When converting from PowerShell, `ls2pwsh` expects the string form of `PSStyle.FileInfo`. Its PowerShell output is emitted as assignment statements.

## Examples

On Windows, `$env:LS_COLORS` is not set by default. If you don't already have an `LS_COLORS` value, you can generate one with tools such as [vivid](https://github.com/sharkdp/vivid).

Apply `LS_COLORS` in PowerShell:

```powershell
Invoke-Expression (ls2pwsh $env:LS_COLORS | Out-String)
```

Generate `LS_COLORS` with `vivid` and apply it in PowerShell:

```powershell
Invoke-Expression (vivid generate nord | ls2pwsh | Out-String)
```

| Before | After |
|---|---|
| <img width="271" height="271" src="https://github.com/user-attachments/assets/50be774d-c655-4a6e-94c5-8c1ef7e0577e" /> | <img width="271" height="271" src="https://github.com/user-attachments/assets/5cac4f2f-9184-46d9-a82b-d544da17a7ff" /> |

Set `LS_COLORS` from `PSStyle.FileInfo`:

```powershell
$env:LS_COLORS = ($PSStyle.FileInfo | ls2pwsh).Trim()
```

## Disclaimer

`ls2pwsh` is not affiliated with, endorsed by, or sponsored by Microsoft or the PowerShell project.
