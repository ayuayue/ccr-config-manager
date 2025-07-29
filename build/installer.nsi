; NSIS Installer Script for Claude Config Manager
; Define your application name
!define APPNAME "Claude Config Manager"
!define COMPANYNAME "Claude Config Manager"
!define DESCRIPTION "Configuration Manager for Claude Code Router"
!define INSTALLDIR "$APPDATA\${APPNAME}"

; Main Install settings
Name "${APPNAME}"
InstallDir "${INSTALLDIR}"
InstallDirRegKey HKCU "Software\${APPNAME}" ""
OutFile "dist/${APPNAME}-${VERSION}-windows-amd64-installer.exe"

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
  
  ; Add files
  File "build\bin\claudeConfigManager-windows-amd64.exe"
  File "build\bin\assets\*.*"
  
  ; Create directories
  SetOutPath "$INSTDIR\assets"
  File "build\bin\assets\*.*"
  
  ; Store installation folder
  WriteRegStr HKCU "Software\${APPNAME}" "" $INSTDIR
  
  ; Create uninstaller
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  
  ; Add to Add/Remove Programs
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" \
                   "DisplayName" "${APPNAME}"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" \
                   "UninstallString" "$INSTDIR\Uninstall.exe"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" \
                   "DisplayIcon" "$INSTDIR\claudeConfigManager-windows-amd64.exe,0"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" \
                   "DisplayVersion" "${VERSION}"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" \
                   "Publisher" "${COMPANYNAME}"
  
  ; Create desktop shortcut
  CreateShortCut "$DESKTOP\${APPNAME}.lnk" "$INSTDIR\claudeConfigManager-windows-amd64.exe"
  
  ; Create start menu shortcut
  CreateDirectory "$SMPROGRAMS\${APPNAME}"
  CreateShortCut "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk" "$INSTDIR\claudeConfigManager-windows-amd64.exe"
  CreateShortCut "$SMPROGRAMS\${APPNAME}\Uninstall.lnk" "$INSTDIR\Uninstall.exe"
SectionEnd

; Uninstaller Section
Section "Uninstall"
  ; Remove files
  Delete "$INSTDIR\claudeConfigManager-windows-amd64.exe"
  Delete "$INSTDIR\assets\*.*"
  RMDir "$INSTDIR\assets"
  
  ; Remove uninstaller
  Delete "$INSTDIR\Uninstall.exe"
  
  ; Remove desktop shortcut
  Delete "$DESKTOP\${APPNAME}.lnk"
  
  ; Remove start menu shortcut
  Delete "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk"
  Delete "$SMPROGRAMS\${APPNAME}\Uninstall.lnk"
  RMDir "$SMPROGRAMS\${APPNAME}"
  
  ; Remove installation directory
  RMDir "$INSTDIR"
  
  ; Remove registry keys
  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}"
  DeleteRegKey HKCU "Software\${APPNAME}"
SectionEnd