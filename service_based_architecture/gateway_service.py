from flask import Flask, redirect, url_for, request, jsonify
import requests

app = Flask(__name__)

WEB_PROTOCOL = "http://"

HOST_URL = "127.0.0.1"

USER_MANAGEMENT_PORT = "5001"
PROJECT_MANAGEMENT_PORT = "5002"
REPO_MANAGEMENT_PORT = "5003"

USER_SERVICE_URL = f"{WEB_PROTOCOL}{HOST_URL}:{USER_MANAGEMENT_PORT}"
PROJECT_SERVICE_URL = f"{WEB_PROTOCOL}{HOST_URL}:{PROJECT_MANAGEMENT_PORT}"
REPO_SERVICE_URL = f"{WEB_PROTOCOL}{HOST_URL}:{REPO_MANAGEMENT_PORT}"

@app.route('/user/<int:user_id>', methods=['GET'])
def get_user(user_id):
    response = requests.get(f'{USER_SERVICE_URL}/user/{user_id}')
    return response.json(), response.status_code


@app.route('/users', methods=['GET'])
def get_users():
    response = requests.get(f'{USER_SERVICE_URL}/users')
    return response.json(), response.status_code


@app.route('/project/<int:project_id>', methods=['GET'])
def get_project(project_id):
    response = requests.get(f'{PROJECT_SERVICE_URL}/project/{project_id}')
    return response.json(), response.status_code


@app.route('/projects', methods=['GET'])
def get_projects():
    response = requests.get(f'{PROJECT_SERVICE_URL}/projects')
    return response.json(), response.status_code


@app.route('/repo/<int:repo_id>', methods=['GET'])
def get_repo(repo_id):
    response = requests.get(f'{REPO_SERVICE_URL}/repo/{repo_id}')
    return response.json(), response.status_code


@app.route('/repos', methods=['GET'])
def get_repos():
    response = requests.get(f'{REPO_SERVICE_URL}/repos')
    return response.json(), response.status_code


@app.route('/user/username=<username>', methods=['POST'])
def add_user(username):
    response = requests.post(f'{USER_SERVICE_URL}/user/username={username}')
    return response.json(), response.status_code


@app.route('/project/add', methods=['POST'])
def add_project():
    name = request.args.get('name')
    project_type = request.args.get('type')
    response = requests.post(f'{PROJECT_SERVICE_URL}/project/add?name={name}&type={project_type}')
    return response.json(), response.status_code


@app.route('/repo/add', methods=['POST'])
def add_repo():
    name = request.args.get('name')
    response = requests.post(f'{REPO_SERVICE_URL}repo/add?name={name}')
    return response.json(), response.status_code


@app.route('/project_repo/add', methods=['POST'])
def add_repo_to_project():
    project_name = request.args.get('project_name')
    repo_id = request.args.get('repo_id')
    response = requests.post(f'{PROJECT_SERVICE_URL}/project_repo/add?project_name={project_name}&repo_id={repo_id}')
    return response.json(), response.status_code


@app.route('/user_repo/add', methods=['POST'])
def add_repo_to_user():
    username = request.args.get('username')
    repo_id = int(request.args.get('repo_id'))

    response = requests.post(f'{USER_SERVICE_URL}/user_repo/add?username={username}&repo_id={repo_id}')

    return response.json(), response.status_code


if __name__ == '__main__':
    app.run(debug=True, port=5000)
