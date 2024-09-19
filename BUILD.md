# Build instructions:

(for windows)

```
cd vue
Remove-Item "dist" -Force -Recurse
npm run build
cd ..
go build -o lethal_backups.exe .

```
