import os
import subprocess
import random
import string
import sys
import ipaddress
import time
import configparser

import requests

# ===================== UTILITY FUNCTIONS =====================

def info(message): print(f"[INFO] {message}")
def log(msg): print(f"\033[1;32m[INFO]\033[0m {msg}")
def warn(msg): print(f"\033[1;33m[WARN]\033[0m {msg}")
def error(msg): print(f"\033[1;31m[ERROR]\033[0m {msg}")
def line(): print("=" * 50)

containers = ["postgres", "backend", "ml", "frontend", "interceptor"]
CONFIG_FILE = "gasha.conf"

# ===================== CONFIG HANDLER =====================

def generate_secret(length=32, hex=False):
    if hex:
        return ''.join(random.choices('0123456789abcdef', k=length))
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

def load_or_generate_config():
    config = configparser.ConfigParser()
    if os.path.exists(CONFIG_FILE):
        config.read(CONFIG_FILE)

    updated = False

    if "secrets" not in config:
        config["secrets"] = {}
        updated = True

    secrets = config["secrets"]

    if "DB_PASSWORD" not in secrets:
        secrets["DB_PASSWORD"] = "postgres"
        updated = True

    if "JWT_SECRET_KEY" not in secrets:
        secrets["JWT_SECRET_KEY"] = generate_secret(64, hex=True)
        updated = True

    if "WSKEY" not in secrets:
        secrets["WSKEY"] = generate_secret(64, hex=True)
        updated = True

    if "CAPTCHA_SECRET" not in secrets:
        secrets["CAPTCHA_SECRET"] = generate_secret(32)
        updated = True

    if updated:
        with open(CONFIG_FILE, "w") as configfile:
            config.write(configfile)

    return secrets

# ===================== NETWORK AND DOCKER HELPERS =====================

def run_command(command):
    try:
        output = subprocess.check_output(command, shell=True, stderr=subprocess.STDOUT)
        return output.decode().strip()
    except subprocess.CalledProcessError as e:
        warn(f"Command failed: {e.output.decode().strip()}")
        return None


def get_all_ips():
    output = subprocess.getoutput("ip -o -4 addr list")
    ips = [line.split()[3].split('/')[0] for line in output.splitlines()]
    try:
        public_ip = requests.get('https://api.ipify.org', timeout=5).text
        if public_ip not in ips:
            ips.append(public_ip)
    except requests.RequestException:
        pass

    return ips


def select_private_ip():
    ips = get_all_ips()
    if not ips:
        error("No private IPs found.")
        sys.exit(1)
    print("Available private IPs for frontend binding:")
    for idx, ip in enumerate(ips, start=1):
        print(f"{idx}) {ip}")
    print("=" * 50)
    choice = input(f"Choose IP [1-{len(ips)}]: ")
    if not choice.isdigit() or int(choice) < 1 or int(choice) > len(ips):
        error("Invalid selection!")
        sys.exit(1)
    return ips[int(choice) - 1]

def container_exists(name):
    output = subprocess.getoutput("docker ps -a --format '{{.Names}}'")
    return name in output.splitlines()

def remove_container(name):
    if container_exists(name):
        info(f"Removing {name} container...")
        run_command(f"docker rm -f {name}")

def remove_image(name):
    info(f"Removing image {name}...")
    run_command(f"docker rmi -f {name}")

def update_pg_hba_conf():
    log("Updating pg_hba.conf...")
    subprocess.run("docker restart postgres", shell=True)

# ===================== DEPLOYMENT =====================

