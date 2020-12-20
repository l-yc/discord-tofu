build:
	go build -o discord-tofu ./main.go

pack: build
	zip -r tofu.zip \
		Makefile \
		discord-tofu \
		config.toml \
		pics/assets/ \
		answer/autorespond/tofu-ai/*.{py,pickle}

install:
	pip install --user nltk
	python -c "import nltk;\
		nltk.download('punkt');\
		nltk.download('wordnet');\
		nltk.download('averaged_perceptron_tagger');\
		nltk.download('stopwords')"
