# How to use

## Check archive sha256 (optional)

This step is optional but highly recommended, the output of the command below shoud be the same as in the .sha256sum files

Otherwise DO NOT follow the next steps, open an issue if you have the time.

```bash
sha256sum spotify-backup.tar.gz
```

## Extract the archive with tar

```bash
tar xzvf spotify-backup.tar.gz
```

## Run

```bash
./spotify-backup
```
