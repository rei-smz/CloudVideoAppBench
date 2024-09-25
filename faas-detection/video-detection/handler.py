import asyncio
import logging
import os
import shutil
import time
import json
from tempfile import NamedTemporaryFile
from minio import Minio
from minio.error import S3Error
from . import ssd
from . import resnet

def get_secret(key):
    with open("/var/openfaas/secrets/{}".format(key)) as f:
        return f.read().strip()

if os.name == 'nt':
  SHELL = 'pwsh'
else:
  SHELL = 'sh'

MINIO_ENDPOINT = os.getenv('MINIO_ENDPOINT')
MINIO_ACCESS_KEY = get_secret("minio-access-key")
MINIO_SECRET_KEY = get_secret("minio-secret-key")
BUCKET_NAME = os.getenv('BUCKET_NAME')

minio_client = Minio(
    MINIO_ENDPOINT,
    access_key=MINIO_ACCESS_KEY,
    secret_key=MINIO_SECRET_KEY,
    secure=False
)

def handle(event, context):
    if event.method != "POST":
      return {
        "statusCode": 405,
        "body": "Method not allowed"
      }
    
    try:
        body = json.loads(event.body.decode('utf-8'))
    except json.JSONDecodeError:
        return {"statusCode": 400, "body": "Invalid JSON"}
    
    path = body.get('path')
    obj_name = body.get('object')
    args = body.get('args', {})
    mode = args.get('mode', 'lw')
    current_time = time.strftime("%Y-%m-%d-%H%M%S", time.localtime())

    # get video
    tmp_path = current_time + '-' + path.split('/')[-1]
    os.mkdir(tmp_path)
    try:
        minio_client.fget_object(BUCKET_NAME, path + '/' + obj_name, tmp_path + '/' + obj_name)
    except S3Error as e:
        logging.error(f"MinIO Error: {e}")
        return {"statusCode": 500, "body": f"Error downloading file: {e}"}

    if not path:
        return {"statusCode": 400, "body": "Path is required"}

    response = None
    if mode == 'lw':
        response = asyncio.run(ssd.run(tmp_path + '/' + obj_name))
    elif mode == 'hp':
        response = asyncio.run(resnet.run(tmp_path + '/' + obj_name))
    else:
        return {"statusCode": 400, "body": "Mode is invalid"}
    
    return {
        "statusCode": 200,
        "body": response,
        "headers": {
            "Content-Type": "application/json"
        }
    }
