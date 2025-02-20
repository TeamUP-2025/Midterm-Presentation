import json
import os

class TeamData:
    def __init__(self, filename="team_data.json"):
        self.filename = filename
        self._initialize_file()

    def _initialize_file(self):
        """Initialize the JSON file if it does not exist."""
        if not os.path.exists(self.filename):
            with open(self.filename, "w") as file:
                json.dump([], file)

    def _read_data(self):
        """Read and return the team members from the JSON file."""
        with open(self.filename, "r") as file:
            return json.load(file)

    def _write_data(self, members):
        with open(self.filename, "w") as file:
            json.dump(members, file)

    def add_member(self, member):
        members = self._read_data()
        if member not in members:
            members.append(member)
            self._write_data(members)
            return True
        return False

    def remove_member(self, member):
        members = self._read_data()
        if member in members:
            members.remove(member)
            self._write_data(members)
            return True
        return False

    def list_members(self):
        return self._read_data()


class TeamManager:
    def __init__(self, repository: TeamData):
        self.repository = repository

    def add_member_to_team(self, member):
        if self.repository.add_member(member):
            return f"Member '{member}' added successfully."
        return f"Member '{member}' already exists in the team."

    def remove_member_from_team(self, member):
        if self.repository.remove_member(member):
            return f"Member '{member}' removed successfully."
        return f"Member '{member}' not found in the team."

    def get_team_members(self):
        return self.repository.list_members()


class TerminalPresentation:
    def __init__(self, domain: TeamManager):
        self.domain = domain

    def add_member(self, member):
        result = self.domain.add_member_to_team(member)
        print(result)

    def remove_member(self, member):
        result = self.domain.remove_member_from_team(member)
        print(result)

    def show_team_members(self):
        members = self.domain.get_team_members()
        if members:
            print("Team Members: " + ", ".join(members))
        else:
            print("The team is currently empty.")


if __name__ == "__main__":
    repository = TeamData()
    domain = TeamManager(repository)
    presentation = TerminalPresentation(domain)

    presentation.add_member("Alice")
    presentation.add_member("Bob")
    presentation.show_team_members()
