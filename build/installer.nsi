; NSIS Installer Script for Claude Config Manager
; Define your application name
!define APPNAME "ClaudeConfigManager"
!define COMPANYNAME "ClaudeConfigManager"
!define DESCRIPTION "Configuration Manager for Claude Code Router"
!define INSTALLDIR "$APPDATA\${APPNAME}"

; Main Install settings
Name "Claude Config Manager"
InstallDir "${INSTALLDIR}"
InstallDirRegKey HKCU "Software\${APPNAME}" ""
OutFile "claude-config-manager-setup.exe"

; Set working directory to build directory
!verbose 3

; Use compression
SetCompressor /SOLID lzma

; Request application privileges for Windows Vista
RequestExecutionLevel user

; Include Modern UI
!include "MUI2.nsh"

; Interface Settings
!define MUI_ABORTWARNING

; Pages
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; Languages
!insertmacro MUI_LANGUAGE "English"

; Installer Sections
Section "Install"
  SetOutPath "$INSTDIR"
  
  ; Add files (with .exe extension for Windows)
  File "bin\claudeConfigManager-windows_amd64.exe"
  
  ; Store installation folder
  WriteRegStr HKCU "Software\${APPNAME}" "" $INSTDIR
  
  ; Create uninstaller
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  
  ; Add to Add/Remove Programs
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" \
                   "DisplayName" "Claude Config Manager"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" \
                   "UninstallString" "$INSTDIR\Uninstall.exe"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" \
                   "DisplayIcon" "$INSTDIR\claudeConfigManager-windows_amd64.exe,0"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" \
                   "DisplayVersion" "${VERSION}"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" \
                   "Publisher" "${COMPANYNAME}"
  
  ; Create desktop shortcut
  CreateShortCut "$DESKTOP\Claude Config Manager.lnk" "$INSTDIR\claudeConfigManager-windows_amd64.exe"
  
  ; Create start menu shortcut
  CreateDirectory "$SMPROGRAMS\Claude Config Manager"
  CreateShortCut "$SMPROGRAMS\Claude Config Manager\Claude Config Manager.lnk" "$INSTDIR\claudeConfigManager-windows_amd64.exe"
  CreateShortCut "$SMPROGRAMS\Claude Config Manager\Uninstall.lnk" "$INSTDIR\Uninstall.exe"
SectionEnd

; Uninstaller Section
Section "Uninstall"
  ; Remove files
  Delete "$INSTDIR\claudeConfigManager-windows_amd64.exe"
  
  ; Remove uninstaller
  Delete "$INSTDIR\Uninstall.exe"
  
  ; Remove desktop shortcut
  Delete "$DESKTOP\Claude Config Manager.lnk"
  
  ; Remove start menu shortcut
  Delete "$SMPROGRAMS\Claude Config Manager\Claude Config Manager.lnk"
  Delete "$SMPROGRAMS\Claude Config Manager\Uninstall.lnk"
  RMDir "$SMPROGRAMS\Claude Config Manager"
  
  ; Remove installation directory
  RMDir "$INSTDIR"
  
  ; Remove registry keys
  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}"
  DeleteRegKey HKCU "Software\${APPNAME}"
SectionEnd