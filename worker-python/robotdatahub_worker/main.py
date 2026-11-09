import logging
import os
import time
import requests
from robotdatahub_worker.jobs.runner import execute

logging.basicConfig(level=logging.INFO, format='%(asctime)s %(levelname)s %(message)s')
BASE = os.environ.get('SERVER_URL', 'http://server:8080').rstrip('/')
TOKEN = os.environ['API_TOKEN']
HEADERS = {'X-Worker-Token': TOKEN}

def once():
    r = requests.post(f'{BASE}/internal/v1/jobs/claim', headers=HEADERS, timeout=20)
    if r.status_code == 204:
        return False
    r.raise_for_status()
    job = r.json()['data']
    try:
        result = execute(job['dataset_id'])
        payload = {'result': result, 'error': ''}
    except Exception as exc:
        logging.exception('Job %s failed', job['id'])
        payload = {'result': None, 'error': str(exc)[:2000]}
    r = requests.post(f"{BASE}/internal/v1/jobs/{job['id']}/complete", json=payload, headers=HEADERS, timeout=20)
    r.raise_for_status()
    return True

if __name__ == '__main__':
    while True:
        try:
            if not once():
                time.sleep(3)
        except requests.RequestException:
            logging.exception('Cannot contact server')
            time.sleep(5)
