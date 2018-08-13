build:
	rm -rf backup.go || true
	rm -rf cron_script.sh || true
	cp  ../backup.go backup.go
	cp  ../cron_script.sh cron_script.sh
	docker build -t nutellinoit/sidecar-backup-mysql:dev .
	rm -rf backup.go || true
	rm -rf cron_script.sh || true
