# ���ñ���
$binName = "ethereum-wallet-generator-nodes.exe"
$zipName = "ethereum-wallet-generator-nodes-v9.9.9-windows-amd64.zip"
$downloadUrl = "https://github.com/seth-shi/ethereum-wallet-generator-nodes/releases/download/v9.9.9/ethereum-wallet-generator-nodes-v9.9.9-windows-amd64.zip"

# �����ļ�
if (Test-Path $zipName) {
    Remove-Item -Path $zipName -Force
    Write-Host "�ļ� $zipName �ѱ�ɾ����"
}
Invoke-WebRequest -Uri $downloadUrl -OutFile $zipName

# ��ѹ�ļ�
if (Test-Path $binName) {
    Remove-Item -Path $binName -Force
    Write-Host "�ļ� $binName �ѱ�ɾ����"
}
Expand-Archive -Force -Path $zipName -DestinationPath .

# У���ļ��Ƿ���ȷ
Write-Host "Զ���ļ� md5"
Invoke-WebRequest -Uri "$downloadUrl.md5" -OutFile "$zipName.md5"
Get-FileHash -Algorithm MD5 -Path "$zipName.md5" | Select-Object -ExpandProperty Hash
Write-Host "�����ļ� md5"
Get-FileHash -Algorithm MD5 -Path $zipName | Select-Object -ExpandProperty Hash

# ɾ��ѹ���ļ�
Remove-Item -Path "$zipName.md5" -Force
Remove-Item -Path $zipName -Force

Write-Host "�������"
