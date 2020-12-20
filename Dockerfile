FROM python:latest

RUN mkdir -p /flask && \
    pip install flask && \
    pip install pymysql


COPY flask/ /flask/

CMD ["python", "/flask/main.py"]