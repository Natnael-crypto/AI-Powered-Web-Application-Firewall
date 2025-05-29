import os
import subprocess

def info(message): print(f"[INFO] {message}")
def warn(message): print(f"[WARN] {message}")
def error(message): print(f"[ERROR] {message}")

containers = ["postgres", "backend", "ml", "frontend", "interceptor"]

def run_command(command):
    try:
        output = subprocess.check_output(command, shell=True, stderr=subprocess.STDOUT)
        return output.decode().strip()
    except subprocess.CalledProcessError as e:
        warn(f"Command failed: {e.output.decode().strip()}")
        return None

def container_exists(name):
    result = run_command(f"docker ps -a --format '{{{{.Names}}}}' | grep -w {name}")
    return result is not None

def stop_container(name):
    if container_exists(name):
        info(f"Stopping {name}...")
        run_command(f"docker stop {name}")

def start_container(name):
    if container_exists(name):
        info(f"Starting {name}...")
        run_command(f"docker start {name}")

def remove_container(name):
    if container_exists(name):
        info(f"Removing {name} container...")
        run_command(f"docker rm -f {name}")

def remove_image(name):
    info(f"Removing image {name}...")
    run_command(f"docker rmi -f {name}")

def update_pg_hba_conf():
    info("Updating pg_hba.conf for external access...")
    run_command("docker exec postgres sh -c \"echo 'host all all 0.0.0.0/0 scram-sha-256' >> /var/lib/postgresql/data/pg_hba.conf\"")
    run_command("docker exec postgres pg_ctl reload")

def deploy_waf():
    info("Deploying GASHA WAF...")
    # Replace these with your actual deployment commands
    run_command("docker run -d --name postgres postgres:latest")
    update_pg_hba_conf()
    run_command("docker run -d --name backend my-backend")
    run_command("docker run -d --name ml my-ml")
    run_command("docker run -d --name frontend my-frontend")
    run_command("docker run -d --name interceptor my-interceptor")

def start_all():
    for name in containers:
        start_container(name)

def stop_all():
    for name in containers:
        stop_container(name)

def restart_all():
    for name in containers:
        stop_container(name)
        start_container(name)

def remove_all_except_db():
    for name in containers:
        if name != "postgres":
            remove_container(name)
            remove_image(name)

def remove_all_including_db():
    for name in containers:
        remove_container(name)
        remove_image(name)

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

def menu():
    while True:
        print("=" * 50)
        print("📦 GASHA WAF Deployment Menu")
        print("=" * 50)
        print("1) Deploy WAF")
        print("2) Start All")
        print("3) Stop All")
        print("4) Restart All")
        print("5) Remove All Except DB")
        print("6) Remove All Including DB")
        print("7) Remove Selected Container")
        print("8) Exit")
        print("=" * 50)
        try:
            choice = input("Choose an option [1-8]: ").strip()
        except EOFError:
            error("Input failed. Please run this script in a terminal.")
            break

        if choice == '1':
            deploy_waf()
        elif choice == '2':
            start_all()
        elif choice == '3':
            stop_all()
        elif choice == '4':
            restart_all()
        elif choice == '5':
            remove_all_except_db()
        elif choice == '6':
            remove_all_including_db()
        elif choice == '7':
            remove_selected_container()
        elif choice == '8':
            print("Exiting...")
            break
        else:
            error("Invalid option. Please choose from 1 to 8.")

if __name__ == "__main__":
    menu()
