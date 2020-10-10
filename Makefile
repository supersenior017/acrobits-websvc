.PHONY: venv

NAME := acrobits-websvc
VENV := ./venv

all: venv

define make_venv
	python3 -m venv --prompt $(NAME) $(1)
	( \
		source $(1)/bin/activate; \
		pip install -r "./src/requirements.txt"; \
		deactivate; \
	)
endef

venv:
	$(call make_venv,$(VENV))

clean:
	rm -rf $(VENV)/

docker:
	docker build -t $(NAME) "./src/"

clean-docker:
	docker rmi $(NAME)