def deploy_waf():
    line()
    log("Deploying GASHA WAF")
    line()

    secrets = load_or_generate_config()

    DB_PASSWORD = secrets["DB_PASSWORD"]
    JWT_SECRET_KEY = secrets["JWT_SECRET_KEY"]
    WSKEY = secrets["WSKEY"]
    CAPTCHA_SECRET = secrets["CAPTCHA_SECRET"]

    print("Secrets loaded (from gasha.conf or generated):")
    print(f"DB_PASSWORD: {DB_PASSWORD}")
    print(f"JWT_SECRET_KEY: {JWT_SECRET_KEY}")
    print(f"WSKEY: {WSKEY}")
    print(f"CAPTCHA_SECRET: {CAPTCHA_SECRET}")

    # frontend_ip = select_private_ip()
    os.system("docker network inspect waf_network >/dev/null 2>&1 || docker network create waf_network")

    if not container_exists("postgres"):
        log("Starting PostgreSQL...")
        os.system(f"docker run -d --name postgres --network waf_network --network-alias db "
                  f"-v waf_pgdata:/var/lib/postgresql/data -e POSTGRES_USER=postgres "
                  f"-e POSTGRES_PASSWORD={DB_PASSWORD} -e POSTGRES_DB=waf_db "
                  f"--restart unless-stopped postgres:16-alpine")
        update_pg_hba_conf()
    else:
        warn("PostgreSQL already exists. Skipping...")
    time.sleep(5)

    if not container_exists("backend"):
        log("Starting Backend...")
        os.system(f"docker run -d --name backend --network waf_network --network-alias backend "
                  f"-e DB_USER=postgres -e DB_PASSWORD={DB_PASSWORD} -e DB_HOST=db -e DB_PORT=5432 "
                  f"-e DB_NAME=waf_db -e DB_SSLMODE=disable -e WSKEY={WSKEY} -e JWT_SECRET_KEY={JWT_SECRET_KEY} "
                  f"--publish 127.0.0.1:8080:8080 --restart unless-stopped natnaelcrypto/waf-backend:latest")
    else:
        warn("Backend already exists. Skipping...")

    if not container_exists("ml"):
        log("Starting ML API...")
        os.system("docker run -d --name ml --network waf_network --network-alias mlapi "
                  "-e DEBUG=true -e BACKEND_API_URL=http://backend:8080/ "
                  "-e BACKEND_API_ML_MODELS_PATH=ml/models -e BACKEND_API_TYPE_ANALYSIS_PATH=ml/submit-analysis "
                  "-e ANOMALY_PREDICTOR_MODEL_PATH=ml_models/random_forest_model_v.0.2.1.pkl "
                  "-e TYPE_PREDICTOR_MODEL_PATH=ml_models/type_predictor_rf.joblib "
                  "-e BAD_WORDS_BY_TYPE_DIR=./words/words_by_type/ "
                  "-e COMMON_BAD_WORDS_PATH=./words/bad_words.txt "
                  "--restart unless-stopped --publish 5000:8090 natnaelcrypto/ml:latest")
    else:
        warn("ML API already exists. Skipping...")

    if not container_exists("frontend"):
        log(f"Starting Frontend ...")
        os.system(f"docker run -d --name frontend --network waf_network "
                  f"-e VITE_BACKEND_URL=http://backend:8080 "
                  f"-p 5173:5173 --restart unless-stopped natnaelcrypto/waf-frontend:latest")
    else:
        warn("Frontend already exists. Skipping...")

    if not container_exists("interceptor"):
        log("Starting Interceptor...")
        os.system(f"docker run -d --name interceptor --network host "
                  f"-e BACKENDURL=http://localhost:8080 "
                  f"-e MLHOSTURL=http://localhost:5000 "
                  f"-e CAPTCHA_SECRET={CAPTCHA_SECRET} -e WSKEY={WSKEY} "
                  "--restart unless-stopped natnaelcrypto/interceptor:latest")
    else:
        warn("Interceptor already exists. Skipping...")

# ===================== LIFECYCLE CONTROLS =====================

def start_all(): os.system("docker start postgres backend ml frontend interceptor")
def stop_all(): os.system("docker stop postgres backend ml frontend interceptor")
def restart_all(): os.system("docker restart postgres backend ml frontend interceptor")

def remove_all_except_db():
    os.system("docker rm -f backend ml frontend interceptor")
    os.system("docker rmi -f natnaelcrypto/waf-backend:latest natnaelcrypto/ml:latest "
              "natnaelcrypto/waf-frontend:latest natnaelcrypto/interceptor:latest")

def remove_all_including_db():
    remove_all_except_db()
    os.system("docker rm -f postgres")
    os.system("docker rmi -f postgres:16-alpine")
    os.system("docker volume rm waf_pgdata")

def remove_selected_container():
    print("Select container to remove:")
    for idx, name in enumerate(containers, 1):
        print(f"{idx}) {name}")
    try:
        choice = int(input(f"Choose [1-{len(containers)}]: "))
        if 1 <= choice <= len(containers):
            selected = containers[choice - 1]
            remove_container(selected)
            remove_image(selected)
        else:
            error("Invalid selection.")
    except ValueError:
        error("Invalid input.")

# ===================== MAIN MENU =====================

def menu():
    while True:
        line()
        print("📦 GASHA WAF Deployment Menu")
        line()
        print("1) Deploy WAF")
        print("2) Start All")
        print("3) Stop All")
        print("4) Restart All")
        print("5) Remove All Except DB")
        print("6) Remove All Including DB")
        print("7) Remove Selected Container")
        print("8) Exit")
        choice = input("Choose an option [1-8]: ")
        if choice == '1': deploy_waf()
        elif choice == '2': start_all()
        elif choice == '3': stop_all()
        elif choice == '4': restart_all()
        elif choice == '5': remove_all_except_db()
        elif choice == '6': remove_all_including_db()
        elif choice == '7': remove_selected_container()
        elif choice == '8':
            log("Exiting...")
            break
        else:
            error("Invalid option. Choose 1–8.")

menu()

