PWD = $(shell pwd)

docker:
	podman build -t discord-tofu .

run:
	podman run --mount type=bind,source=$(PWD)/config.toml,target=/app/config.toml:Z discord-tofu

build:
	go build -o discord-tofu ./main.go

pack: build
	zip -r tofu.zip \
		Makefile \
		discord-tofu \
		config.toml \
		pics/assets/ \
		brain/tofu-ai/*.{py,pickle}

install:
	python3 -m pip install --user nltk
	python3 -c "import nltk;\
		nltk.download('punkt');\
		nltk.download('wordnet');\
		nltk.download('averaged_perceptron_tagger');\
		nltk.download('stopwords');\
		nltk.download('omw-1.4')"
