# Example using ziputils as a file client (communicating with the example server)

# Windows example
- Build with `go build -o=client.exe`
- Run with (be sure to replace the arguments - `SERVER`, `PORT`, `-l` and `-r`):
    ```
    client.exe -m UPLOADFOLDER^
      -s "http://SERVER:PORT"^
      -l "\path\to\local_file_to_upload"^
      -r "\path\to\remote_file_to_save_to"
    ```


# Linus example
- Build with `go build -o=client`
- Add executable flag with `chmod +x client`
- Run with (be sure to replace the arguments - `SERVER`, `PORT`, `-l` and `-r`):
    ```
    client -m UPLOADFOLDER \
      -s "http://SERVER:PORT" \
      -l "/path/to/local_file_to_upload" \
      -r "/path/to/remote_file_to_save_to"
    ```
