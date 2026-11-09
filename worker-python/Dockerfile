FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY robotdatahub_worker ./robotdatahub_worker
CMD ["python","-m","robotdatahub_worker.main"]
