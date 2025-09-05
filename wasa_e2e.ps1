
param(
  [string]$BaseUrlApi = "http://localhost:3000",
  [string]$BaseUrlWeb = "http://localhost:8080",
  [string]$BackendTag = "wasa-backend:final",
  [string]$FrontendTag = "wasa-frontend:final",
  [string]$UserA = "Maria",
  [string]$UserB = "Luca"
)

$ErrorActionPreference = "Stop"

function Step($m){ Write-Host "`n==> $m" -ForegroundColor Cyan }
function Ok($m){ Write-Host $m -ForegroundColor Green }
function Warn($m){ Write-Warning $m }
function Fail($m){ Write-Host $m -ForegroundColor Red }

function Try-Http($method, $url, $headers=$null, $body=$null){
  try {
    if ($null -ne $body) {
      return Invoke-RestMethod -Method $method -Uri $url -Headers $headers -ContentType 'application/json' -Body $body -TimeoutSec 5
    } else {
      return Invoke-RestMethod -Method $method -Uri $url -Headers $headers -TimeoutSec 5
    }
  } catch {
    return $null
  }
}

# 0) Execution policy hint (non blocca lo script, solo info)
try { $null = Get-ExecutionPolicy } catch { Warn "Se vedi errori di policy, esegui: Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass" }

# 1) Build backend
if (Test-Path "./Dockerfile.backend") {
  Step "Build backend image ($BackendTag)"
  docker build -f Dockerfile.backend -t $BackendTag . | Out-Host
} else {
  Fail "Dockerfile.backend non trovato nella cartella corrente."
  exit 1
}

# 2) Build frontend (se presente)
$frontendBuilt = $false
if (Test-Path "./Dockerfile.frontend") {
  Step "Build frontend image ($FrontendTag) usando ./Dockerfile.frontend"
  docker build -f Dockerfile.frontend -t $FrontendTag . | Out-Host
  $frontendBuilt = $true
} elseif (Test-Path "./webui/Dockerfile") {
  Step "Build frontend image ($FrontendTag) usando ./webui/Dockerfile"
  docker build -f ./webui/Dockerfile -t $FrontendTag ./webui | Out-Host
  $frontendBuilt = $true
} else {
  Warn "Dockerfile frontend non trovato. Salto la build FE."
}

# 3) Network, stop & run
Step "Crea network 'wasa-net' (se non esiste)"
docker network create wasa-net 2>$null | Out-Null

Step "Stop & remove container precedenti"
docker rm -f wasa-backend 2>$null | Out-Null
if ($frontendBuilt) { docker rm -f wasa-frontend 2>$null | Out-Null }

Step "Run backend su :3000 (user 1000)"
docker run -d --name wasa-backend --network wasa-net --user 1000 -p 3000:3000 `
  -e PORT=3000 -e GIN_MODE=release $BackendTag | Out-Host

# 4) Wait backend ready
Step "Attendo backend pronto su $BaseUrlApi"
$ready = $false
for ($i=0; $i -lt 30; $i++){
  $ping = Try-Http "GET" ($BaseUrlApi + "/debug")
  if ($ping -ne $null) { $ready = $true; break }
  Start-Sleep -Seconds 1
}
if (-not $ready) {
  Warn "Non ho ricevuto /debug; provo un GET generico sulla root"
  $root = Try-Http "GET" ($BaseUrlApi + "/")
  if ($root -eq $null) {
    Fail "Backend non raggiungibile su $BaseUrlApi. Controlla 'docker logs --tail=100 wasa-backend'."
    exit 1
  }
}
Ok "Backend OK"

# 5) Run frontend (se buildato)
if ($frontendBuilt) {
  Step "Run frontend su :8080 (nginx)"
  docker run -d --name wasa-frontend --network wasa-net -p 8080:80 $FrontendTag | Out-Host

  Step "Attendo frontend pronto su $BaseUrlWeb"
  $webok = $false
  for ($i=0; $i -lt 30; $i++){
    $html = Try-Http "GET" $BaseUrlWeb
    if ($html -ne $null) { $webok = $true; break }
    Start-Sleep -Seconds 1
  }
  if ($webok) { Ok "Frontend OK" } else { Warn "Frontend non ha risposto 200 su $BaseUrlWeb (forse buildato senza Dockerfile? Verifica)"; }
}

# 6) Smoke test backend
Step "Smoke test backend: login + sendMessage + list"
$respA = Invoke-RestMethod -Method Post -Uri ($BaseUrlApi + "/session") `
  -ContentType 'application/json' -Body ('"' + $UserA + '"')
$tokA = $respA.identifier; if (-not $tokA) { $tokA = $respA.token }
if (-not $tokA) { Fail "Non riesco a estrarre il token"; exit 1 }
$HeadersA = @{ Authorization = "Bearer $tokA" }
Ok ("Token A: " + $tokA.Substring(0,20) + "...")

# sendMessage (auto-create chat con UserB)
$bodySend = @{ recipientusername = $UserB; message = @{ text = "Ciao da " + $UserA } } | ConvertTo-Json
$msg = Invoke-RestMethod -Method Post -Uri ($BaseUrlApi + "/api/conversations/messages") `
  -Headers $HeadersA -ContentType 'application/json' -Body $bodySend
$msg | ConvertTo-Json -Depth 6 | Write-Host

# list convos
$convos = Invoke-RestMethod -Method Get -Uri ($BaseUrlApi + "/api/conversations") -Headers $HeadersA
$convos | ConvertTo-Json -Depth 6 | Write-Host

# 7) (Opzionale) Smoke via frontend proxy
if ($frontendBuilt) {
  Step "Test via frontend: GET / (static) + /api/conversations (proxy)"
  $home = Try-Http "GET" $BaseUrlWeb
  if ($home -ne $null) { Ok "GET / (frontend) OK" } else { Warn "GET / (frontend) non OK" }

  $convFE = Try-Http "GET" ($BaseUrlWeb + "/api/conversations") $HeadersA
  if ($convFE -ne $null) {
    Ok "Proxy FE → BE OK (GET /api/conversations via frontend)"
    $convFE | ConvertTo-Json -Depth 6 | Write-Host
  } else {
    Warn "Il proxy /api dal frontend non risponde. Possibili cause: nginx.conf non punta a 'wasa-backend:3000' o usa host.docker.internal."
    Warn "Il backend comunque è OK (vedi smoke diretto)."
  }
}

Ok "E2E completato."
