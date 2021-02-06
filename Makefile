.PHONY: venv

NAME    := acrobits-websvc

venvdir := ./venv
appdir  := ./app

all: venv

define make_venv
	python3 -m venv --prompt $(NAME) $(1)
	( \
		source $(1)/bin/activate; \
		pip install -r $(appdir)/requirements.txt; \
		deactivate; \
	)
endef

venv:
	$(call make_venv,$(venvdir))

clean:
	rm -rf $(venvdir)/

docker:
	docker build -t $(NAME) $(appdir)/

clean-docker:
	docker rmi $(NAME)
