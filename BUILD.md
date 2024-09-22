# Build instructions:

(for windows)

```
cd vue
Remove-Item "dist" -Force -Recurse
npm run build
cd ..
go build -ldflags "-w -s" -o lethal_backups.exe .
```

```
