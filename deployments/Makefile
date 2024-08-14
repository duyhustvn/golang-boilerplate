.PHONY: kafka-up


## KAFKA
kafka-zookeeper-up:
	cd kafka && docker compose rm -f kafka-kraft && docker compose build kafka-zookeeper && docker compose up zookeeper kafka-zookeeper

kafka-zookeeper-up-bg:
	cd kafka && docker compose rm -f kafka-kraft && docker compose build kafka-zookeeper && docker compose up zookeeper kafka-zookeeper -d

kafka-zookeeper-down:
	cd kafka && docker compose down

kafka-kraft-up-bg:
	cd kafka && docker compose build kafka-kraft && docker compose up kafka-kraft -d



## TIG
tig-up:
	cd TIG && docker compose up

tig-up-bg:
	cd TIG && docker compose build && docker compose up -d

telegraf-build:
	cd TIG && docker compose build telegraf

telegraf-up:
	make telegraf-build && cd TIG && docker compose up telegraf

telegraf-up-bg:
	cd TIG && make telegraf-build && docker compose up telegraf -d

grafana-up:
	cd TIG && docker compose up grafana


grafana-up-bg:
	cd TIG && docker compose up grafana -d
