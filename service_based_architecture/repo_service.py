from flask import Flask, jsonify, request

app = Flask(__name__)

repos = {}


@app.route('/repo/<int:repo_id>', methods=['GET'])
def get_repo(repo_id):
    if repo_id in repos:
        return jsonify(repos[repo_id]), 200
    return jsonify({"error": "Repo not found"}), 404


@app.route('/repos', methods=['GET'])
def get_repos():
    return jsonify(list(repos.values())), 200


@app.route('/repo/add', methods=['POST'])
def add_repo():
    name = request.args.get('name')
    repo_id = len(repos) + 1
    repos[repo_id] = {"id": repo_id, "name": name}
    return jsonify({"message": "Repo added", "repo": repos[repo_id]}), 201


if __name__ == '__main__':
    app.run(debug=True, port=5003)
