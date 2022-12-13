# mje

## Build

```
$ make build
$ ./bin/mje --help

```
## Usage

### Backup by user-id
To backup your prompts json and images you can run the following.

```
# Grab your browser Cookie session token from midjourney.com "__Secure-next-auth.session-token"
$ MIDJOURNEY_TOKEN="ey....."
$ mje midjourney backup-jobs --token "$MIDJOURNEY_TOKEN" --user-id=264088956247867392 --output backup-jobs/upscale --type upscale

INFO[0003] SKIPPING - File already exists for url: https://storage.googleapis.com/dream-machines-output/be82ce3d-ad30-4602-8397-d28a3b5d36a0/0_0.png : backup-jobs/upscale/264088956247867392/images/2022-11-24T10:47:03-06:00_be82ce3d-ad30-4602-8397-d28a3b5d36a0_0_0.png
INFO[0000] Fetched recent jobs user-id=264088956247867392 page=0 total-jobs=50
INFO[0000] Fetched recent jobs user-id=264088956247867392 page=1 total-jobs=100
INFO[0000] Fetched recent jobs user-id=264088956247867392 page=2 total-jobs=150
INFO[0000] Fetched recent jobs user-id=264088956247867392 page=3 total-jobs=200
INFO[0001] Fetched recent jobs user-id=264088956247867392 page=4 total-jobs=250
INFO[0001] Fetched recent jobs user-id=264088956247867392 page=5 total-jobs=300
INFO[0001] Fetched recent jobs user-id=264088956247867392 page=6 total-jobs=350
INFO[0001] Fetched recent jobs user-id=264088956247867392 page=7 total-jobs=400
INFO[0002] Fetched recent jobs user-id=264088956247867392 page=8 total-jobs=450
INFO[0002] Fetched recent jobs user-id=264088956247867392 page=9 total-jobs=500
INFO[0002] Fetched recent jobs user-id=264088956247867392 page=10 total-jobs=550
INFO[0002] Fetched recent jobs user-id=264088956247867392 page=11 total-jobs=600
INFO[0003] Fetched recent jobs user-id=264088956247867392 page=12 total-jobs=650
INFO[0003] Fetched recent jobs user-id=264088956247867392 page=13 total-jobs=700
INFO[0003] Fetched recent jobs user-id=264088956247867392 page=14 total-jobs=750
INFO[0004] Fetched recent jobs user-id=264088956247867392 page=15 total-jobs=800
INFO[0004] Fetched recent jobs user-id=264088956247867392 page=16 total-jobs=850
INFO[0004] Fetched recent jobs user-id=264088956247867392 page=17 total-jobs=885
...
INFO[0004] DONE Fetching image: https://storage.googleapis.com/dream-machines-output/67c71903-e28c-4b2c-a8ee-d080081060b5/0_0.png
INFO[0004] DONE Fetching image: https://storage.googleapis.com/dream-machines-output/67c71903-e28c-4b2c-a8ee-d080081060b5/0_0.png
```

### Recent Jobs
```
# Grab your browser Cookie session token from midjourney.com "__Secure-next-auth.session-token"
$ MIDJOURNEY_TOKEN="ey....."
$ mje midjourney recent-jobs --token "$MIDJOURNEY_TOKEN"  --user-id=264088956247867392 --fetch-all-jobs | tee /tmp/dougnukem-recent-midjourney-jobs.json
```
