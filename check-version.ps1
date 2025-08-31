# PowerShell 脚本：检查本地版本是否为最新
# 读取本地版本号
$localVersion = Get-Content -Path "version.txt" -Raw
# 获取远程版本号（从 GitHub 主分支 version.txt）
$remoteVersion = Invoke-RestMethod -Uri "https://raw.githubusercontent.com/Asice-Cloud/tz-gin-template/master/version.txt"

Write-Host "本地版本: $localVersion"
Write-Host "最新版本: $remoteVersion"

if ($localVersion -eq $remoteVersion) {
    Write-Host "当前已是最新版本。"
} else {
    Write-Host "有新版本可用，请及时更新！"
}
