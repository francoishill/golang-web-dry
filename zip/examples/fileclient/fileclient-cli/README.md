# Example using ziputils as a file client (communicating with the example server)

## Commands
Commands are currently `UPLOAD`, `DOWNLOAD`, `DELETE`, `STATS`, `MOVE`.

All commands require
- the server url `-s` flag as well as​
- the remote path `-r` flag​​

Then only the four upload/download commands require
- the local path `-l` flag too

## Windows example
- Build with `go build -o=client.exe`
- Run with (be sure to replace the arguments - `SERVER`, `PORT`, `-l` and `-r`):
    ```
    client.exe -m UPLOADFOLDER^
      -s "http://SERVER:PORT"^
      -l "\path\to\local_file_to_upload"^
      -r "\path\to\remote_file_to_save_to"
    ```


## Linux example
- Build with `go build -o=client`
- Add executable flag with `chmod +x client`
- Run with (be sure to replace the arguments - `SERVER`, `PORT`, `-l` and `-r`):
    ```
    client -m UPLOADFOLDER \
      -s "http://SERVER:PORT" \
      -l "/path/to/local_file_to_upload" \
      -r "/path/to/remote_file_to_save_to"
    ```
