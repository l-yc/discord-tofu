#!/bin/sh
pip install --user nltk
python -c "import nltk;\
  nltk.download('punkt');\
  nltk.download('wordnet');\
  nltk.download('averaged_perceptron_tagger')"
