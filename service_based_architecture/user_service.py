from flask import Flask, jsonify, request

app = Flask(__name__)

users = {}


@app.route('/user/<int:user_id>', methods=['GET'])
def get_user(username):
    if username in users:
        # Fetch the user's repos by repo_id
        user = users[username]
        user_repos = [repos[repo_id] for repo_id in user['repos']]  # Get repos associated with user
        user['repos'] = user_repos  # Include repos in the response
        return jsonify(user), 200
    return jsonify({"error": "User not found"}), 404


@app.route('/users', methods=['GET'])
def get_users():
    return jsonify(list(users.values())), 200


@app.route('/user/username=<username>', methods=['POST'])
def add_user(username):
    user_id = len(users) + 1
    users[user_id] = {"id": user_id, "username": username}
    return jsonify({"message": "User added", "user": users[user_id]}), 201


@app.route('/user_repo/add', methods=['POST'])
def add_repo_to_user():
    username = request.args.get('username')
    repo_id = int(request.args.get('repo_id'))
    for user in users.values():
        if user['username'] == username:
            if 'repos' not in user:
                user['repos'] = []
            user['repos'].append(repo_id)
            return jsonify({"message": "Repo added to user", "user": user}), 200
    return jsonify({"error": "User not found"}), 404


if __name__ == '__main__':
    app.run(debug=True, port=5001)
