from flask import Flask, jsonify, request

app = Flask(__name__)

projects = {}


@app.route('/project/<int:project_id>', methods=['GET'])
def get_project(project_id):
    if project_id in projects:
        # Fetch the project's repos by repo_id
        project = projects[project_id]
        project_repos = [repos[repo_id] for repo_id in project['repos']]  # Get repos associated with project
        project['repos'] = project_repos  # Include repos in the response
        return jsonify(project), 200
    return jsonify({"error": "Project not found"}), 404



@app.route('/projects', methods=['GET'])
def get_projects():
    return jsonify(list(projects.values())), 200


@app.route('/project/add', methods=['POST'])
def add_project():
    name = request.args.get('name')
    project_type = request.args.get('type')
    project_id = len(projects) + 1
    projects[project_id] = {"id": project_id, "name": name, "type": project_type}
    return jsonify({"message": "Project added", "project": projects[project_id]}), 201


@app.route('/project_repo/add', methods=['POST'])
def add_repo_to_project():
    project_name = request.args.get('project_name')
    repo_id = int(request.args.get('repo_id'))
    for project in projects.values():
        if project['name'] == project_name:
            if 'repos' not in project:
                project['repos'] = []
            project['repos'].append(repo_id)
            return jsonify({"message": "Repo added to project", "project": project}), 200
    return jsonify({"error": "Project not found"}), 404


if __name__ == '__main__':
    app.run(debug=True, port=5002)
