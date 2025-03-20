FROM python:3.9.7-slim
COPY server /
RUN apt-get update && apt-get install -y netcat-openbsd bash && rm -rf /var/lib/apt/lists/*
RUN python -m unittest tests/test_common.py
ENTRYPOINT ["/bin/sh"]
